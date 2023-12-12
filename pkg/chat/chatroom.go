package chat

import (
    "fmt"
    "io"
    "os"
    "time"

    "golang.org/x/net/websocket"
)


// users is a map to store registered users.
var users = make(map[string]*User)

// ChatRoom represents a chat room.
type ChatRoom struct {
	users    map[*User]bool
	messages *os.File
}

// newChatRoom creates and initializes a new ChatRoom.
func NewChatRoom() *ChatRoom {
	// Open or create a file for storing chat messages.
	messages, err := os.OpenFile("messages.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error newChatRoom:", err)
		return nil
	}

	return &ChatRoom{
		users:    make(map[*User]bool),
		messages: messages,
	}
}


// handleWs handles WebSocket connections in the chat room.
func (r *ChatRoom) HandleWs(ws *websocket.Conn) {
	// A new user has connected.
	fmt.Println("A new user has connected:", ws.RemoteAddr())

	// Use a channel to signal the completion of login
	loginDone := make(chan struct{})
	defer close(loginDone)

	var user *User
	var err error

	go func() {
		// Attempt to log in the user.
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
			if err := writeUsersToCSV([]*User{user}, "users.csv"); err != nil {
				fmt.Println("Error writing user to CSV:", err)
			}
		} else {
			fmt.Println("User logged in:", user.name)
		}

		// Add the user to the chat room.
		r.users[user] = true

		// Send chat history to the new user
		info, err := os.Stat("messages.txt")
		if err != nil {
			fmt.Println("Error info:", err)
			return
		}

		if info.Size() > 0 {
			file, err := os.OpenFile("messages.txt", os.O_RDWR, 0644)
			if err != nil {
				fmt.Println("Error open:", err)
				return
			}

			// Read and send chat history to the user.
			history := make([]byte, 1024)
			_, err = file.Read(history)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error read:", err)
					return
				}
			}

			if _, err = user.ws.Write(history); err != nil {
				fmt.Println("Error send:", err)
				return
			}

			fmt.Println("Previous messages have been sent to", ws.RemoteAddr(), ":\n", string(history))
		}

		// Start listening for incoming messages
		r.listen(user)
	}
}



// listen listens for incoming messages from a user.
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

// sendToAll sends a message to all connected users.
func (r *ChatRoom) sendToAll(msg []byte) {
	for user := range r.users {
		_, err := user.ws.Write(msg)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
