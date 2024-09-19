package model

// CreateUserParams holds the parameters for creating a new user.
type CreateUserKafkaParams struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int64  `json:"role"`
}
