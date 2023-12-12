package chat

import (
	"strings"
	"fmt"
	"strconv"
	
	"golang.org/x/net/websocket"

)


type User struct {
    name       string
    password   string
    registered bool
    ws         *websocket.Conn
}



// loginUser reads user credentials from the WebSocket.
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

			// Write users to CSV file
			if err := writeUsersToCSV([]*User{user}, "users.csv"); err != nil {
				return nil, err
			}

			return user, nil
		}
		return nil, fmt.Errorf("User already registered")
	}

	return nil, fmt.Errorf("Invalid action")
}

func (u *User) toCSVRecord() []string {
	return []string{u.name, u.password, strconv.FormatBool(u.registered)}
}
