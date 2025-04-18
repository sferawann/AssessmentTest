package model

import (
	"time"
)

type Nasabah struct {
	ID        int       `gorm:"column:id;primaryKey" json:"id"`
	Nama      string    `gorm:"column:nama" json:"nama"`
	NIK       string    `gorm:"column:nik" json:"nik"`
	NoHP      string    `gorm:"column:no_hp" json:"no_hp"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Nasabah) TableName() string {
	return "nasabah"
}
