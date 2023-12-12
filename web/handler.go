package web

import ("net/http")


func HandleIndex(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "chat.html")
}


func HandleChat(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "chat.html")
}

