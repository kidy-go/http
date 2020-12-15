package http

import (
	"strconv"
	"strings"
)

type Headers struct {
	ContentType   string
	ContentLength int
	Connection    string
	Host          string
	Referer       string
	UserAgent     string
	Length        int
	Payload       map[string][]string
}

func ReadHeaders(buffer []byte) Headers {
	headers := Headers{ContentLength: -1}

	if pos := strings.Index(string(buffer), "\r\n\r\n"); pos > -1 {
		headers.Length = pos
		headers.ContentLength = pos
		for _, s := range strings.Split(string(buffer[:pos]), "\r\n") {
			if ex := strings.Split(s, ":"); len(ex) > 1 {
				n, v := strings.Trim(ex[0], " "), strings.Trim(ex[1], " ")

				switch strings.ToLower(n) {
				case "content-type":
					headers.ContentType = v
				case "content-length":
					if n, e := strconv.Atoi(v); e == nil {
						headers.ContentLength = n
					}
				case "connection":
					headers.Connection = v
				case "host":
					headers.Host = v
				case "referer":
					headers.Referer = v
				case "user-agent":
					headers.UserAgent = v
				}

				var val []string
				for _, _v := range strings.Split(v, ";") {
					val = append(val, strings.Trim(_v, " "))
				}
				headers.Set(strings.ToLower(n), val)
			}
		}
	}
	return headers
}

func (h Headers) Set(name string, value []string) {
	if nil == h.Payload {
		h.Payload = make(map[string][]string)
	}

	if _, ok := h.Payload[name]; ok {
		h.Payload[name] = []string{}
	}
	h.Payload[name] = append(h.Payload[name], value...)
}

func (h Headers) Get(name string) ([]string, bool) {
	if v, ok := h.Payload[name]; ok {
		return v, true
	}
	return nil, false
}
