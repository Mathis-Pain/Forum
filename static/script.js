 function openRegisterModal() {
        document.getElementById("registerModal").style.display = "block";
      }

      function closeRegisterModal() {
        document.getElementById("registerModal").style.display = "none";
      }

      function handleRegister(event) {
        event.preventDefault();
        const password = document.getElementById("registerPassword").value;
        const confirmPassword =
          document.getElementById("confirmPassword").value;

        if (password !== confirmPassword) {
          alert("Les mots de passe ne correspondent pas");
          return;
        }

        alert("Inscription réussie !");
        closeRegisterModal();
      }

      function showLoginModal() {
        alert("Modal de connexion à implémenter");
      }

      // Fermer le modal en cliquant à l'extérieur
      window.onclick = function (event) {
        const modal = document.getElementById("registerModal");
        if (event.target === modal) {
          closeRegisterModal();
        }
      };