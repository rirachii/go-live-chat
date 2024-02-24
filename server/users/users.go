package users




type User struct {
	userID string;
	userDisplayName string;
	userEmail string;
}


func (user User) ID() string {
	return user.userID
}


func (user User) DisplayName() string {
	return user.userDisplayName
}

func (user User) Email() string {
	return user.userEmail
}
