package ns

// UserID is a identifier of user.
type UserID string

// String returns the string representation of this user ID.
func (u *UserID) String() string {
	return string(*u)
}

// UserName is a name of user.
type UserName string

// NewUserName generate UserName from string
func NewUserName(n string) UserName {
	return UserName(n)
}

// User is a owner of node.
type User struct {
	ID   UserID   `json:"id"`
	Name UserName `json:"name"`
}

// UserReader defines how to extract users.
type UserReader interface {
	GetByID(id UserID) (*User, error)
}
