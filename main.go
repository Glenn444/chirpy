package main

import (
	"io"
	"log"
	"net/http"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Server struct{
	Address string
	Handler Handler
}
func Hello(w http.ResponseWriter,req *http.Request){
	io.WriteString(w, "This is my website")
}
func main() {
	
	mux := http.NewServeMux()
	//rh := http.RedirectHandler("tobitresearchconsulting.com",307)
	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/",fileServer)

	log.Print("Listening...")

	http.ListenAndServe(":8080", mux)

}