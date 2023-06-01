package model

// User represents a user in the system.
type User struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Password string `json:"password"`
}

// TableName returns the name of the corresponding database table for the User model.
func (u *User) TableName() string {
	return "user"
}
