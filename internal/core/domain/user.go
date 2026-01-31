package domain

type User struct {
	ID       uint
	Name     string
	Email    *string
	Phone    string
	Password string
	Timestamp
}

func (User) TableName() string {
	return "auth.users"
}
