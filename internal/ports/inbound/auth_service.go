package inbound

type AuthenService interface {
	Login(email, password string) (string, error)
	Authorize(userID, action string) (bool, error)
}
