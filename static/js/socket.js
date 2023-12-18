function initSocket(username, typeMessage) {
  window.localSocket = new WebSocket(
    "ws://localhost:8080/ws?username=" +
      username +
      "&typeMessage=" +
      typeMessage
  );

  window.localSocket.addEventListener("open", function (event) {
    console.log("WebSocket connection established:", event);
  });

  window.localSocket.addEventListener("message", function (event) {

    const jsonData = isBase64(event.data)
      ? JSON.parse(atob(event.data))
      : JSON.parse(event.data);

    if (jsonData.roomID !== undefined) {
      const roomID = jsonData.roomID;
      sessionStorage.setItem("gobusters_current_room_id", roomID);
    } 
  });
}

const isBase64 = (str) => {
  try {
    return btoa(atob(str)) === str;
  } catch (err) {
    return false;
  }
};

function initPrivateRoomSocket(username, username2) {
  
}

function sendSocketMessage(message) {
  window.localSocket.send(message);
}
