package engine

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Post struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Message   string `json:"message"`
    CreatedAt string `json:"created_at"` 
}

type PostMessage struct {
    Username string `json:"username"`
    Message  string `json:"message"`
}

type Message struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Content   string    `json:"content"`
    CreatedAt string    `json:"created_at"`
}

