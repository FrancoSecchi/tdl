let socket;

export function initSocket() {
  socket = new WebSocket("ws://localhost:8080/ws");
}