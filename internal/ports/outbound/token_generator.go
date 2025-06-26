package outbound

type TokenGenerator interface {
	GenerateToken(userID string) (string, error)
}
