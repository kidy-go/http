package http

import (
	"bufio"
	"io"
)

type Conn struct {
	buf *bufio.ReadWriter
	rwc io.ReadWriter
}

func newConn(rwc io.ReadWriter, server bool) *Conn {
	conn := &Conn{
		rwc: rwc,
	}

	readBufSize, writeSize := defaultBufSize, 0
	if readBufSize < minReadBufSize {
		readBufSize = minReadBufSize
	}

	conn.buf = bufio.NewReadWriter(
		bufio.NewReaderSize(rwc, readBufSize),
		bufio.NewWriterSize(rwc, writeSize),
	)

	return conn
}
