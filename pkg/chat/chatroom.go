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
var Users = make(map[string]*User)

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
	
	params := ws.Request().URL.Query()
    	username := params.Get("username")

	newUser := &User{
            name: username,
            ws:   ws,
	}

	// Add the user to the chat room.
	r.users[newUser] = true

	// Start listening for incoming messages
	r.listen(newUser)
}

// listen escucha los mensajes entrantes de un usuario.
func (r *ChatRoom) listen(user *User) {
    data := make([]byte, 1024)
    for {
        n, err := user.ws.Read(data)
        if err != nil {
            if err == io.EOF {
                fmt.Println("El usuario", user.name, "cerró la conexión.")
            } else {
                fmt.Println("Error ws:", err)
            }

            // Elimina al usuario y cierra la conexión
            delete(r.users, user)
            user.ws.Close()
            return
        }

        msg := []byte(user.ws.RemoteAddr().String() + ": " + string(data[:n]) + "\n" + time.Now().String() + "\n")
        fmt.Println(string(msg))
        r.sendToAll(msg)

        if _, err = r.messages.WriteString(string(msg)); err != nil {
            fmt.Println("Error write:", err)
            // Puedes decidir si quieres cerrar la conexión del usuario aquí en caso de un error de escritura.
            // user.ws.Close()
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
