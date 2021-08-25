package token

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// Maker interface is for managing tokens
type Maker interface {
	Create(id int64) (*TokenDetails, error)
	ExtractTokenMetadata(r *http.Request) (*AccessDetails, error)
	VerifyToken(tokenstring string) (*jwt.Token, error)
	RefreshToken(refreshToken string) (tokens map[string]string, err error)
}
