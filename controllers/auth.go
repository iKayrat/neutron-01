package controllers

// import (
// 	"fmt"

// 	_ "github.com/jackc/pgx"
// 	"github.com/jinzhu/gorm"
// 	"github.com/qor/auth"
// )

// type Authentication struct {
// 	Id       uint `gorm:"primarykey, autoIncrement"`
// 	Name     string
// 	Email    string `gorm:"unique"`
// 	Password string
// }

// var (
// 	dsn = "host=localhost user=postgres password=kaak dbname=auth port=5432 sslmode=disable"

// 	db, _ = gorm.Open("postgres", dsn)

// 	// Initialize Auth with configuration
// 	Auth = auth.New(&auth.Config{
// 		DB: db,
// 	})
// )

// func init() {
// 	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	// if err != nil {
// 	// 	fmt.Println("gorm: ", err)
// 	// }

// 	// _, err := gorm.Open(postgres.New(postgres.Config{
// 	// 	DriverName: "postgresql",
// 	// 	DSN:        dsn,
// 	// }))

// 	// fmt.Println(err)

// 	db.DB().Ping()
// 	db.AutoMigrate(&Authentication{})

// 	a := Authentication{Name: "Ciril", Email: "cer", Password: "123"}
// 	result := db.Create(&a)

// 	fmt.Println(result.Error)
// 	fmt.Println(result.RowsAffected)
// }
