package http

import (
	"strconv"
	"strings"
)

type HeaderItem struct {
	Value []string
}

func (h HeaderItem) String() string {
	return strings.Join(h.Value, ";")
}

func (h HeaderItem) Append(value []string) {
	h.Value = append(h.Value, value...)
}

type Headers struct {
	Connection string
	Host       string
	Referer    string
	UserAgent  string
	Length     int
	Boundary   string
	Raw        string
	Payload    map[string]HeaderItem
}

func ReadHeaders(buffer []byte) Headers {
	headers := Headers{Length: -1, Raw: ""}

	if pos := strings.Index(string(buffer), "\r\n\r\n"); pos > -1 {
		headers.Length = pos
		for _, s := range strings.Split(string(buffer[:pos]), "\r\n") {
			if sn := strings.Index(s, ":"); sn > -1 {
				headers.Raw += s + "\r\n"
				n, v := strings.Trim(s[:sn], " "), strings.Trim(s[sn+1:], " ")

				switch strings.ToLower(n) {
				case "content-type":
				case "content-length":
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
		h.Payload = make(map[string]HeaderItem)
	}

	if _, ok := h.Payload[name]; !ok {
		h.Payload[name] = HeaderItem{Value: value}
	}
	h.Payload[name].Append(value)
}

func (h Headers) Get(name string) (HeaderItem, bool) {
	v, ok := h.Payload[name]
	return v, ok
}

func (h Headers) ContentLength() (length int) {
	length = h.Length
	if val, ok := h.Get("content-length"); ok {
		if n, e := strconv.Atoi(val.String()); e == nil {
			length = n
			return
		}
	}
	return
}

func (h Headers) ContentType() (cType string) {
	if v, ok := h.Get("content-type"); ok {
		cType := v.String()
		if bn := strings.Index(cType, "boundary="); bn > -1 {
			h.Boundary = cType[bn+9:]
			cType = cType[:bn]
		}
	}
	return
}

func (h Headers) String() string {
	return h.Raw
}

type FormValue struct {
	ContentType string
}
