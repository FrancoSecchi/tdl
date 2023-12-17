$(document).ready(function () {
  $("#switchMessage").click(function () {
    $(".form-container").toggleClass(
      "animate__fadeInLeft animate__fadeInRight"
    );
    $(".form-container").toggleClass("form-left form-right");

    $("#switchAction").text(
      $("#switchAction").text() === "¿No tenes cuenta? Registrate"
        ? "¿Ya tenes cuenta? Inicia sesión"
        : "¿No tenes cuenta? Registrate"
    );
  });

  $("#signupForm").submit(function (event) {
    event.preventDefault();
    register();
  });

  $("#loginForm").submit(function (event) {
    event.preventDefault();
    login();
  });
});

function login() {
  const username = $("#loginUsername").val();
  const password = $("#loginPassword").val();

  $.ajax({
    type: "POST",
    url: "/login",
    data: { username: username, password: password },
    success: function (response) {
      if (response.success) {
          sessionStorage.setItem("gobusters_user", username)

          window.location.href = "/chat";
      }
    },
    error: function (error) {
      alert("Credenciales invalidas");
    },
  });
}

function register() {
  const username = $("#signupUsername").val();
  const password = $("#signupPassword").val();

  $.ajax({
    type: "POST",
    url: "/register",
    data: { username: username, password: password },
    success: function (response) {
      if (response.success) {
        sessionStorage.setItem("gobusters_user", username);
        window.location.href = "/chat";
      }
    },
    error: function (error) {
      alert("Problemas al intentar registrarse");
    },
  });
}
