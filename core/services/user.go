package services

type UserService interface {
	Signin(username, password string) (result string, err error)
}
