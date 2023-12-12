let socket

function initSocket(username, successCallback) {
  socket = new WebSocket("ws://localhost:8080/ws?username=" + username);

  // Evento que se activa cuando la conexión se establece con éxito
  socket.onopen = function (event) {
    console.log("WebSocket connection established:", event);
    // Llamada al callback de éxito si se proporciona
    if (successCallback) {
      successCallback();
    }
  };
}