# Forum

Ce projet est une application de forum développée en Go avec une base de données MySQL.

## Prérequis

Avant de lancer le serveur, assurez-vous d'avoir les éléments suivants installés sur votre machine :

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Installation

1. Clonez ce dépôt sur votre machine locale :
   ```bash
   git clone https://github.com/24PADROL/ForUm.git
   cd ForUm

2. Lancer le serveur
Construisez et démarrez les conteneurs avec Docker Compose :
    ```bash	
    docker-compose up --build

Une fois les conteneurs démarrés, accédez à l'application dans votre navigateur à l'adresse suivante :
    ```bash
    http://localhost:8080
