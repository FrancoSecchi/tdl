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
});
