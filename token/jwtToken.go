package token

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/twinj/uuid"
	"neutron0.1/models"
)

var client *redis.Client

// redis connection
func init() {
	dsn := "localhost:6379"

	client = redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println("Redis Client err: ", err)
	}
}

const (
	Key                  string = "iambatman"
	RefreshKey           string = "opaopalolalola1125778432"
	DefaultExpireSeconds int    = 3600 // 1 hour
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUuid string
	Userid     uint64
}

type JwtManager struct {
	Maker
	TokenDetails
}

func NewJwtManager() Maker {
	return &JwtManager{}
}

// Create FindId()-new token with claims
func (m *JwtManager) Create(id int64) (*TokenDetails, error) {
	user, err := models.FindById(id)
	if err != nil {
		fmt.Println("email not found", err.Error())
	}

	// Create the Claims
	td := &TokenDetails{}

	td.AtExpires = time.Now().Add(time.Second * time.Duration(DefaultExpireSeconds)).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = user.Id
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(Key))
	fmt.Println("*td.AccessToken: ", td.AccessToken)
	if err != nil {
		fmt.Println("generate json web token failed - error: ", err)
		return nil, err
	}

	fmt.Printf("* token -%T = %[1]s\n", td.AccessToken)
	fmt.Println("* token will be expired at: ", time.Unix(td.AtExpires, 0))

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_id"] = td.RefreshUuid
	rtClaims["user_id"] = user.Id
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(RefreshKey))
	fmt.Println("*td.RefreshToken: ", td.RefreshToken)
	if err != nil {
		fmt.Println("generate refresh token failed - error: ", err)
		return nil, err
	}
	return td, nil
}

func (m *JwtManager) VerifyToken(tokenString string) (*jwt.Token, error) {
	return VerifyToken(tokenString)
}

func (m *JwtManager) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {

	//extract from request header
	req := ExtractToken(r)

	token, err := VerifyToken(req)
	if err != nil {
		return nil, err
	}

	if !TokenValid(token) {
		fmt.Println("TokenValid:", TokenValid(token))
		return nil, fmt.Errorf("TokenValid: %t", TokenValid(token))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	// fmt.Println("**claims", claims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
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

	return nil, err
}

func (m *JwtManager) RefreshToken(refreshToken string) (tokens map[string]string, err error) {
	refreshSecretKey := []byte(RefreshKey)

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing metgod: %v,", token.Header["alg"])
		}
		return []byte(refreshSecretKey), nil
	})
	if err != nil {
		fmt.Println("Refresh Token expired:", err)
		return nil, err
	}

	if !TokenValid(token) {
		fmt.Println(err)
		return nil, fmt.Errorf("TokenValid: %t", TokenValid(token))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Println("claims: ", claims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_id"].(string)
		fmt.Println("rfid:", refreshUuid)

		if !ok {
			fmt.Println("refreshuuid err:", err)
			return nil, err
		}
		userId, usrerr := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			fmt.Println(err)
			return nil, usrerr
		}
		deleted, delerr := DeleteAuth(refreshUuid)
		if delerr != nil || deleted == 0 {
			fmt.Println("unauthorized:", delerr)
			return nil, delerr
		}

		var m JwtManager
		newToken, newtokErr := m.Create(int64(userId))
		if err != nil {
			fmt.Println("generating token err:", err)
			return nil, newtokErr
		}

		err = CreateAuth(uint(userId), newToken)
		if err != nil {
			fmt.Println(err)
			return
		}

		tokens := map[string]string{
			"access_token":  newToken.AccessToken,
			"refresh_token": newToken.RefreshToken,
		}

		return tokens, nil

	} else {
		fmt.Println("refresh expired")
		return nil, fmt.Errorf("refresh expired")
	}
}

func CreateAuth(userid uint, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

func FetchAuth(authD *AccessDetails) (uint64, error) {
	userid, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func DeleteTokens(authDetails *AccessDetails) error {
	refreshUuid := fmt.Sprintf("%s %d", authDetails.AccessUuid, authDetails.Userid)
	fmt.Println("*refresh sprintf: ", refreshUuid)
	//delete Access Token
	deleteAt, err := client.Del(authDetails.AccessUuid).Result()
	fmt.Println("deleteAt: ", deleteAt)
	if err != nil {
		return err
	}
	//delete Refresh Token
	deleteRt, err := client.Del(refreshUuid).Result()
	fmt.Println("deleteRt: ", deleteRt)
	if err != nil {
		return err
	}
	if deleteAt != 1 && deleteRt != 1 {
		return errors.New("something wrong deleting tokens")
	}
	return nil
}
