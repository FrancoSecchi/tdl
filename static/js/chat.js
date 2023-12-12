const gobusters_user = sessionStorage.getItem("gobusters_user");

if (!gobusters_user) {
  window.location.href = "/";
}

$(document).ready(function () {
  initSocket(gobusters_user);
  setChatHistory();
  $("#input-message").on("keydown", function (event) {
    if (event.key === "Enter") {
      sendChatMessage();
    }
  });

  window.localSocket.onmessage = function (event) {
    const data = event.data;
    const parsedMessage = JSON.parse(data);
    console.log(parsedMessage);
    if (parsedMessage.user != gobusters_user)
      appendMessage(
        parsedMessage.user == gobusters_user,
        parsedMessage.user,
        parsedMessage.message,
        parsedMessage.time
      );
  };
});

function setChatHistory() {
  $.ajax({
    url: "/getChatHistory",
    method: "GET",
    success: function (history) {
      history.forEach(function (message) {
        appendMessage(
          message.user === gobusters_user,
          message.user,
          message.message,
          message.time
        );
      });
    },
    error: function (error) {
      console.error("Error al obtener el historial de chat:", error);
    },
  });
}

function getTime() {
  const now = new Date();
  const hours = now.getHours().toString().padStart(2, "0");
  const minutes = now.getMinutes().toString().padStart(2, "0");
  return `${hours}:${minutes}`;
}

function sendChatMessage() {
  var messageText = $("#input-message").val();
  if (messageText.trim() !== "") {
    appendMessage(true, gobusters_user, messageText, getTime());
    sendSocketMessage(messageText);
  }
}

function appendMessage(isFromCurrentUser, user, messageText, time) {
  var classMessage = isFromCurrentUser ? "message-sent" : "message-received";
  var newMessage = $("<div>").addClass("message " + classMessage);
  var newMessageText = $("<p>").addClass("message-text").text(messageText);
  var userText = $("<p>").addClass("username").text(user);
  var timeText = $("<p>").addClass("time").text(time);
  console.log(user, messageText, time);
  if (!isFromCurrentUser) {
    newMessage.append(userText);
  }
  newMessage.append(newMessageText);
  newMessage.append(timeText);
  $(".chat-content").append(newMessage);
  $("#input-message").val("");
  var chatContent = $(".chat-content");
  chatContent.scrollTop(chatContent[0].scrollHeight);
}
