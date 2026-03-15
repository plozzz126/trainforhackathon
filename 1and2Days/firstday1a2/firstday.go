package first1
import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Работает?"))
	
}


func first1(){
	http.HandleFunc("/test", handler)
	http.ListenAndServe(":8080", nil)
	fmt.Println("прост что бы не ругало")
}