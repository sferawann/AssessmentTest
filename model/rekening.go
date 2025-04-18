package model

import "time"

type Rekening struct {
	ID         int       `gorm:"column:id;primaryKey" json:"id"`
	NasabahID  int       `gorm:"column:nasabah_id" json:"nasabah_id"`
	NoRekening string    `gorm:"column:no_rekening" json:"no_rekening"`
	Saldo      float64   `gorm:"column:saldo" json:"saldo"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`

	Nasabah    Nasabah     `gorm:"foreignKey:NasabahID;references:ID" json:"nasabah"`
	Transaksis []Transaksi `gorm:"foreignKey:RekeningID;references:ID" json:"transaksi"`
}

func (Rekening) TableName() string {
	return "rekening"
}
