
const gobusters_user = sessionStorage.getItem("gobusters_user");

if (!gobusters_user) {
  window.location.href = "/"
}


$(document).ready(function () {
  initSocket(gobusters_user);
  $("#input-message").on("keydown", function (event) {
    if (event.key === "Enter") {
      sendChatMessage();
    }
  });

  window.localSocket.onmessage = function (event) {
    const data = event.data;
    const parsedMessage = JSON.parse(data);
    console.log(gobusters_user, parsedMessage.user);
    if (parsedMessage.user != gobusters_user)
      appendMessage(
        parsedMessage.user == gobusters_user,
        parsedMessage.message
      );
  };
});

function sendChatMessage() {
  var messageText = $("#input-message").val();
  if (messageText.trim() !== "") {
    appendMessage(true, messageText)
    sendSocketMessage(messageText);
  }
}

function appendMessage(isFromCurrentUser, messageText) {
  var classMessage = isFromCurrentUser ? "message-sent" : "message-received";
  var newMessage = $("<div>").addClass("message " + classMessage);
  var newMessageText = $("<p>").addClass("message-text").text(messageText);
  newMessage.append(newMessageText);
  $(".chat-content").append(newMessage);
  $("#input-message").val("");
}

