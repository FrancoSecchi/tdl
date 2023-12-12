 

function initSocket(username, successCallback) {
  window.socket = new WebSocket("ws://localhost:8080/ws?username=" + username);

  // Evento que se activa cuando la conexión se establece con éxito
  window.socket.onopen = function (event) {
    sessionStorage.setItem("socket", socket);
    console.log("WebSocket connection established:", event);
  };
}

function sendSocketMessage(message) {
  console.log(message)
  window.socket.send(message);
}
