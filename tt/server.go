package http

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strconv"
	//"strings"
	"sync"
	//"time"
)

type Server struct {
	addr string

	// 为来到此服务器的请求指定 context 上下文 ，不设就是 context.Background()
	// 如果BaseContext为nil，则默认值为context.Background()。
	// 如果为非nil，则它必须返回非nil上下文
	Context func(net.Listener) context.Context

	// 互斥锁 保证资源的安全
	mux sync.Mutex

	// 监听socket表
	//listeners map[*net.Listener]struct{}

	ln net.Listener

	TLSConfig *tls.Config

	// 存活的客户端链接表
	// activeConn map[*conn]struct{}
}

func (s *Server) ListenAndServe() error {
	if "" == s.addr {
		s.addr = ":http"
	}
	fmt.Println("Listen", s.addr)

	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	ln = ln.(*net.TCPListener)
	return s.Serve(ln)
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	//fmt.Println(s.shuttingDown())
	if "" == s.addr {
		s.addr = ":https"
	}
	fmt.Println("Listen", s.addr)

	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	//sln := ln.(*net.TCPListener)
	return s.ServeTLS(ln, certFile, keyFile)
}

func (s *Server) Serve(ln net.Listener) error {
	s.ln = ln //.(*net.TCPListener)
	defer s.ln.Close()

	for {
		// rw, _ := s.ln.Accept()
		// c := newConn(rw, s)
		// go c.serve()
		rwc, _ := ln.Accept()

		go serve(rwc)
	}
}

// cloneTLSConfig returns a shallow clone of cfg, or a new zero tls.Config if
// cfg is nil. This is safe to call even if cfg is in active use by a TLS
// client or server.
func cloneTLSConfig(cfg *tls.Config) *tls.Config {
	if cfg == nil {
		return &tls.Config{}
	}
	return cfg.Clone()
}

func strSliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func (s *Server) ServeTLS(ln net.Listener, certFile, keyFile string) error {
	//defer s.ln.Close()

	config := cloneTLSConfig(s.TLSConfig)
	if !strSliceContains(config.NextProtos, "http/1.1") {
		config.NextProtos = append(config.NextProtos, "http/1.1")
	}

	configHasCert := len(config.Certificates) > 0 || config.GetCertificate != nil
	if !configHasCert || certFile != "" || keyFile != "" {
		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
	}

	tlsLn := tls.NewListener(ln, config)
	return s.Serve(tlsLn)
}

func serve(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("remoteAddr:", remoteAddr)
	localAddr := conn.LocalAddr().String()
	fmt.Println("localAddr:", localAddr)

	cr := &cReader{c: conn}
	cr.cc = bufio.NewReader(conn)
	cr.setInfiniteReadLimit()
	//var buffer []byte
	for {
		//line, prefix, err := cr.cc.ReadLine()
		//fmt.Println(len(line), prefix, err)
		//time.Sleep(time.Second)
		//buf := make([]byte, 4096)
		//n, _ := cr.Read(buf)
		//buffer = append(buffer, buf[:n]...)
		//fmt.Println(n)
		sResponse(conn, cr)
		cr.backgroundRead()
		return
	}
}

var textprotoReaderPool sync.Pool

func sResponse(conn net.Conn, cr *cReader) {
	response := "HTTP/2.0 200 OK\r\n"

	cc := cr.cc
	tp := textproto.NewReader(cc)
	s, _ := tp.ReadLine()
	fmt.Println(s)
	mimeHeader, _ := tp.ReadMIMEHeader()
	length := 0
	if ln, ok := mimeHeader["Content-Length"]; ok {
		length, _ = strconv.Atoi(ln[0])
	}

	for n, v := range mimeHeader {
		fmt.Printf("Header: %s :\t%v\n", n, v)
	}

	var buffer []byte
	fmt.Println("Length", length)

	var body []byte
	x, buf := 0, make([]byte, 4096)
	for x < length {
		n, e := cr.cc.Read(buf[:])
		if e != nil {
			break
		}
		x += n
		body = append(body, buf[:n]...)
	}
	//body, _ = cr.cc.Peek(length)

	b := fmt.Sprintf("%s <hr /> (%d.%d)\r\n%s", string(body), len(body), length, string(buffer))
	response += fmt.Sprintf("Content-Length: %d\r\n", len(b))
	response += fmt.Sprintf("Content-Type: text/html\r\n")
	//response += fmt.Sprintf("Content-Encoding: gzip\r\n")
	response += "\r\n" + b
	conn.Write([]byte(response))
}

type cReader struct {
	c net.Conn

	cc *bufio.Reader
	// 互斥锁
	mux  sync.Mutex
	cond *sync.Cond

	hasByte bool
	byteBuf [1]byte
	isRead  bool
	aborted bool
	remain  int64
}

func (cr *cReader) lock() {
	cr.mux.Lock()
	if nil == cr.cond {
		cr.cond = sync.NewCond(&cr.mux)
	}
}
func (cr *cReader) unlock() { cr.mux.Unlock() }
func (cr *cReader) wait()   { cr.cond.Wait() }

func (cr *cReader) Read(b []byte) (n int, err error) {
	cr.lock()

	if cr.isRead {
		cr.unlock()
		panic("无效的并发Read")
	}

	// cr.hitReadLimit()
	if cr.remain <= 0 {
		cr.wait()
		cr.unlock()
		return 0, io.EOF
	}

	if len(b) == 0 {
		cr.wait()
		cr.unlock()
		return 0, nil
	}

	if int64(len(b)) > cr.remain {
		b = b[:cr.remain]
	}

	// 如果存在数据
	if cr.hasByte {
		b[0] = cr.byteBuf[0]
		cr.hasByte = false
		//cr.unlock()
		//return 1, nil
	}

	cr.isRead = true
	cr.unlock()
	n, err = cr.c.Read(b)

	cr.lock()
	cr.isRead = false
	if err != nil {
		fmt.Println(err)
	}

	cr.remain -= int64(n)
	cr.unlock()

	cr.cond.Broadcast()
	return n, err
}
func (cr *cReader) hitReadLimit() bool { return cr.remain <= 0 }

// 设置读取长度
func (cr *cReader) setReadLimit(remain int64) { cr.remain = remain }
func (cr *cReader) setInfiniteReadLimit()     { cr.remain = 1<<63 - 1 }

func (cr *cReader) backgroundRead() {
	n, err := cr.c.Read(cr.byteBuf[:])
	cr.lock()
	if 1 == n && err != io.EOF {
		cr.hasByte = true
	}
	if err != nil && err == io.EOF {
		fmt.Println("bgr> ", err)
		cr.wait()
	}
	cr.aborted = false
	cr.isRead = false
	cr.unlock()
	cr.cond.Broadcast()
}
