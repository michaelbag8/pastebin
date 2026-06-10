package main


import(
	"net/http"
)

func main() {
	http.HandleFunc("POST   /pastes     ", pasteHandler)     
http.HandleFunc("GET    /pastes     ", pasteHandler)     
http.HandleFunc("GET    /pastes/{id}", pasteHandler)     
http.HandleFunc("DELETE /pastes/{id}", pasteHandler) 
}