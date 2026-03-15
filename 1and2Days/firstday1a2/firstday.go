package main

import (
	"fmt"
	"encoding/json"
	"net/http"
)

type SumReq struct {
	A int `json:"a"`
	B int `json:"b"`
}

type SumRes struct {
	Result int `json:"result"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"Hl": "Hello World",
	}
	json.NewEncoder(w).Encode(resp)
}

func sumHandler(w http.ResponseWriter, r *http.Request) {

	var req SumReq
	json.NewDecoder(r.Body).Decode(&req)

	result := req.A + req.B

	resp := SumRes{
		Result: result,
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/sum", sumHandler)
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}