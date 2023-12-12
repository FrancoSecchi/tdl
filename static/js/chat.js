
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
});

function sendChatMessage() {
  var messageText = $("#input-message").val();
  if (messageText.trim() !== "") {
    var newMessage = $("<div>").addClass("message message-sent");
    var newMessageText = $("<p>").addClass("message-text").text(messageText);
    newMessage.append(newMessageText);
    $(".chat-content").append(newMessage);
    $("#input-message").val("");
    sendSocketMessage(messageText);
  }
}

