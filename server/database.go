package engine

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error

    // Modifier ces valeurs selon ta configuration MySQL
    user := "root"
    password := "8426"
    host := "127.0.0.1"
    port := "3306"
    database := "forum"

    // Chaîne de connexion MySQL
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database)

    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Erreur de connexion à MySQL :", err)
    }

    // Vérifier la connexion
    if err := DB.Ping(); err != nil {
        log.Fatal("Impossible de pinger MySQL :", err)
    }

    createTables()
}

// Création des tables
func createTables() {
    usersTable := `CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(100) UNIQUE NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        role ENUM('user', 'moderator', 'admin') DEFAULT 'user'
    );`

    postsTable := `CREATE TABLE IF NOT EXISTS posts (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT,
        title VARCHAR(255) NOT NULL,
        content TEXT NOT NULL,
        image VARCHAR(255),
        category VARCHAR(100),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
    );`

    commentsTable := `CREATE TABLE IF NOT EXISTS comments (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT,
        post_id INT,
        content TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
    );`

    likesTable := `CREATE TABLE IF NOT EXISTS likes (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT,
        post_id INT,
        is_like TINYINT(1),
        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
    );`

    notificationsTable := `CREATE TABLE IF NOT EXISTS notifications (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT,
        message TEXT NOT NULL,
        is_read TINYINT(1) DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
    );`

    moderationTable := `CREATE TABLE IF NOT EXISTS moderation (
        id INT AUTO_INCREMENT PRIMARY KEY,
        moderator_id INT,
        post_id INT,
        action VARCHAR(100) NOT NULL,
        reason TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(moderator_id) REFERENCES users(id) ON DELETE SET NULL,
        FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
    );`

    tables := []string{usersTable, postsTable, commentsTable, likesTable, notificationsTable, moderationTable}

    for _, table := range tables {
        _, err := DB.Exec(table)
        if err != nil {
            log.Fatal("Erreur lors de la création des tables :", err)
        }
    }
    log.Println("Base de données MySQL et tables créées avec succès !")
}
