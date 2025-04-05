# Utiliser une image de base pour Go
FROM golang:1.23

# Définir le répertoire de travail dans le conteneur
WORKDIR /app

# Copier les fichiers du projet dans le conteneur
COPY . .

# Télécharger les dépendances
RUN go mod tidy

# Compiler l'application
RUN go build -o main .

# Exposer le port 8080
EXPOSE 8080

# Commande pour exécuter l'application
CMD ["./main"]