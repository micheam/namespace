package ns

// UID is a identifier of user.
type UID string

// String returns the string representation of this user ID.
func (u *UID) String() string {
	return string(*u)
}

// UserName is a name of user.
type UserName string

// User is a owner of node.
type User struct {
	ID   UID      `json:"id"`
	Name UserName `json:"name"`
}
