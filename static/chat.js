$(document).ready(function () {
  $("#input-message").on("keydown", function (event) {
    if (event.key === "Enter") {
      sendMessage();
    }
  });
});

function sendMessage() {
  var messageText = $("#input-message").val();
  if (messageText.trim() !== "") {
    var newMessage = $("<div>").addClass("message message-sent");
    var newMessageText = $("<p>").addClass("message-text").text(messageText);
    newMessage.append(newMessageText);
    $(".chat-content").append(newMessage);
    $("#input-message").val("");
  }
}
