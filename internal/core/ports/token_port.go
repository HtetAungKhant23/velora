// outbound (secondary adapter)

package ports

type TokenService interface {
	Generate(userID, email string) (string, error)
}
