package chat

import (
"fmt"
"io"
"os"
"time"
"math/rand"

"golang.org/x/net/websocket"
)


// users is a map to store registered users.
var users = make(map[string]*User)

// ChatRoom represents a chat room.
type ChatRoom struct {
	id int
	users    map[*User]bool
	messages *os.File
}

// newChatRoom creates and initializes a new ChatRoom.
func NewChatRoom() *ChatRoom {
	// Inicializar la semilla del generador de números aleatorios
	rand.Seed(time.Now().UnixNano())

	// Open or create a file for storing chat messages.
	messages, err := os.OpenFile("messages.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error newChatRoom:", err)
		return nil
	}

	return &ChatRoom{
		id: rand.Intn(100000),
		users:    make(map[*User]bool),
		messages: messages,
	}
}


// handleWs handles WebSocket connections in the chat room.
func (r *ChatRoom) HandleWs(ws *websocket.Conn) {
	// A new user has connected.
	fmt.Println("A new user has connected:", ws.RemoteAddr())

	var user *User = getUserFromWebSocket(ws)
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
	}

	// Start listening for incoming messages
	r.listen(user)

}

/**Función de utilidad para obtener un usuario a partir de la conexión WebSocket
* Each web connection has a unique websocket object.
* A new websocket instance is created each time a client connects to the server.representing that particular connection. 
* As a result, if there are two people connected, each will have their own websocket object. A single Conn.
*/
func getUserFromWebSocket(ws *websocket.Conn) *User {
	for _, user := range users {
		if user.ws == ws {
			return user
		}
	}
	return nil
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
