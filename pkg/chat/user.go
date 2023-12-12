package chat

import (
	"fmt"
	"strconv"

	"golang.org/x/net/websocket"

)

const USERS_FILE = "users.csv"




type RegistrationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User string `json:"user"`

}

type User struct {
    name       string
    password   string
    registered bool
    ws         *websocket.Conn
}

func Login(username string, password string) (*User, error) {

	if user, ok := users[username]; ok && user.password == password {
			return user, nil
	}
	return nil, fmt.Errorf("Credenciales erroneas")
}

// Register realiza la lógica de registro.
func Register(username string, password string) (*User, error) {
	if _, ok := users[username]; !ok {
		user := &User{
			name:       username,
			password:   password,
			registered: true,
		}

		users[username] = user

		if err := appendUsersToCSV([]*User{user}, USERS_FILE); err != nil {
			return nil, fmt.Errorf("Error al escribir el usuario en CSV: %v", err)
		}

		return user, nil
	}

	// Usuario ya registrado
	return nil, fmt.Errorf("Usuario ya registrado")
}

func (u *User) toCSVRecord() []string {
	return []string{u.name, u.password, strconv.FormatBool(u.registered)}
}
