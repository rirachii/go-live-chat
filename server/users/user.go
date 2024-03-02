package users




type User struct {
	UserID string;
	UserDisplayName string;
	UserEmail string;
}


func (user User) ID() string {
	return user.UserID
}


func (user User) DisplayName() string {
	return user.UserDisplayName
}

func (user User) Email() string {
	return user.UserEmail
}
