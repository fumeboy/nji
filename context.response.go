package nji

import "net/http"

type Response struct {
	Writer http.ResponseWriter
}

func (r *Response) String(status int, data string) (err error){
	r.Writer.WriteHeader(status)
	_, err = r.Writer.Write([]byte(data))
	return
}