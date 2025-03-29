package main

import (
	"database/sql"
	_ "fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
)

type TestForm struct {
	Email   string `form:"email" json:"email" binding:"required"`
	Message string `form:"message" json:"message" binding:"required"`
}

type EventA struct {
	Title       string
	Description string
	// Date        string // todo: Make this into a date time object
	// Location    string // todo: Strongly typed address?
}

type EventDb struct {
	db *sql.DB
}

func OpenConnection() (*EventDb, error) {
	const file string = "test.db"

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	const create string = `
		CREATE TABLE IF NOT EXISTS events (
		id INTEGER NOT NULL PRIMARY KEY,
		title TEXT
		description TEXT
		);`

	if _, err := db.Exec(create); err != nil {
		return nil, err
	}

	return &EventDb{
		db: db,
	}, nil
}

func (self *EventDb) Insert(event EventA) (int, error) {
	res, err := self.db.Exec("INSERT INTO activities VALUES(NULL,?,?);", event.Title, event.Description)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func main() {
	// db, error := OpenConnection()
	// if error != nil {
	// 	fmt.Println(error)
	// 	return
	// }

	// id, error := db.Insert(EventA { Title: "awesome sauce", Description: "Saucy awesomeness" })
	// if error != nil {
	// 	fmt.Println(error)
	// 	return
	// }
	// fmt.Println(id)
	// fmt.Println("success!")
	


	app := fiber.New(fiber.Config{
		Views: html.New("./templates", ".tmpl"),
	})

	app.Get("/", func(ctx fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{})
	})

	app.Post("/test-form", func(ctx fiber.Ctx) error {
		var form TestForm
		if err := ctx.Bind().Body(&form); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString("Bad Request")
		}

		return ctx.Render("test-form-rsp", fiber.Map{
			"email":   form.Email,
			"message": form.Message,
		})
	})

	app.Listen(":3000")
}
