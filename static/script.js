/////////// Modal de connexion
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


/////////// Likes et dislikes
function toggleLike(button) {
            const isActive = button.classList.contains('active');
            const countSpan = button.querySelector('span:last-child');
            let count = parseInt(countSpan.textContent);
            
            if (isActive) {
                button.classList.remove('active');
                countSpan.textContent = count - 1;
            } else {
                button.classList.add('active');
                countSpan.textContent = count + 1;
                
                // Retirer le dislike si actif
                const post = button.closest('.post');
                const dislikeBtn = post.querySelector('.action-btn:nth-child(2)');
                if (dislikeBtn.classList.contains('active')) {
                    dislikeBtn.classList.remove('active');
                    const dislikeCount = dislikeBtn.querySelector('span:last-child');
                    dislikeCount.textContent = parseInt(dislikeCount.textContent) - 1;
                }
            }
        }

        function toggleDislike(button) {
            const isActive = button.classList.contains('active');
            const countSpan = button.querySelector('span:last-child');
            let count = parseInt(countSpan.textContent);
            
            if (isActive) {
                button.classList.remove('active');
                countSpan.textContent = count - 1;
            } else {
                button.classList.add('active');
                countSpan.textContent = count + 1;
                
                // Retirer le like si actif
                const post = button.closest('.post');
                const likeBtn = post.querySelector('.action-btn:first-child');
                if (likeBtn.classList.contains('active')) {
                    likeBtn.classList.remove('active');
                    const likeCount = likeBtn.querySelector('span:last-child');
                    likeCount.textContent = parseInt(likeCount.textContent) - 1;
                }
            }
        }
        
        

/////////// Modal de réponse
function showResponseBox(button) {
            const modal = document.getElementById('responseModal');
            modal.style.display = 'flex';
            document.getElementById('responseText').focus();
        }

        function closeResponseModal() {
            const modal = document.getElementById('responseModal');
            modal.style.display = 'none';
            document.getElementById('responseText').value = '';
        }

        function submitResponse() {
            const text = document.getElementById('responseText').value.trim();
            if (text) {
                alert('Réponse publiée ! (Simulation)');
                closeResponseModal();
            } else {
                alert('Veuillez saisir une réponse.');
            }
        }

        // Fermer le modal en cliquant à l'extérieur
        document.getElementById('responseModal').addEventListener('click', function(e) {
            if (e.target === this) {
                closeResponseModal();
            }
        });