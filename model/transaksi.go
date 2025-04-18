package model

import "time"

type Transaksi struct {
	ID             int       `gorm:"column:id;primaryKey" json:"id"`
	RekeningID     int       `gorm:"column:rekening_id" json:"rekening_id"`
	Nominal        float64   `gorm:"column:nominal" json:"nominal"`
	JenisTransaksi string    `gorm:"column:jenis_transaksi" json:"jenis_transaksi"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`

	Rekening Rekening `gorm:"foreignKey:RekeningID;references:ID" json:"rekening"`
}

func (Transaksi) TableName() string {
	return "transaksi"
}
