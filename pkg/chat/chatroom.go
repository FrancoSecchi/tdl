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

type ChatRoom struct {
	id int
	isPrivate bool
	users    map[*User]bool
}

type ChatMessage struct {
	User    string `json:"user"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

type ChatRoomPool struct {
    rooms map[int]*ChatRoom
}


func NewChatRoomPool() *ChatRoomPool {
    return &ChatRoomPool{
        rooms: make(map[int]*ChatRoom),
    }
}

func (pool *ChatRoomPool) GetRoomByID(roomID int) *ChatRoom {
    return pool.rooms[roomID]
}

func AddRoomToPool(p *ChatRoomPool, r *ChatRoom) {
	p.rooms[r.id] = r
}


// newChatRoom creates and initializes a new ChatRoom.
func NewChatRoom(id int) *ChatRoom {
	rand.Seed(time.Now().UnixNano())

	if (id == 0) {
		id = rand.Intn(100000)
	}

	return &ChatRoom{
		id: id,
		users: make(map[*User]bool),
	}
}

//CreatePrivateChat: Creates a private chat between 2 users
func createPrivateChat(user1, user2 *User) *ChatRoom {
	rand.Seed(time.Now().UnixNano())
	privateChat := &ChatRoom{
		id: rand.Intn(100000),
		isPrivate: true,
		users:     make(map[*User]bool),
	}

	privateChat.users[user1] = true
	privateChat.users[user2] = true

	user1.privateRooms[privateChat.id] = privateChat
	user2.privateRooms[privateChat.id] = privateChat

	return privateChat
}

// handleWs handles WebSocket connections in the chat room.
func (r *ChatRoom) HandleWs(ws *websocket.Conn) {
	params := ws.Request().URL.Query()
	username := params.Get("username")
	fmt.Println("A new user has connected:", username, " - Remote Address:", ws.RemoteAddr())

	newUser := &User{
		name: username,
		ws:   ws,
	}
	r.users[newUser] = true

	message := map[string]int{"roomID": r.id}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	if err := websocket.JSON.Send(ws, jsonMessage); err != nil {
		fmt.Println("Error sending roomID to the client:", err)
		return
	}

	fmt.Println("RoomID sent to the client")

	r.listen(newUser)
}

// listen is a method of the ChatRoom type that listens for incoming messages from a user's WebSocket connection.
// It continuously reads data from the WebSocket and processes it.
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
            delete(r.users, user)
            user.ws.Close()
            return
        }

	  msg := ChatMessage{
		User:    user.name,
		Message: string(data[:n]),
		Time:    time.Now().Format("15:04"),
 	 }
	  jsonData, err := json.Marshal(msg)
        r.sendToAll(jsonData)

	  messageToSave :=  []string {msg.User, msg.Message, msg.Time}


        if _, err = writeChatHistory("global_chat.csv", messageToSave, true); err != nil {
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
