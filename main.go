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
	http.HandleFunc("/", web.HandleChatRoom)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/ws", websocket.Handler(chatRoom.HandleWs))
	fmt.Println("Gobusters Chat Application")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
