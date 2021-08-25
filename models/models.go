package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"neutron0.1/util"

	_ "github.com/lib/pq"
)

type User struct {
	Id       uint      `json:"id"`
	Name     string    `json:"name"`
	Lastname string    `json:"lastname"`
	Username string    `json:"username"`
	Email    string    `json:"email" orm:"index,size(191)"`
	Password string    `json:"password"`
	Created  time.Time `json:"created_on" orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `json:"updated_on" orm:"auto_now;type(datetime)"`
}

type InputUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BasicCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func init() {
	_ = orm.RegisterDriver("postgres", orm.DRPostgres)
	_ = orm.RegisterDataBase("default", "postgres", "dbname=neutron0.1 host=localhost user=postgres password=kaak port=5432 sslmode=disable")
	// _ = orm.RegisterDataBase("clientToken", "postgres", "dbname=neutron0.1 host=localhost user=postgres password=kaak port=5432 sslmode=disable")

	orm.RegisterModel(new(User))
}

func CreateNew(name, lastname, username, email, password string) (uid int64, err error) {
	o := orm.NewOrm()
	var user User
	user.Name = name
	user.Lastname = lastname
	user.Username = username
	user.Email = email

	hash, err := util.HashPassword(password)
	if err != nil {
		return -1, err
	}
	user.Password = hash

	uid, insertErr := o.Insert(&user)
	if insertErr != nil {
		return -1, errors.New("failed to insert user to database")
	}
	return uid, nil
}

func FindById(id int64) (user *User, err error) {
	o := orm.NewOrm()
	userId := User{Id: uint(id)}
	e := o.Read(&userId)

	if e == orm.ErrNoRows {
		return nil, errors.New("user not found")
	} else if e == nil {
		return &userId, nil
	} else {
		return nil, errors.New("unknown error occured")
	}

}

func FindByEmail(email string) (user *User, err error) {
	o := orm.NewOrm()
	u := User{Email: email}
	e := o.Read(&u, "Email")

	if e == orm.ErrNoRows {
		return nil, errors.New("user not found")
	} else if e == nil {
		return &u, nil
	} else {
		return nil, errors.New("unknown error occured")
	}
}

func CheckUser(email, password string) (user *User, err error) {
	u, err := FindByEmail(email)
	fmt.Println(u)
	if err == nil {
		if ok := util.CheckPasswordAndHash(password, u.Password); !ok {
			return nil, errors.New("email and password doesn't match")
		}
		return u, nil

	} else {
		return nil, err
	}

}
