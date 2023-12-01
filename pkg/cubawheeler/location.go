package cubawheeler

import "gorm.io/gorm"

type Location struct {
	gorm.Model
	ID   string  `json:"id" gorm:"primaryKey;varchar(36);not null"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
	User string  `json:"-"`
}

func (l *Location) BeforeSave(*gorm.DB) error {
	if l.ID == "" {
		l.ID = NewID().String()
	}
	return nil
}
