package handler

import (
	"fmt"
	"net/http"
)

type PanicHandler struct{}

func NewPanicHandler() *PanicHandler {
	return &PanicHandler{}
}

func (p *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("PANIC!!!")
	fmt.Println("do-panic")
}
