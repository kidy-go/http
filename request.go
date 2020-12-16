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
	Boundary      string
	Length        int
	Raw           string
	Payload       map[string][]string
}

func ReadHeaders(buffer []byte) Headers {
	headers := Headers{ContentLength: -1, Raw: ""}

	if pos := strings.Index(string(buffer), "\r\n\r\n"); pos > -1 {
		headers.Length = pos
		headers.ContentLength = pos
		for _, s := range strings.Split(string(buffer[:pos]), "\r\n") {
			if sn := strings.Index(s, ":"); sn > -1 {
				headers.Raw += s + "\r\n"
				n, v := strings.Trim(s[:sn], " "), strings.Trim(s[sn+1:], " ")

				switch strings.ToLower(n) {
				case "content-type":
					if bn := strings.Index(v, "boundary="); bn > -1 {
						headers.Boundary = v[bn+9:]
					}
					headers.ContentType = strings.Split(v, ";")[0]
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

				val := []string{}
				for _, xv := range strings.Split(v, ";") {
					val = append(val, strings.Trim(xv, " "))
				}
				headers.Set(strings.ToLower(n), val)
			}
		}
	}
	return headers
}

func (h *Headers) Set(name string, value []string) {
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

type FormValue struct {
	ContentType string
}
