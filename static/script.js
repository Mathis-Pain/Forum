function openLoginModal() {
        document.getElementById("loginModal").style.display = "block";
      }

      function closeLoginModal() {
        document.getElementById("loginModal").style.display = "none";
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
        closeLoginModal();
      }

      // Fermer le modal en cliquant à l'extérieur
      window.onclick = function (event) {
        const modal = document.getElementById("loginModal");
        if (event.target === modal) {
          closeLoginModal();
        }
      };
