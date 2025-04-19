// struct representation of all the data that we're storing
package main

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name    string
	Email   string `gorm:"uniqueIndex"`
	Age     int
	Profile Profile
	Posts   []Post
}

type Group struct {
	gorm.Model
	Name            string
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
	MeetupId  string `gorm:"uniqueIndex"`
	Dynamic   bool
	GroupName string
}

func ExampleDbStuff () {
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
}
