package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	for i := 0; i < 100; i++ {
		go func() {
			var buf = make([]byte, 200)
			var i int
			for j := 0; j < 100_000; j++ {
				r, err := http.Post("http://127.0.0.1:8003/json", "application/json", strings.NewReader(`{"A":"Hello ", "B": "World!"}`))
				if err != nil {
					fmt.Println(err)
					return
				}
				if i, err = r.Body.Read(buf); err != io.EOF {
					fmt.Println(err)
					return
				}
				if string(buf[:i]) != "Hello World!" {
					fmt.Println(123)
					return
				}
				r.Body.Close()
			}
		}()
	}
	time.Sleep(19*time.Second)
}
