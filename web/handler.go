package web

import (
    "net/http"
	"encoding/json"
    "fmt"

    "gobusters-chat-app/pkg/chat"
)


func HandleIndex(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "view/index.html")
}


func HandleChat(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "view/chat.html")
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {

}


type RegistrationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    string `json:"user,omitempty"`
}

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
		response := RegistrationResponse{
			Success: false,
			Message: err.Error(),
		}
		sendJSONResponse(w, response, http.StatusBadRequest)
		return
	}
    response := RegistrationResponse{
			Success: true,
			Message: "Se registro exitosamente",
		}
    sendJSONResponse(w, response, 200)
    return
	//http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}