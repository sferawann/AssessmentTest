package repository

import (
	"errors"

	"github.com/sferawann/go-bank-api/model"
	"github.com/sferawann/go-bank-api/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type NasabahRepository interface {
	Create(newNasabah model.Nasabah) (model.Nasabah, error)
	FindByNIK(nik string) (model.Nasabah, error)
	FindByNoHP(nohp string) (model.Nasabah, error)
}

type nasabahRepository struct {
	db *gorm.DB
}

func (r *nasabahRepository) Create(newNasabah model.Nasabah) (model.Nasabah, error) {
	utils.Log.WithFields(logrus.Fields{
		"nama":   newNasabah.Nama,
		"nik":    newNasabah.NIK,
		"no_hp":  newNasabah.NoHP,
		"action": "create nasabah",
		"layer":  "repository",
	}).Info("Mencoba membuat nasabah baru")

	result := r.db.Create(&newNasabah)
	if result.Error != nil {
		utils.Log.WithError(result.Error).WithFields(logrus.Fields{
			"nama":   newNasabah.Nama,
			"nik":    newNasabah.NIK,
			"no_hp":  newNasabah.NoHP,
			"action": "create nasabah",
			"layer":  "repository",
		}).Error("Gagal membuat nasabah baru")
		return model.Nasabah{}, result.Error

	}
	utils.Log.WithFields(logrus.Fields{
		"id":     newNasabah.ID,
		"nama":   newNasabah.Nama,
		"nik":    newNasabah.NIK,
		"no_hp":  newNasabah.NoHP,
		"action": "create nasabah",
		"layer":  "repository",
	}).Info("Nasabah baru berhasil dibuat")
	return newNasabah, nil
}

func (r *nasabahRepository) FindByNIK(nik string) (model.Nasabah, error) {
	utils.Log.WithFields(logrus.Fields{
		"nik":    nik,
		"action": "FindByNIK",
		"layer":  "repository",
	}).Info("Mencoba mencari nasabah berdasarkan NIK")
	var nasabah model.Nasabah
	err := r.db.Where("nik = ?", nik).First(&nasabah).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Nasabah{}, nil // bukan error, hanya tidak ditemukan
	}
	utils.Log.WithFields(logrus.Fields{
		"id":     nasabah.ID,
		"nama":   nasabah.Nama,
		"nik":    nik,
		"action": "FindByNIK",
		"layer":  "repository",
	}).Info("Berhasil menemukan nasabah berdasarkan NIK")
	return nasabah, err
}

func (r *nasabahRepository) FindByNoHP(nohp string) (model.Nasabah, error) {
	utils.Log.WithFields(logrus.Fields{
		"no_hp":  nohp,
		"action": "FindByNoHP",
		"layer":  "repository",
	}).Info("Mencoba mencari nasabah berdasarkan Nomor HP")
	var nasabah model.Nasabah
	err := r.db.Where("no_hp = ?", nohp).First(&nasabah).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Nasabah{}, nil
	}
	utils.Log.WithFields(logrus.Fields{
		"id":     nasabah.ID,
		"nama":   nasabah.Nama,
		"no_hp":  nohp,
		"action": "FindByNoHP",
		"layer":  "repository",
	}).Info("Berhasil menemukan nasabah berdasarkan Nomor HP")
	return nasabah, err
}

func NewNasabahRepository(db *gorm.DB) NasabahRepository {
	return &nasabahRepository{db}
}
