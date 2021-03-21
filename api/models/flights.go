package models

import (
	"fmt"
	"time"
)

type Flight struct {
	ID          uint      `gorm:"primary_key"`
	Name        string    `json:"name" gorm:"type:varchar(255)"`
	Number      string    `json:"number" gorm:"type:varchar(255)"`
	Scheduled   time.Time `json:"scheduled"`
	Arrival     time.Time `json:"arrival"`
	Departure   time.Time `json:"departure"`
	Destination string    `json:"destination" gorm:"type:varchar(255)"`
	Fare        float64   `json:"fare"`
	Duration    int       `json:"duration"`
}

type Search struct {
	Name        string
	Scheduled   string
	Departure   string
	Destination string
}

// TODO
func (f *Flight) Validate() error {
	if f.Name == "" {
		return fmt.Errorf("Please provide flight name")
	}
	if f.Number == "" {
		return fmt.Errorf("Please provide flight number")
	}
	// if f.Scheduled.IsZero() {
	// 	return fmt.Errorf("Please provide flight scheduled time")
	// }
	// if f.Arrival.IsZero() {
	// 	return fmt.Errorf("Please provide flight arrival time")
	// }
	// if f.Departure.IsZero() {
	// 	return fmt.Errorf("Please provide flight departure time")
	// }
	if f.Destination == "" {
		return fmt.Errorf("Please provide flight destination")
	}
	if f.Fare == 0 {
		return fmt.Errorf("Please provide flight fare")
	}
	if f.Duration == 0 {
		return fmt.Errorf("Please provide esimated flight duration")
	}

	return nil
}

func (f *Flight) Create() error {
	return db.Create(&f).Error
}

func (f *Flight) Delete(id int) error {
	return db.Delete(&Flight{}, id).Error
}

func (f *Flight) Update(id int) error {
	f.ID = uint(id)
	return db.Model(&Flight{}).Updates(f).Error
}

func (f *Flight) Find(s map[string]interface{}) ([]Flight, error) {
	var flights []Flight

	err := db.Where(s).Find(&flights).Error

	return flights, err
}
