package chat

import (
"fmt"
"io"
"time"
"encoding/json"
"math/rand"
"strconv"

"golang.org/x/net/websocket"
)


const GLOBAL_CHAT_ID = 1

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

func CreatePrivateChat(users []string, ws *websocket.Conn) *ChatRoom {
   rand.Seed(time.Now().UnixNano())

   privateChat := &ChatRoom{
      id:        rand.Intn(100000),
      isPrivate: true,
      users:     make(map[*User]bool),
   }

   for _, username := range users {
      user := Users[username]
	var userWs *websocket.Conn


      if user != nil {
         if !privateChat.users[user] {
		if (user.name == users[0]) {
			userWs = ws
		}
            user.ws = userWs
            privateChat.users[user] = true
            if user.privateRooms == nil {
               user.privateRooms = make(map[int]*ChatRoom)
            }
            user.privateRooms[privateChat.id] = privateChat
         }
      }
   }

   return privateChat
}


// GetOrCreatePrivateRoomBetweenUsers verifica si ya existe una sala de chat privada entre dos usuarios.
// Si existe, devuelve la sala de chat privada; de lo contrario, crea una nueva sala y la devuelve.
func GetOrCreatePrivateRoomBetweenUsers(username1, username2 string, ws *websocket.Conn) (*ChatRoom, bool) {
    // Verificar si ya existe una sala de chat privada entre los usuarios
    user1 := Users[username1]
    user2 := Users[username2]
    existingRoom := findPrivateRoom(user1, user2)
    if existingRoom != nil {
        return existingRoom, true
    }
    fmt.Println("Llego?:", username1)
    return CreatePrivateChat([]string{username1, username2}, ws), false
}

// findPrivateRoom busca una sala de chat privada entre dos usuarios.
// Devuelve la sala de chat privada si la encuentra; de lo contrario, devuelve nil.
func findPrivateRoom(user1, user2 *User) *ChatRoom {
    // Iterar sobre las salas privadas de user1 y verificar si user2 está en alguna de ellas
    for _, room := range user1.privateRooms {
        if room.users[user2] {
		fmt.Println("Room 1:", room)
            return room
        }
    }

    // Iterar sobre las salas privadas de user2 y verificar si user1 está en alguna de ellas
    for _, room := range user2.privateRooms {
        if room.users[user1] {
		fmt.Println("Room 2:", room)
            return room
        }
    }

    // No se encontró ninguna sala privada existente
    return nil
}

// handleWs handles WebSocket connections in the chat room.
func (r *ChatRoom) HandleWs(ws *websocket.Conn) {
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

	params := ws.Request().URL.Query()
	username := params.Get("username")
	if (!r.isPrivate) {	
		fmt.Println("A new user has connected:", username, " - Remote Address:", ws.RemoteAddr())
		newUser := &User{
			name: username,
			ws:   ws,
			privateRooms: make(map[int]*ChatRoom),
		}
		r.users[newUser] = true
		r.listen(newUser)
	} else {
		user, err := GetUser(username)
		if (err != nil) {
			fmt.Println("Error getting the user: ", err)
		}	
		
		if (user.ws == nil) {
			user.ws = ws
		}
		fmt.Println("Usuario: ", user)
		r.users[user] = true
		r.listen(user)
	}

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

	  filename := "global_chat.csv"
	  if (r.isPrivate) {
		filename = strconv.Itoa(r.id) + "_chat.csv"
	  }

        if _, err = writeChatHistory(filename, messageToSave, true); err != nil {
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
