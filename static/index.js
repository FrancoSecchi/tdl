$(document).ready(function () {
  $("#switchMessage").click(function () {
    $(".form-container").toggleClass(
      "animate__fadeInLeft animate__fadeInRight"
    );
    $(".form-container").toggleClass("form-left form-right");

    // Cambia el texto del mensaje
    $("#switchAction").text(
      $("#switchAction").text() === "Registrate"
        ? "Inicia sesión"
        : "Registrate"
    );
  });

  $("#signupForm").submit(function (event) {
    event.preventDefault();
    register();
  });
});

function register() {
  const username = $("#signupUsername").val();
  const password = $("#signupPassword").val();

  console.log($("#signupForm").serialize());
  $.ajax({
    type: "POST",
    url: "/register",
    data: { username: username, password: password },
    success: function (response) {
      if (response.success) {
        
      }
    },
    error: function (error) {
      console.log("Error en la solicitud AJAX:", error);
    },
  });
}
