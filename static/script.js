document.addEventListener("DOMContentLoaded", () => {
  // A helper function to manage modals
  function setupModal(
    modalId,
    openBtnClass,
    closeBtnClass,
    outsideClickable = true
  ) {
    const modal = document.getElementById(modalId);
    if (!modal) return;

    const openButtons = document.querySelectorAll(openBtnClass);
    const closeButtons = document.querySelectorAll(closeBtnClass);

    // Open the modal
    openButtons.forEach((btn) => {
      btn.addEventListener("click", (e) => {
        e.preventDefault(); // Prevent default link/button behavior if needed
        modal.style.display = "flex";
      });
    });

    // Close the modal
    closeButtons.forEach((btn) => {
      btn.addEventListener("click", () => {
        modal.style.display = "none";
      });
    });

    // Close the modal when clicking outside
    if (outsideClickable) {
      window.addEventListener("click", (event) => {
        if (event.target === modal) {
          modal.style.display = "none";
        }
      });
    }
  }

  // --- Modal de connexion ---
  setupModal("loginModal", ".log-me-in", ".close-btn"); // Assuming you have buttons with these classes

  // Handle registration form submission
  const registerForm = document.getElementById("registerForm");
  if (registerForm) {
    registerForm.addEventListener("submit", (e) => {
      const password = document.getElementById("registerPassword").value;
      const confirmPassword = document.getElementById("confirmPassword").value;

      if (password !== confirmPassword) {
        alert("Les mots de passe ne correspondent pas.");
        return;
      }

      alert("Inscription réussie !");
      document.getElementById("loginModal").style.display = "none";
    });
  }

  // --- Modal de réponse ---
  setupModal("responseModal", ".response-btn", ".close-btn"); // Assuming buttons with these classes

  // Handle response submission
  const responseForm = document.getElementById("responseForm");
  if (responseForm) {
    responseForm.addEventListener("submit", (e) => {
      e.preventDefault();
      const responseText = document.getElementById("responseText").value.trim();

      if (responseText) {
        alert("Réponse publiée ! (Simulation)");
        document.getElementById("responseModal").style.display = "none";
        document.getElementById("responseText").value = "";
      } else {
        alert("Veuillez saisir une réponse.");
      }
    });
  }

  // --- Modal de confirmation de suppression ---
  const confirmationModal = document.getElementById("confirmationModal");
  if (confirmationModal) {
    let formToSubmit = null;

    // Open modal on delete button click
    document.querySelectorAll(".delete-btn").forEach((button) => {
      button.addEventListener("click", (event) => {
        event.preventDefault(); // Prevent default form submission
        formToSubmit = event.target.closest("form");
        confirmationModal.style.display = "flex";
      });
    });

    // Confirm delete
    const confirmBtn = document.getElementById("confirmDelete");
    if (confirmBtn) {
      confirmBtn.addEventListener("click", () => {
        confirmationModal.style.display = "none";
        if (formToSubmit) {
          formToSubmit.submit();
        }
      });
    }

    // Cancel delete
    const cancelBtn = document.getElementById("cancelDelete");
    if (cancelBtn) {
      cancelBtn.addEventListener("click", () => {
        confirmationModal.style.display = "none";
        formToSubmit = null;
      });
    }
  }

  // --- Edit/Validate button functionality ---
  document.querySelectorAll(".modifier-btn").forEach((button) => {
    button.addEventListener("click", (event) => {
      const row = event.target.closest("tr");
      if (!row) return;

      // Toggle visibility of fields and buttons
      row
        .querySelectorAll(".user-edit")
        .forEach((input) => input.classList.remove("is-hidden"));
      row
        .querySelectorAll(".userfield")
        .forEach((span) => span.classList.add("is-hidden"));
      row.querySelector(".validate-btn")?.classList.remove("is-hidden");
      row.querySelector(".cancel-btn")?.classList.remove("is-hidden");
      row.querySelector(".modifier-btn")?.classList.add("is-hidden");
    });
  });

  document.querySelectorAll(".cancel-btn").forEach((button) => {
    button.addEventListener("click", (event) => {
      const row = event.target.closest("tr");
      if (!row) return;

      // Toggle visibility of fields and buttons back
      row
        .querySelectorAll(".user-edit")
        .forEach((input) => input.classList.add("is-hidden"));
      row
        .querySelectorAll(".userfield")
        .forEach((span) => span.classList.remove("is-hidden"));
      row.querySelector(".validate-btn")?.classList.add("is-hidden");
      row.querySelector(".cancel-btn")?.classList.add("is-hidden");
      row.querySelector(".modifier-btn")?.classList.remove("is-hidden");
    });
  });

  document.querySelectorAll(".validate-btn").forEach((button) => {
    button.addEventListener("click", (event) => {
      const row = event.target.closest("tr");
      if (!row) return;

      // Toggle visibility of fields and buttons back
      row
        .querySelectorAll(".user-edit")
        .forEach((input) => input.classList.add("is-hidden"));
      row
        .querySelectorAll(".userfield")
        .forEach((span) => span.classList.remove("is-hidden"));
      row.querySelector(".validate-btn")?.classList.add("is-hidden");
      row.querySelector(".cancel-btn")?.classList.add("is-hidden");
      row.querySelector(".modifier-btn")?.classList.remove("is-hidden");
    });
  });
});
