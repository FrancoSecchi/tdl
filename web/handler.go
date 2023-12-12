package web

import (
    "net/http"
    "encoding/json"
    "fmt"

    "gobusters-chat-app/pkg/chat"
)

// UserActionResponse represents the response structure for user actions.
type UserActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    string `json:"user,omitempty"`
}

// HandleIndex serves the index.html page.
func HandleIndex(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "view/index.html")
}


// HandleChat serves the chat.html page.
func HandleChat(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "view/chat.html")
}

// HandleGetChatHistory retrieves the chat history and returns it as JSON.
func HandleGetChatHistory(w http.ResponseWriter, r *http.Request) {
	chatMessages, err := chat.GetChatHistoryData("chats/global_chat.csv")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Devolver el historial como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatMessages)
}

// HandleLogin handles user login.
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Obtener datos del formulario
	username := r.FormValue("username")
	password := r.FormValue("password")

    _, err := chat.Login(username, password)
	if err != nil {
		response := UserActionResponse{
			Success: false,
			Message: err.Error(),
		}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
    response := UserActionResponse{
			Success: true,
			Message: "Se ingresó exitosamente",
		}
    sendJSONResponse(w, response, 200)
    return
}

// HandleRegister handles user registration.
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Obtener datos del formulario
	username := r.FormValue("username")
	password := r.FormValue("password")

    _, err := chat.Register(username, password)
	if err != nil {
		response := UserActionResponse{
			Success: false,
			Message: err.Error(),
		}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
    response := UserActionResponse{
			Success: true,
			Message: "Se registro exitosamente",
		}
    sendJSONResponse(w, response, 200)
    return
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}