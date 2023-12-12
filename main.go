package main

import (
	"fmt"
	"log"
	"net/http"

	"gobusters-chat-app/pkg/chat"
	"gobusters-chat-app/web"
	"golang.org/x/net/websocket"
)


func main() {
	chatRoom := chat.NewChatRoom()
	http.HandleFunc("/chat", web.HandleChat)
	http.HandleFunc("/", web.HandleIndex)
	http.HandleFunc("/register", web.HandleRegister)
	http.HandleFunc("/login", web.HandleLogin)
	http.Handle("/ws", websocket.Handler(chatRoom.HandleWs))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	fmt.Println("Gobusters Chat Application")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func setRoutes() {
}
