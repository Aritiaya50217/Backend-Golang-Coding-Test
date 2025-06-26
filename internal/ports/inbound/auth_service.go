package inbound

type AuthenService interface {
	Login(email, password string) (string, error)
}
