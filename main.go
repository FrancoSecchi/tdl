package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

type ChatRoom struct {
	users    map[*User]bool
	messages *os.File
}

func newChatRoom() *ChatRoom {
	messages, err := os.OpenFile("messages.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644) //func OpenFile(name string, flag int, perm FileMode) (*File, error)
	if err != nil {
		fmt.Println("Error newChatRoom:", err)
		return nil
	}

	return &ChatRoom{
		users:    make(map[*User]bool),
		messages: messages,
	}
}

func (r *ChatRoom) handleWs(ws *websocket.Conn) {
	fmt.Println("Se ha conectado un nuevo usuario:", ws.RemoteAddr())

	user := &User{
		ws: ws,
	}

	r.users[user] = true

	info, err := os.Stat("messages.txt")
	if err != nil {
		fmt.Println("Error info:", err)
		return
	}

	if info.Size() > 0 {
		file, err := os.OpenFile("messages.txt", os.O_RDWR, 0644) //func OpenFile(name string, flag int, perm FileMode) (*File, error)
		if err != nil {
			fmt.Println("Error open:", err)
			return
		}

		history := make([]byte, 1024)
		_, err = file.Read(history) //func (f *File) Read(b []byte) (n int, err error)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error read:", err)
				return
			}
		}

		if _, err = user.ws.Write(history); err != nil { //func (ws *Conn) Write(msg []byte) (n int, err error)
			fmt.Println("Error send:", err)
			return
		}

		fmt.Println("Los mensajes anteriores han sido enviados a", ws.RemoteAddr(), ":\n", string(history))
	}

	r.listen(user)
}

func (r *ChatRoom) listen(user *User) {
	data := make([]byte, 1024)
	for {
		n, err := user.ws.Read(data) //func (ws *Conn) Read(msg []byte) (n int, err error)
		if err != nil {
			delete(r.users, user)
			if err == io.EOF {
				fmt.Println("El usuario", user.ws.RemoteAddr(), "se ha desconectado. Quedan", len(r.users), "usuarios conectados.")
			} else {
				fmt.Println("Error ws:", err)
			}
			return
		}

		msg := []byte(user.ws.RemoteAddr().String() + ": " + string(data[:n]) + "\n" + time.Now().String() + "\n")
		fmt.Println(string(msg))
		r.sendToAll(msg)

		if _, err = r.messages.WriteString(string(msg)); err != nil { //func (b *Writer) WriteString(s string) (int, error)
			fmt.Println("Error write:", err)
			return
		}
	}
}

func (r *ChatRoom) sendToAll(msg []byte) {
	for user := range r.users {
		if _, err := user.ws.Write(msg); err != nil { //func (ws *Conn) Write(msg []byte) (n int, err error)
			fmt.Println("Error sendAll:", err)
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
	//type Handler interface { ServeHTTP(ResponseWriter, *Request) }
	//func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request)
	//ServeHTTP implements the http.Handler interface for a WebSocket

	fmt.Println("Gobusters Chat Application")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
