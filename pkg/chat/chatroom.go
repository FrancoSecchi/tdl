package chat

import (
"fmt"
"io"
"time"
"encoding/json"
"math/rand"

"golang.org/x/net/websocket"
)


// users is a map to store registered users.
var Users = make(map[string]*User)

// ChatRoom represents a chat room.
type ChatRoom struct {
	id int
	users    map[*User]bool
}

type ChatMessage struct {
	User    string `json:"user"`
	Message string `json:"message"`
	Hora    string `json:"hora"`
}


// newChatRoom creates and initializes a new ChatRoom.
func NewChatRoom() *ChatRoom {
	// Inicializar la semilla del generador de números aleatorios
	rand.Seed(time.Now().UnixNano())

	return &ChatRoom{
		id: rand.Intn(100000),
		users:    make(map[*User]bool),
	}
}


// handleWs handles WebSocket connections in the chat room.
func (r *ChatRoom) HandleWs(ws *websocket.Conn) {
	// A new user has connected.
	
	params := ws.Request().URL.Query()
    	username := params.Get("username")
	fmt.Println("A new user has connected:", username, " - Remote Address:", ws.RemoteAddr())

	newUser := &User{
            name: username,
            ws:   ws,
	}
	r.users[newUser] = true
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

	  msg := ChatMessage{
		User:    user.name,
		Message: string(data[:n]),
		Hora:    time.Now().Format("15:04"),
 	 }
	  jsonData, err := json.Marshal(msg)
        r.sendToAll(jsonData)

	  messageToSave :=  []string {msg.User, msg.Message, msg.Hora}


        if _, err = writeChatHistory("global_chat.csv",messageToSave, true); err != nil {
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
