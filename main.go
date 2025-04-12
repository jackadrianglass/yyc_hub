package main

import (
	_ "database/sql"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/gofiber/fiber/v3"
	_ "github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
)

// DB Models

type User struct {
	gorm.Model
	Name       string
	Email      string `gorm:"uniqueIndex"`
	Age        int
	Profile    Profile
	Posts      []Post
}

type Group struct {
	gorm.Model
	Name string
	MeetupGroupName *string
}

type Profile struct {
	gorm.Model
	UserID uint
	Bio    string
}

type Post struct {
	gorm.Model
	UserID  uint
	Title   string
	Content string
}

type Event struct {
	gorm.Model
	Title       string
	Description string
	Date        time.Time
	Location    string
}

type MeetupEvent struct {
	gorm.Model
	// todo: Link it to events table
	MeetupId string `gorm:"uniqueIndex"`
	Dynamic bool
	GroupName string
}

// front end stuff

type TestForm struct {
	Email   string `form:"email" json:"email" binding:"required"`
	Message string `form:"message" json:"message" binding:"required"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{}, &Profile{}, &Post{})

	db.Create(&User{
		Name:    "John",
		Email:   "john@example.com",
		Age:     30,
		Profile: Profile{Bio: "Software developer"},
	})

	var user User
	db.First(&user, 1)
	fmt.Println(user)
	db.First(&user, "email = ?", "john@example.com")
	fmt.Println(user)

	db.Preload("Profile").Preload("Posts").Find(&user)

	db.Model(&user).Update("Name", "John Doe")
	db.Model(&user).Updates(User{Name: "John Doe", Age: 31})

	db.Delete(&user)

	// app := fiber.New(fiber.Config{
	// 	Views: html.New("./templates", ".tmpl"),
	// })

	// app.Get("/", func(ctx fiber.Ctx) error {
	// 	return ctx.Render("index", fiber.Map{})
	// })

	// app.Post("/test-form", func(ctx fiber.Ctx) error {
	// 	var form TestForm
	// 	if err := ctx.Bind().Body(&form); err != nil {
	// 		return ctx.Status(fiber.StatusBadRequest).SendString("Bad Request")
	// 	}

	// 	return ctx.Render("test-form-rsp", fiber.Map{
	// 		"email":   form.Email,
	// 		"message": form.Message,
	// 	})
	// })

	// app.Listen(":3000")
}
