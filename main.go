package main

import (
	"fmt"
	"io"
	"net/http"
)

func proxy_hello(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Forwarding request...")
	httpClient := &http.Client{}
	proxyReq, _ := http.NewRequest(req.Method, "http://localhost:1234/hello", nil)
	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	fmt.Println("Request executed: " + string(bodyBytes))
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Service executing request...")
	fmt.Fprintf(w, "Hooray!")
}

func make_proxy() {
	http.HandleFunc("/proxy", proxy_hello)
	http.ListenAndServe(":8080", nil)
}

func main() {
	go make_proxy()
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":1234", nil)
}
