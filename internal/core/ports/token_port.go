// outbound (secondary adapter)

package ports

type Claims struct {
	UserID string
	Email  string
}

type TokenService interface {
	Generate(userID, email string) (string, error)
	Validate(tokenStr string) (Claims, error)
}
