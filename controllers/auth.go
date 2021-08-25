package controllers

import (
	"encoding/json"
	"fmt"

	"neutron0.1/models"
)

func (c *UserController) Auth() {
	var u models.User
	// fmt.Println("->ctx i ->", string(c.Controller.Ctx.Input.RequestBody))
	err := json.Unmarshal([]byte(c.Ctx.Input.RequestBody), &u)
	if err != nil {
		fmt.Println("unmarshal error: ", err)
		return
	}

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

	token, err := newjwt.Create(int64(user.Id))
	if err != nil {
		errResponse := ErrResponse{
			Message: err.Error(),
		}
		c.Data["json"] = errResponse
		c.ServeJSON()
		c.StopRun()
	}

	// saveErr := util.CreateAuth(user.Id, token)
	// if saveErr != nil {
	// 	errResponse := ErrResponse{
	// 		Message: "Failed to create auth",
	// 	}
	// 	c.Data["json"] = errResponse
	// }

	successRes := AuthResponse{
		Message: "User created successfully",
		User:    user,
		Token:   token,
	}
	c.Data["json"] = successRes
	c.ServeJSON()
	c.StopRun()

}
