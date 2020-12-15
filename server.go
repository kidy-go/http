// server.go kee > 2020/12/14

package http

import (
	// "errors"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

type FrameHeader struct {
	Length uint32
	Type   string
	Flags  string
}

type Server struct {
}

func ListenAndServeTLS(addr string, pem string, key string) error {
	crt, err := tls.LoadX509KeyPair(pem, key)
	if err != nil {
		return err
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader

	listen, err := tls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		go handlerConnection(conn)
	}
}

func ListenAndServe(addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		go handlerConnection(conn)
	}
}

func handlerConnection(conn net.Conn) {
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("conn>read (fail):", err)
			return
		}

		defer conn.Close()

		fmt.Printf("%v\n%s\n", buffer[:n], buffer[:n])
		fmt.Println(buffer[:9], string(buffer[:9]))

		body := []byte("Hello world")
		length := len(body)

		response := "HTTP/1.1 200 OK\r\n"
		response += "Date: Tue, 10 Jul 2020 00:00:00 GMT\r\n"
		response += fmt.Sprintf("Content-Length: %d\r\n", length)
		response += fmt.Sprintf("Content-Type: text/html\r\n\r\n")
		response += fmt.Sprintf("%s", body)
		fmt.Println("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
		conn.Write([]byte(response))
		return
	}
}
