package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
)

// Declare a map to store registered users
var users = make(map[string]*User)

type ChatRoom struct {
	users map[*User]bool
}

func newChatRoom() *ChatRoom {
	return &ChatRoom{
		users: make(map[*User]bool),
	}
}

func (r *ChatRoom) handleWs(ws *websocket.Conn) {
	fmt.Println("A new user has connected:", ws.RemoteAddr())

	user, err := loginUser(ws)
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
	defer func() {
		// Remove the user when the WebSocket connection is closed
		delete(r.users, user)
		fmt.Println("User disconnected:", user.name)
	}()

	r.listen(ws)
}

func writeUserToCSV(user *User, filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write user data
	if err := writer.Write(user.toCSVRecord()); err != nil {
		return err
	}

	return nil
}

var registeredUsers []*User

func loginUser(ws *websocket.Conn) (*User, error) {

	// Read action:username:password from the user
	var credentials string
	err := websocket.Message.Receive(ws, &credentials)
	if err != nil {
		return nil, err
	}

	// Split the credentials into action, username, and password
	parts := strings.Split(credentials, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid credentials format")
	}

	action, username, password := parts[0], parts[1], parts[2]

	switch action {
	case "login":
		// In a real-world scenario, you would validate against a database.
		if user, ok := users[username]; ok && user.password == password {
			return user, nil
		}
	case "register":
		// Check if the user is not already registered
		if _, ok := users[username]; !ok {
			user := &User{
				name:       username,
				password:   password,
				registered: true,
				ws:         ws,
			}

			// Store the registered user
			users[username] = user
			registeredUsers = append(registeredUsers, user)

			// Write users to CSV file
			if err := writeUsersToCSV(registeredUsers, "users.csv"); err != nil {
				return nil, err
			}

			return user, nil
		}
		return nil, fmt.Errorf("User already registered")
	}

	return nil, fmt.Errorf("Invalid action")
}

func (r *ChatRoom) listen(ws *websocket.Conn) {
	for {
		var msg string
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(msg)
		r.sendToAll([]byte(msg))
	}
}

func (r *ChatRoom) sendToAll(msg []byte) {
	for user := range r.users {
		_, err := user.ws.Write(msg) // func (ws *Conn) Write(msg []byte) (n int, err error)
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

func writeUsersToCSV(users []*User, filename string) error {
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

type User struct {
	name       string
	password   string
	registered bool
	ws         *websocket.Conn
}

func handleChatRoom(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	chat := newChatRoom()
	http.HandleFunc("/", handleChatRoom)
	http.Handle("/ws", websocket.Handler(chat.handleWs))
	fmt.Println("Gobusters Chat Application")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
