package web

import ("net/http")


func HandleIndex(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "view/index.html")
}


func HandleChat(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "view/chat.html")
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
}
