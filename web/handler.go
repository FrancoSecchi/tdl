package web

import ("net/http")


func HandleChatRoom(w http.ResponseWriter, r *http.Request) {
    // Servir la interfaz de chat
    http.ServeFile(w, r, "index.html")
}