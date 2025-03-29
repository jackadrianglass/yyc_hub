// struct representation of all the data that we're storing
package main;
import (
	"time"
)

// todo: maybe use gorm when we want it
type Event struct {
	ID               int64     // `gorm:"primaryKey;column:id"`
	EventDate        time.Time // `gorm:"column:event_date"`
	EventLocation    string    // `gorm:"column:event_location"`
	EventDescription string    // `gorm:"column:event_description"`
	EventGroupName   string    // `gorm:"column:event_group_name"`
	Dynamic          bool      // `gorm:"column:dynamic"`
	EventID          string    // `gorm:"column:event_id"`
	GroupName        string    // `gorm:"column:group_name"`
}

