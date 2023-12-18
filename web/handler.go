package web

import (
    "net/http"
    "encoding/json"
    "fmt"
    "net/url"
    //"strconv"


    "gobusters-chat-app/pkg/chat"
    "golang.org/x/net/websocket"
)

const GLOBAL_CHAT = "GLOBAL_CHAT"
const PRIVATE_CHAT = "PRIVATE_CHAT"

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

	queryParams := r.URL.Query()
      roomID := queryParams.Get("roomID")
	filepath := "chats/global_chat.csv"

	if(roomID != "1") {
		filepath = "chats/" + roomID + "_chat.csv"
	}

	chatMessages, err := chat.GetChatHistoryData(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatMessages)
}

// HandleLogin handles user login.
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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


// NewConnectWsHandler devuelve un http.Handler para la conexión inicial WebSocket.
func NewConnectWsHandler(poolRooms *chat.ChatRoomPool) http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		r := ws.Request()
		params := r.URL.Query()
		HandleConnectWs(ws, poolRooms, params)
	})
}


// HandleConnectWs maneja la solicitud de conexión WebSocket inicial.
func HandleConnectWs(ws *websocket.Conn, poolRooms *chat.ChatRoomPool, params url.Values) {
	typeMessage := params.Get("typeMessage")
	if (typeMessage != PRIVATE_CHAT) {
		globalChat := poolRooms.GetRoomByID(chat.GLOBAL_CHAT_ID)
		globalChat.HandleWs(ws)
	}  else {
		username := params.Get("username")
		userTarget := params.Get("targetUser")
		privateRoom, roomAlreadyExisted := chat.GetOrCreatePrivateRoomBetweenUsers(username, userTarget, ws)
		if (!roomAlreadyExisted) {
			chat.AddRoomToPool(poolRooms, privateRoom);
		}
		privateRoom.HandleWs(ws)
	}
}

// HandleRegister handles user registration.
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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