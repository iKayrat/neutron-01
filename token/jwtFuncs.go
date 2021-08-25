package token

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func TokenValid(tokenstring *jwt.Token) bool {

	if _, ok := tokenstring.Claims.(jwt.Claims); !ok && !tokenstring.Valid {
		return false
	}
	return true
}

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	// fmt.Println("**bearer Token:", bearerToken)
	// fmt.Println("**strArr:", strArr)
	if len(strArr) == 2 {
		// fmt.Println("**strArr[1]:", strArr[1])
		return strArr[1]
	}
	return ""
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	BitfielldCheckToken(tokenString)

	hmacSamplesecret := []byte(Key)
	// fmt.Println("**tokenstring verify:", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// fmt.Println(hmacSamplesecret)
		return hmacSamplesecret, nil
	})
	// fmt.Println("***token: ", token)
	if err != nil {
		// fmt.Println("**verify1 err: ", err)
		return nil, err
	}
	// fmt.Println("**tokenParse verify1:", token)

	return token, nil
}

func BitfielldCheckToken(tokenstring string) {
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		return []byte(Key), nil
	})
	if token.Valid {
		fmt.Println("***Token is Valid***")
	} else if v, ok := err.(*jwt.ValidationError); ok {
		if v.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if v.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			fmt.Println("Timing is everything")
		} else {
			fmt.Println("Couldn't handle this token!")
		}
	} else {
		fmt.Println("Couldn't handle this token!")
	}
}

func getUuid(token *jwt.Token) (*AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return nil, fmt.Errorf("token not valid")
	}

	accessUuid, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, fmt.Errorf("claims uuid is emty")
	}

	userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
	if err != nil {
		return nil, err
	}

	return &AccessDetails{
		AccessUuid: accessUuid,
		Userid:     userId,
	}, nil

}
