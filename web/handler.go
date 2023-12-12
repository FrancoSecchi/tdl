package web

import ("net/http")


func HandleChatRoom(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
}