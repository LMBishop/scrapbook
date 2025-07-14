package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator struct {
	secretKey []byte
	parser    *jwt.Parser
}

func NewAuthenticator(secretKey []byte) *Authenticator {
	parser := jwt.NewParser(jwt.WithIssuer("scrapbook"), jwt.WithExpirationRequired())

	a := &Authenticator{
		secretKey: secretKey,
		parser:    parser,
	}

	return a
}

func (a *Authenticator) NewJwt() (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "scrapbook",
			"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		})

	return t.SignedString(a.secretKey)
}

func (a *Authenticator) VerifyJwt(token string) error {
	_, err := a.parser.Parse(token, func(t *jwt.Token) (interface{}, error) { return a.secretKey, nil })
	return err
}
