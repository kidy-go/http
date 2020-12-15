// server.go kee > 2020/12/14

package http

import (
	// "errors"
	"fmt"
	"io"
	"net"
)

type Server struct {
	addr    string
	ln      net.Listener
	headers Headers
}

func (s *Server) ListenAndServe() error {
	if "" == s.addr {
		s.addr = ":http"
	}

	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.ln = ln.(*net.TCPListener)
	return s.Serve()
}

func (s *Server) Serve() error {
	defer s.ln.Close()

	for {
		rwc, err := s.ln.Accept()
		if err != nil {
			continue
			//return err
		}
		go handlerConnection(rwc)
	}
}

func handlerConnection(conn net.Conn) {
	defer conn.Close()

	connReader(conn)
	response := ""
	response += "HTTP/1.1 200 OK\r\n"
	response += fmt.Sprintf("Content-Type: text/html\r\n\r\n")
	response += "Hello world"
	conn.Write([]byte(response))
}

func connReader(conn net.Conn) (buffer []byte) {
	length, maxSize := 0, 0
	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}

		buffer = append(buffer, buf[:n]...)
		length += n
		if maxSize <= 0 {
			headers := ReadHeaders(buffer)
			maxSize = headers.ContentLength
		}
		if maxSize > 0 && maxSize-1 < length {
			return
		}
	}
	return
}
