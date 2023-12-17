package chat

import (
	"fmt"
	"strconv"

	"golang.org/x/net/websocket"

)

const USERS_FILE = "users.csv"


type User struct {
    name       string
    password   string
    registered bool 
    privateRooms map[int]*ChatRoom
    ws *websocket.Conn
}

func GetUser(username string) (*User, error) {
	allUsers, _ := getUsersFromCSV(USERS_FILE); 

	if user, ok := allUsers[username]; ok {
		Users[username] = user
		return user, nil
	}
	return nil, fmt.Errorf("Usuario incorrecto")

}

func Login(username string, password string) (*User, error) {
		
	allUsers, _ := getUsersFromCSV(USERS_FILE); 


	if user, ok := allUsers[username]; ok && user.password == password {
		Users[username] = user
		return user, nil
	}

	return nil, fmt.Errorf("Credenciales erroneas")
}

// Register realiza la lógica de registro.
func Register(username string, password string) (*User, error) {
	if _, ok := Users[username]; !ok {
		user := &User{
			name:       username,
			password:   password,
			registered: true,
		}

		Users[username] = user

		if err := appendUsersToCSV([]*User{user}, USERS_FILE); err != nil {
			return nil, fmt.Errorf("Error al escribir el usuario en CSV: %v", err)
		}

		return Users[username], nil
	}

	// Usuario ya registrado
	return nil, fmt.Errorf("Usuario ya registrado")
}

func (u *User) toCSVRecord() []string {
	return []string{u.name, u.password, strconv.FormatBool(u.registered)}
}
