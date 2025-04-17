package engine

type User struct {
	Name     string `json:"name"`
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Post struct {
	ID 	 	 int 	`json:"id"`
	UserID   int    `json:"user_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Image    string `json:"image"`
	Category string `json:"category"`
}