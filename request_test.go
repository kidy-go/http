package http

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const rawFile = "./form-data.raw"

func TestRequest(t *testing.T) {
	f, e := os.OpenFile(rawFile, os.O_RDONLY, 0755)
	checkErr(e)

	raw, err := ioutil.ReadAll(f)
	checkErr(err)

	fmt.Println("Read Headers:")
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
				fname = fname[:efn]
				rn := strings.Index(fraw, "\r\n\r\n") + len("\r\n\r\n")

				fmt.Println("FileName: ", []string{fname})
				f, e := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0600)
				defer f.Close()
				checkErr(e)
				fmt.Println([]byte(fraw[rn : len(fraw)-2]))
				f.Write([]byte(fraw[rn : len(fraw)-2]))
				break
			}
		}
	}

	// f2, _ := os.OpenFile("timg.jpg", os.O_RDONLY, 0755)
	// img, _ := ioutil.ReadAll(f2)
	// fmt.Println(img)
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
