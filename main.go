package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type ChatRoom struct {
	users map[*User]bool
}

func newChatRoom() *ChatRoom {
	return &ChatRoom{
		users: make(map[*User]bool),
	}
}

func (r *ChatRoom) handleWs(ws *websocket.Conn) {
	fmt.Println("Se ha conectado un nuevo usuario:", ws.RemoteAddr())

	user := &User{
		ws: ws,
	}

	r.users[user] = true
	r.listen(user.ws)
}

func (r *ChatRoom) listen(ws *websocket.Conn) {
	data := make([]byte, 1024)
	for {
		n, err := ws.Read(data) //func (ws *Conn) Read(msg []byte) (n int, err error)
		if err != nil {
			fmt.Println(err)
			return
		}

		msg := data[:n]
		fmt.Println(string(msg))
		r.sendToAll(msg)
	}
}

func (r *ChatRoom) sendToAll(msg []byte) {
	for user := range r.users {
		_, err := user.ws.Write(msg) //func (ws *Conn) Write(msg []byte) (n int, err error)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

type User struct {
	name     string
	password string
	ws       *websocket.Conn
}

func handleChatRoom(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	chat := newChatRoom()
	http.HandleFunc("/", handleChatRoom)                 //func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
	http.Handle("/ws", websocket.Handler(chat.handleWs)) //func Handle(pattern string, handler Handler)
	fmt.Println("Gobusters Chat Application")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
