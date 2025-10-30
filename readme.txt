# 💬 Forum – Projet Web en Go

## 🎯 Objectif du projet

Ce projet consiste à créer un **forum web** en **Go**, permettant la **communication entre utilisateurs**, la **gestion de catégories**, ainsi que les **interactions sociales** (likes/dislikes).  
Le but est de mettre en pratique les fondamentaux du développement web full stack : HTTP, sessions, cookies, bases de données et conteneurisation.

---

## ⚙️ Fonctionnalités principales

- 🧑‍💻 **Authentification des utilisateurs**
  - Inscription avec email, nom d’utilisateur et mot de passe
  - Connexion et gestion de session via cookies
  - Expiration automatique des sessions
  - (Bonus) Chiffrement des mots de passe et UUID pour les sessions

- 💬 **Communication**
  - Création de **posts** et de **commentaires**
  - Association de **catégories** aux posts
  - Affichage des posts et commentaires pour tous les utilisateurs (connectés ou non)

- 👍 **Likes et Dislikes**
  - Un utilisateur connecté peut aimer ou ne pas aimer un post ou un commentaire
  - Le nombre de likes/dislikes est visible par tous

- 🧩 **Filtres**
  - Filtrage des posts par :
    - Catégories  
    - Posts créés  
    - Posts aimés  
  - Les deux derniers filtres sont accessibles uniquement aux utilisateurs connectés

---

## 🗄️ Base de données – SQLite

Le forum utilise **SQLite** pour stocker les données (utilisateurs, posts, commentaires, catégories, etc.).  
Le projet inclut :
- Une initialisation automatique de la base et du schéma
- Des requêtes SQL de type `SELECT`, `CREATE` et `INSERT`

💡 Un **diagramme entité-relation (ERD)** est conseillé pour structurer la base de données.

---

## 🔐 Sécurité et bonnes pratiques

- Gestion complète des erreurs HTTP et techniques  
- Chiffrement possible des mots de passe avec **bcrypt**  
- UUID pour les identifiants de session  
- Nettoyage régulier des sessions expirées  

---

## 🐳 Docker

Le projet est entièrement **conteneurisé avec Docker** :
- Création d’une image serveur Go  
- Gestion simplifiée du déploiement  
- Isolation des dépendances

---

## 🧠 Compétences mises en œuvre

- Langage **Go** (net/http, templates, SQL)  
- **HTML / HTTP / Cookies / Sessions**  
- **SQLite** pour la persistance des données  
- **Docker** pour la conteneurisation  
- Bonnes pratiques de code et gestion des erreurs  

---

## 📦 Technologies et dépendances

- `net/http` – Serveur web Go  
- `database/sql` + `sqlite3` – Base de données  
- `bcrypt` – Chiffrement des mots de passe  
- `uuid` – Gestion unique des sessions  

---

## 🧩 Structure du projet (exemple)

## pour executer dans un docker
docker compose -f docker/compose.yaml up --build -d

## pour lancer un local 
go run ./cmd/server