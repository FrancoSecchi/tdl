package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
)

const USESR_FILE = "users.csv"

type User struct {
	name       string
	password   string
	registered bool
	ws         *websocket.Conn
}

// Declare a map to store registered users
var users = make(map[string]*User)

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
	fmt.Println("A new user has connected:", ws.RemoteAddr())

	// Use a channel to signal the completion of login
	loginDone := make(chan struct{})
	defer close(loginDone)

	var user *User
	var err error

	go func() {
		user, err = loginUser(ws)
		loginDone <- struct{}{}
	}()

	select {
	case <-loginDone:
		// Login completed
		if err != nil {
			fmt.Println("Login/Register error:", err)
			return
		}

		if user.registered {
			fmt.Println("User registered:", user.name)

			// Write user to CSV file
			if err := writeUserToCSV(user, "users.csv"); err != nil {
				fmt.Println("Error writing user to CSV:", err)
			}
		} else {
			fmt.Println("User logged in:", user.name)
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

// Add a method to convert a user to a CSV record
func (u *User) toCSVRecord() []string {
	return []string{u.name, u.password, strconv.FormatBool(u.registered)}
}

func writeUserToCSV(users []*User, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Name", "Password", "Registered"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write user data
	for _, user := range users {
		if err := writer.Write(user.toCSVRecord()); err != nil {
			return err
		}
	}

	return nil
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
