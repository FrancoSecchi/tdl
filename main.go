package main

import (
	"fmt"
	"log"
	"net/http"

	"gobusters-chat-app/pkg/chat"
	"gobusters-chat-app/web"
)




func main() {
	poolRooms := chat.NewChatRoomPool()
	chatRoom := chat.NewChatRoom(1)
	chat.AddRoomToPool(poolRooms, chatRoom);
	
	http.HandleFunc("/chat", web.HandleChat)
	http.HandleFunc("/", web.HandleIndex)
	http.HandleFunc("/register", web.HandleRegister)
	http.HandleFunc("/login", web.HandleLogin)
	http.HandleFunc("/getChatHistory", web.HandleGetChatHistory)
	http.Handle("/ws", web.NewConnectWsHandler(poolRooms))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	fmt.Println("Gobusters Chat Application")
	log.Fatal(http.ListenAndServe(":8080", nil))
}