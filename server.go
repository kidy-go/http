// server.go kee > 2020/12/14

package http

import (
	// "errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

type Server struct {
	addr    string
	ln      net.Listener
	headers Headers
}

func ListenAndServe(addr string) error {
	s := &Server{addr: addr}
	return s.ListenAndServe()
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

	raw := connReader(conn)

	ioutil.WriteFile("./form-data.raw", raw, 0755)
	go reader(raw)

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

func reader(raw []byte) {
	headers := ReadHeaders(raw)
	for name, value := range headers.Payload {
		fmt.Println(name, "=>", value)
	}

	bodyBegin := headers.Length + len("\r\n\r\n")
	body := string(raw[bodyBegin:])

	if "multipart/form-data" == headers.ContentType {
		fmt.Println("Boundary:", headers.Boundary)
		for _, fraw := range strings.Split(body, "--"+headers.Boundary) {
			n := len(fraw)
			if n > 100 {
				n = 100
			}
			fmt.Println(fraw[:n])
			if fn := strings.Index(fraw[:n], "filename="); fn > -1 {
				fn += len("filename=") + 1
				fname := fraw[fn:]
				efn := strings.Index(fname, "\"")
				fname = "x-" + fname[:efn]
				rn := strings.Index(fraw, "\r\n\r\n") + len("\r\n\r\n")

				fmt.Println("FileName: ", []string{fname})
				f, e := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0600)
				defer f.Close()
				_checkErr(e)
				fmt.Println([]byte(fraw[rn : len(fraw)-2]))
				f.Write([]byte(fraw[rn : len(fraw)-2]))
				break
			}
		}
	}
}

func _checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
