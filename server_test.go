// server_test.go kee > 2020/12/15

package http

import (
	"fmt"
	"testing"
)

func TestListen(t *testing.T) {
	//e := ListenAndServeTLS(":443", "test.pem", "test.key")
	e := ListenAndServe(":80")
	fmt.Println(e)
}
