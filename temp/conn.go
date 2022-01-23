package http

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"sync"
	"time"
)

type Header map[string][]string

const maxInt64 = 1<<63 - 1

type conn struct {
	// 连接绑定的Server 不可变,不可为nil
	serv *Server

	// 用于取消连接的上下文
	cancelCtx context.Context

	// rwc是底层网络连接。
	// 它从不被其他类型包装，并且是分配给CloseNotifier调用方的值。
	// 通常类型为*net.TCPConn或*tls.Conn。
	rwc net.Conn

	// r是bufr的读取源。它是rwc的包装器，它提供io.LimitedReader-style限制(在读取请求头时)
	// 支持CloseNotifier的功能。参见*connReader文档。
	r *connReader

	//  bufr从r读取 还是前面的reader，加了缓冲
	bufr *bufio.Reader

	// bufw 写入 checkConnErrorWriter{c}, 当发生错误时填充 werr。
	bufw *bufio.Writer

	remoteAddr string

	// mu guards hijackedv
	mux sync.Mutex

	// 表示这个连接是否被一个带有Hijacker接口的 Handler hijacked了，主要用于切换协议
	hijackedv bool
}

func newConn(rwc net.Conn, serv *Server) *conn {
	c := &conn{
		serv: serv,
		rwc:  rwc,
	}
	return c
}

func (c *conn) serve() {
	c.remoteAddr = c.rwc.RemoteAddr().String()
	c.r = &connReader{conn: c}
	c.bufr = newBufioReader(c.r)
	//buf := make([]byte, 100)
	var buffer []byte
	for {
		//n, _ := c.rwc.Read(buf)
		//buffer = append(buffer, buf[:n]...)
		c.readRequest()
		simpleResponse(c.rwc, buffer)
		c.r.background()
		//c.r.startBackgroundRead()
	}
}

func simpleResponse(conn net.Conn, buffer []byte) {
	response := "HTTP/1.1 200 OK\r\n"

	b := "Hello world"
	response += fmt.Sprintf("Content-Length: %d\r\n", len(b))
	response += fmt.Sprintf("Content-Type: text/html\r\n")
	//response += fmt.Sprintf("Content-Encoding: gzip\r\n")
	response += "\r\n" + b
	fmt.Println(">", len(buffer))

	conn.Write([]byte(response))
}
func (c *conn) hijacked() bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.hijackedv
}

func (c *conn) readRequest() {
	if c.hijacked() {
		fmt.Println("Hijacked")
		return
	}
	c.r.setReadLimit(4096)
	readRequest(c.bufr)
}

func readRequest(b *bufio.Reader) Header {
	tp := newTextprotoReader(b)

	// First line: GET /index.html HTTP/1.0
	var s string
	var err error
	if _, err = tp.ReadLine(); err != nil {
		return nil
	}
	fmt.Println(s)

	mimeHeader, _ := tp.ReadMIMEHeader()
	fmt.Println(len(mimeHeader))
	return Header{}
}

var textprotoReaderPool sync.Pool

func newTextprotoReader(br *bufio.Reader) *textproto.Reader {
	if v := textprotoReaderPool.Get(); v != nil {
		tr := v.(*textproto.Reader)
		tr.R = br
		return tr
	}
	return textproto.NewReader(br)
}

var bufioReaderPool sync.Pool

func newBufioReader(r io.Reader) *bufio.Reader {
	if v := bufioReaderPool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	return bufio.NewReader(r)
}

type connReader struct {
	conn *conn

	mu      sync.Mutex // guards following
	hasByte bool
	byteBuf [1]byte
	cond    *sync.Cond
	inRead  bool
	aborted bool  // set true before conn.rwc deadline is set to past
	remain  int64 // bytes remaining
}

func (cr *connReader) background() bool {
	n, err := cr.conn.rwc.Read(cr.byteBuf[:])
	cr.lock()

	if n == 1 {
		cr.hasByte = true
	}
	if ne, ok := err.(net.Error); ok && cr.aborted && ne.Timeout() {

	} else if err != nil {
		return true
	}
	cr.aborted = false
	cr.inRead = false
	cr.unlock()
	cr.cond.Broadcast()
	return false
}
func (cr *connReader) lock() {
	cr.mu.Lock()
	if cr.cond == nil {
		cr.cond = sync.NewCond(&cr.mu)
	}
}
func (cr *connReader) Read(p []byte) (n int, err error) {
	cr.lock()
	if cr.inRead {
		cr.unlock()
		if cr.conn.hijacked() {
			panic("invalid Body.Read call. After hijacked, the original Request must not be used")
		}
		panic("invalid concurrent Body.Read call")
	}
	if cr.hitReadLimit() {
		cr.unlock()
		return 0, io.EOF
	}
	if len(p) == 0 {
		cr.unlock()
		return 0, nil
	}
	if int64(len(p)) > cr.remain {
		p = p[:cr.remain]
	}
	if cr.hasByte {
		p[0] = cr.byteBuf[0]
		cr.hasByte = false
		cr.unlock()
		return 1, nil
	}
	cr.inRead = true
	cr.unlock()
	n, err = cr.conn.rwc.Read(p)

	cr.lock()
	cr.inRead = false
	if err != nil {
		cr.handleReadError(err)
	}
	cr.remain -= int64(n)
	cr.unlock()

	cr.cond.Broadcast()
	return n, err
}
func (cr *connReader) startBackgroundRead() {
	cr.lock()
	defer cr.unlock()
	if cr.inRead {
		panic("invalid concurrent Body.Read call")
	}
	if cr.hasByte {
		return
	}
	cr.inRead = true
	cr.conn.rwc.SetReadDeadline(time.Time{})
	go cr.backgroundRead()
}
func (cr *connReader) backgroundRead() {
	n, err := cr.conn.rwc.Read(cr.byteBuf[:])
	cr.lock()
	if n == 1 {
		cr.hasByte = true
	}
	if ne, ok := err.(net.Error); err != nil && (!ok || !cr.aborted || !ne.Timeout()) {
		cr.handleReadError(err)
	}
	cr.aborted = false
	cr.inRead = false
	cr.unlock()
	cr.cond.Broadcast()
}
func (cr *connReader) abortPendingRead() {
	cr.lock()
	defer cr.unlock()
	if !cr.inRead {
		return
	}
	cr.aborted = true
	//cr.conn.rwc.SetReadDeadline(aLongTimeAgo)
	for cr.inRead {
		cr.cond.Wait()
	}
	cr.conn.rwc.SetReadDeadline(time.Time{})
}
func (cr *connReader) closeNotify() {
	//res, _ := cr.conn.curReq.Load().(*response)
	//if res != nil && atomic.CompareAndSwapInt32(&res.didCloseNotify, 0, 1) {
	//res.closeNotifyCh <- true
	//}
}
func (cr *connReader) unlock()                   { cr.mu.Unlock() }
func (cr *connReader) setReadLimit(remain int64) { cr.remain = remain }
func (cr *connReader) setInfiniteReadLimit()     { cr.remain = maxInt64 }
func (cr *connReader) hitReadLimit() bool        { return cr.remain <= 0 }
func (cr *connReader) handleReadError(_ error) {
	//cr.conn.cancelCtx()
	cr.closeNotify()
}
