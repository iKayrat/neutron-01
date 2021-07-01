package controllers

import (
	"encoding/json"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"
	"neutron0.1/models"
	"neutron0.1/utils"
)

type UserController struct {
	beego.Controller
}

type AuthResponse struct {
	Message string              `json:"message"`
	User    *models.User        `json:"user"`
	Token   *utils.TokenDetails `json:"token"`
	jwt.StandardClaims
}

type ErrResponse struct {
	Message string `json:"message"`
}

func (c *UserController) ActiveContent(view string) {
	// c.Layout = "basic-layout.tpl"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "header.tpl"
	c.LayoutSections["Footer"] = "footer.tpl"
	c.TplName = view + ".tpl"
}

func (c *UserController) Register() {

	c.ActiveContent("user/register")

	// BodyRequest

	var u models.User
	// fmt.Println("->ctx i ->", string(c.Controller.Ctx.Input.RequestBody))
	err := json.Unmarshal([]byte(c.Ctx.Input.RequestBody), &u)
	if err != nil {
		fmt.Println("unmarshal error: ", err)
		return
	}

	fmt.Println("Input User: ", u)

	id, err := models.CreateNew(u.Name, u.Lastname, u.Username, u.Email, u.Password)
	if err != nil {
		errResponse := ErrResponse{
			Message: err.Error(),
		}
		c.Data["json"] = errResponse
		c.ServeJSON()
		c.StopRun()
	}

	user, err := models.FindById(id)
	if err != nil {
		errResponse := ErrResponse{
			Message: err.Error(),
		}
		c.Data["json"] = errResponse
		c.ServeJSON()
		c.StopRun()
	}

	token, err := utils.GenerateToken(id)
	if err != nil {
		errResponse := ErrResponse{
			Message: "Failed to generate token",
		}
		c.Data["json"] = errResponse
	}

	// req := httplib.Post("http://localhost:8080/api/register")
	// req.Header("Token", token)

	// c.Ctx.SetSecureCookie()("Token", token)

	// req := httplib.Post("http://loclhost:8080/api/register")
	// req.Header("Token", token)

	successRes := AuthResponse{
		Message: "User created successfully",
		User:    user,
		Token:   token,
	}
	c.Data["json"] = successRes
	c.ServeJSON()

}

func (c *UserController) Login() {
	c.ActiveContent("user/login")

	var credentials models.BasicCredentials

	err := json.Unmarshal([]byte(c.Ctx.Input.RequestBody), &credentials)
	if err != nil {
		fmt.Println("unmarshall err: ", err)
		return
	}
	user, err := models.CheckUser(credentials.Email, credentials.Password)
	if err != nil {
		errResponse := ErrResponse{
			Message: err.Error(),
		}
		c.Data["json"] = errResponse
		c.ServeJSON()
		c.StopRun()
	}

	token, err := utils.GenerateToken(int64(user.Id))
	if err != nil {
		errResponse := ErrResponse{
			Message: "Failed to generate token",
		}
		c.Data["json"] = errResponse
	}

	saveErr := utils.CreateAuth(user.Id, token)
	if saveErr != nil {
		errResponse := ErrResponse{
			Message: "Failed to create auth",
		}
		c.Data["json"] = errResponse
	}
	fmt.Println("refreshToken", token.RefreshToken)
	successRes := AuthResponse{
		Message: "User logged in successfully",
		User:    user,
		Token:   token,
	}

	c.Data["json"] = successRes
	c.ServeJSON()
	c.StopRun()
}

func (c *UserController) Logout() {
	auid, err := utils.ExtractTokenMetadata(c.Ctx.Request)
	if err != nil {
		c.Data["json"] = "Unathorized"
		c.ServeJSON()
	}
	err = utils.DeleteTokens(auid)
	if err != nil {
		fmt.Println("delete tokens: ", err)
		c.Data["json"] = "unathorized del"
		c.ServeJSON()
		c.StopRun()
		return
	}
	c.Data["json"] = "Successfully logged out"
	c.ServeJSON()
	c.StopRun()
}

func (c *UserController) Refresh() {
	mapToken := map[string]string{}

	if err := json.Unmarshal([]byte(c.Ctx.Input.RequestBody), &mapToken); err != nil {
		fmt.Println("unmarshal err", err)
		return
	}

	refreshToken := mapToken["refresh_token"]

	tokens, err := utils.RefreshToken(refreshToken)
	if err != nil {
		c.Data["json"] = err
		c.ServeJSON()
	}
	c.Data["json"] = tokens
	c.ServeJSON()
	c.StopRun()
}

// func (c *UserController) Refresh() {
// 	mapToken := map[string]string{}

// 	if err := json.Unmarshal([]byte(c.Ctx.Input.RequestBody), &mapToken); err != nil {
// 		fmt.Println("unmarshal err", err)
// 		return
// 	}
// 	fmt.Println("map Token is: ", mapToken)
// 	refreshSecretKey := utils.RefreshKey
// 	refreshToken := mapToken["refresh_token"]

// 	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing metgod: %v,", token.Header["alg"])
// 		}
// 		return []byte(refreshSecretKey), nil
// 	})
// 	// fmt.Println("token parsed Refresh:", token)
// 	if err != nil {
// 		fmt.Println("Refresh Token expired:", err)
// 		c.Data["json"] = "RefreshToken expired"
// 		c.ServeJSON()
// 		c.StopRun()
// 	}

// 	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
// 		fmt.Println()
// 		c.Data["json"] = "claims err"
// 		c.ServeJSON()
// 		c.StopRun()
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	fmt.Println("claims: ", claims)
// 	if ok && token.Valid {
// 		refreshUuid, ok := claims["refresh_id"].(string)
// 		fmt.Println("rt uid:", refreshUuid)

// 		if !ok {
// 			fmt.Println("refreshuuid err:", err)
// 			c.Data["json"] = "not ok"
// 			c.ServeJSON()
// 			c.StopRun()
// 		}
// 		userId, usrerr := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
// 		if err != nil {
// 			fmt.Println(usrerr)
// 			c.Data["json"] = "Parse err"
// 			c.ServeJSON()
// 			c.StopRun()
// 		}
// 		deleted, delerr := utils.DeleteAuth(refreshUuid)
// 		if delerr != nil || deleted == 0 {
// 			fmt.Println("unauthorized:", delerr)
// 			c.Data["json"] = "DeleteAuth err"
// 			c.ServeJSON()
// 			c.StopRun()
// 		}
// 		newToken, newtokErr := utils.GenerateToken(int64(userId))
// 		if err != nil {
// 			fmt.Println("generating token err:", newtokErr)
// 			c.Data["json"] = "Generate err"
// 			c.ServeJSON()
// 			c.StopRun()
// 		}

// 		err = utils.CreateAuth(uint(userId), newToken)
// 		if err != nil {
// 			fmt.Println(err)
// 			c.Data["json"] = "CreateAuth err"
// 			c.ServeJSON()
// 			c.StopRun()
// 		}

// 		tokens := map[string]string{
// 			"access_token":  newToken.AccessToken,
// 			"refresh_token": newToken.RefreshToken,
// 		}
// 		c.Data["json"] = tokens
// 		c.ServeJSON()
// 		c.StopRun()

// 	} else {
// 		fmt.Println("refresh expired")
// 		c.Data["json"] = err
// 		c.ServeJSON()
// 		c.StopRun()
// 	}
// }
