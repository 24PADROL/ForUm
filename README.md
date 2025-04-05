# 🛸 Forum

Ce projet est une application de forum développée en **Go** avec une base de données **MySQL**, orchestrée via **Docker Compose**. Il permet aux utilisateurs de créer des sujets de discussion et de publier des messages dans un environnement web simple et fonctionnel.

---

## 📦 Prérequis

Avant de lancer le serveur, assurez-vous d’avoir installé les éléments suivants sur votre machine :

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

---

## ✨ Installation

### 1. Clonez le projet

```bash
git clone https://github.com/24PADROL/ForUm.git
cd ForUm
```

### 2. Lancez le serveur

Construisez et démarrez les conteneurs avec Docker Compose :

```bash
docker-compose up --build
```

> 🐳 Cette commande va :
> - Démarrer un conteneur MySQL avec la base `forum`
> - Exécuter le script `database.sql` pour créer les tables
> - Lancer l'application Go qui se connecte à la base de données

---

## 🌐 Accéder à l'application

Une fois les conteneurs lancés avec succès, ouvrez votre navigateur et allez à l'adresse suivante :

👉 [http://localhost:8080](http://localhost:8080)

---

## 📁 Structure du projet

```bash
ForUm/
├── db/
│   └── database.sql         # Script SQL pour créer les tables du forum
├── server/                  # Backend server
├── web/                     # Frontend web
│
├── Dockerfile               # Image de l'application Go
├── docker-compose.yml       # Configuration Docker (MySQL + App)
├── main.go                  # Code principal de l'application
└── README.md                # Ce fichier
```

## 🛡️ Licence

Ce projet est sous licence **MIT**.

---

## 👨‍💻 Auteur

- **24PADROL**
🔗 [GitHub](https://github.com/24PADROL)
- **CasualElf34**
🔗 [GitHub](https://github.com/CasualElf34)
- **yasmine200**
🔗 [GitHub](https://github.com/yasmine200)
- **mkbyx**
🔗 [GitHub](https://github.com/mkbyx)
