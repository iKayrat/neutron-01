package controllers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/astaxie/beego/session"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/gomodule/redigo/redis"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/dgrijalva/jwt-go"
	"neutron0.1/models"
	"neutron0.1/utils"
)

var GlobalSessions *session.Manager

func init() {
	// dsn := "localhost:6379"
	// var err error
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println("redis dial error:", err)
	}
	defer c.Close()

	sessionConfig := &session.ManagerConfig{
		CookieName:      "gosessionid",
		EnableSetCookie: true,
		Gclifetime:      60,
		Maxlifetime:     60,
		Secure:          true,
		CookieLifeTime:  60,
		ProviderConfig:  `127.0.0.1:6379`,
	}

	GlobalSessions, err = session.NewManager("redis", sessionConfig)
	if err != nil {
		fmt.Println("errrrrr sessiossnnfs")
	}
	go GlobalSessions.GC()
}

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
	c.Layout = "basic-layout.tpl"
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

func (c *UserController) Loginsession() {
	c.ActiveContent("user/login")

	sess, err := GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		log.Panic("session: Session start erroro")
	}

	defer sess.SessionRelease(c.Ctx.ResponseWriter)

	// c.Ctx.Request.ParseForm()
	// if c.Ctx.Request.Method == "GET" {
	// 	t, _ := template.ParseFiles("views/user/login.tpl")
	// 	c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/html")
	// 	t.Execute(c.Ctx.ResponseWriter, sess.Get("gosessionid"))

	// } else {
	// 	sess.Set("authtoken", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjA5NjE5MDI0LWMwM2UtNGZkMy05YjZhLWE1MTA5NDA5ZDkzOCIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTYyNjUyMTE3NywidXNlcl9pZCI6M30.7wCKQRCvgN7SHFSCuekMTpwmFIUfDn0ouerCtMf5us0")
	// 	err := c.SetSession("authtoken", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjA5NjE5MDI0LWMwM2UtNGZkMy05YjZhLWE1MTA5NDA5ZDkzOCIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTYyNjUyMTE3NywidXNlcl9pZCI6M30.7wCKQRCvgN7SHFSCuekMTpwmFIUfDn0ouerCtMf5us0")
	// 	if err != nil {
	// 		fmt.Println("set session err", err)
	// 	}
	// 	c.Redirect("/", http.StatusFound)
	// }
	sess.Set("sess", "package cookie set")

	packageCookie := sess.Get("sess")
	fmt.Println("packageCookie: ", packageCookie)

	c.SetSecureCookie("cookie", "auth", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjcyODA5NTYsInJlZnJlc2hfaWQiOiIwOTcwMTY5Yi00MTRkLTQ0ZWMtOGU5Ni1kNzA1MzUwMWM4MGEiLCJ1c2VyX2lkIjozfQ.SuAo0SAdvZB33R1csFchgStZFJdRe1ljYIzJUStygvc")
	getsecurecookie, ok := c.GetSecureCookie("cookie", "auth")
	if !ok {
		fmt.Println("cookie is false")
	}
	fmt.Println("getsecurecookie is: ", getsecurecookie)

	cookie := c.Ctx.Input.Cookie("gosessionid")
	cookieContext := c.Ctx.Input.Cookie("auth")
	cookieget := sess.Get(cookie)
	getsess := c.GetSession("gosessionid")
	fmt.Println("cookie is: ", cookie)
	fmt.Println("cookieContext is: ", cookieContext)
	fmt.Println("cookieget is: ", cookieget)
	fmt.Println("getsess is: ", getsess)

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
