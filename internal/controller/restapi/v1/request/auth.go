package request

// Register is the request body for POST /auth/register.
type Register struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// Login is the request body for POST /auth/login.
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Refresh is the request body for POST /auth/refresh and POST /auth/logout.
type Refresh struct {
	RefreshToken string `json:"refresh_token"`
}
