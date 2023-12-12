 

function initSocket(username, successCallback) {
  window.localSocket = new WebSocket("ws://localhost:8080/ws?username=" + username);

  // Evento que se activa cuando la conexión se establece con éxito
  window.localSocket.onopen = function (event) {
    console.log("WebSocket connection established:", event);
  };
}

function sendSocketMessage(message) {
  window.localSocket.send(message);
}