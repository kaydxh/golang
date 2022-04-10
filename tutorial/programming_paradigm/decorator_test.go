package tutorial

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

//http 中间件采用修饰器编程规范
type HttpHandlerDecorator func(http.HandlerFunc) http.HandlerFunc

func Handler(h http.HandlerFunc, decors ...HttpHandlerDecorator) http.HandlerFunc {
	for i := range decors {
		d := decors[len(decors)-1-i] // iterate in reverse
		h = d(h)
	}
	return h
}

func WithServerHeader(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithServerHeader()")
		w.Header().Set("Server", "HelloServer v0.0.1")
		h(w, r)
	}
}

func WithAuthCookie(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithAuthCookie()")
		cookie := &http.Cookie{Name: "Auth", Value: "Pass", Path: "/"}
		http.SetCookie(w, cookie)
		h(w, r)
	}
}

func WithBasicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithBasicAuth()")
		//cookie, _ := r.Cookie("Auth")

		/*
			if err != nil || cookie.Value != "Pass" {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		*/
		h(w, r)
	}
}

func WithDebugLog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithDebugLog")
		r.ParseForm()
		/*
					log.Println(r.Form)
			//		log.Println("path", r.URL.Path)
			//		log.Println("scheme", r.URL.Scheme)
			//		log.Println(r.Form["url_long"])
					for k, v := range r.Form {
						log.Println("key:", k)
						log.Println("val:", strings.Join(v, ""))
					}
		*/
		h(w, r)
	}
}

//handler function
func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved Request %s from %s\n", r.URL.Path, r.RemoteAddr)
	fmt.Fprintf(w, "Hello, World! "+r.URL.Path)
}

func TestHttpMiddleware(t *testing.T) {
	http.HandleFunc("/v4/hello", Handler(hello,
		WithServerHeader, WithBasicAuth, WithDebugLog))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
