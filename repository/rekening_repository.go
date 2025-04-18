package repository

import (
	"errors"

	"github.com/sferawann/go-bank-api/model"
	"github.com/sferawann/go-bank-api/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RekeningRepository interface {
	Create(newRekening model.Rekening) (model.Rekening, error)
	FindByNasabahID(nasabahID int) (model.Rekening, error)
	FindByNoREK(noREK string) (model.Rekening, error)
	UpdateSaldo(UpdateRekening model.Rekening) (model.Rekening, error)
}

type rekeningRepository struct {
	db *gorm.DB
}

func (r *rekeningRepository) Create(newRekening model.Rekening) (model.Rekening, error) {
	utils.Log.WithFields(logrus.Fields{
		"nasabah_id":  newRekening.NasabahID,
		"no_rekening": newRekening.NoRekening,
		"saldo":       newRekening.Saldo,
		"action":      "create rekening",
		"layer":       "repository",
	}).Info("Mencoba membuat rekening baru")
	result := r.db.Create(&newRekening)
	if result.Error != nil {
		utils.Log.WithFields(logrus.Fields{
			"nasabah_id":  newRekening.NasabahID,
			"no_rekening": newRekening.NoRekening,
			"saldo":       newRekening.Saldo,
			"error":       result.Error,
			"action":      "create rekening",
			"layer":       "repository",
		}).Error("Gagal membuat rekening")
		return model.Rekening{}, result.Error
	}
	utils.Log.WithFields(logrus.Fields{
		"id":          newRekening.ID,
		"nasabah_id":  newRekening.NasabahID,
		"no_rekening": newRekening.NoRekening,
		"saldo":       newRekening.Saldo,
		"action":      "create rekening",
		"layer":       "repository",
	}).Info("Berhasil membuat rekening baru")
	return newRekening, nil
}

func (r *rekeningRepository) FindByNasabahID(nasabahID int) (model.Rekening, error) {
	utils.Log.WithFields(logrus.Fields{
		"nasabah_id": nasabahID,
		"action":     "FindByNasabahID",
		"layer":      "repository",
	}).Info("Mencari rekening berdasarkan NasabahID")
	var rekening model.Rekening
	err := r.db.Where("nasabah_id = ?", nasabahID).First(&rekening).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.Log.WithFields(logrus.Fields{
			"nasabah_id": nasabahID,
			"action":     "FindByNasabahID",
			"layer":      "repository",
		}).Warn("Rekening tidak ditemukan")
		return model.Rekening{}, nil
	}
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"nasabah_id": nasabahID,
			"error":      err,
			"action":     "FindByNasabahID",
			"layer":      "repository",
		}).Error("Gagal mencari rekening")
		return model.Rekening{}, err
	}
	utils.Log.WithFields(logrus.Fields{
		"nasabah_id":  nasabahID,
		"no_rekening": rekening.NoRekening,
		"action":      "FindByNasabahID",
		"layer":       "repository",
	}).Info("Rekening ditemukan")
	return rekening, nil
}

func (r *rekeningRepository) UpdateSaldo(UpdateRekening model.Rekening) (model.Rekening, error) {
	result := r.db.Save(&UpdateRekening)
	if result.Error != nil {
		return model.Rekening{}, result.Error
	}
	return UpdateRekening, nil
}

func (r *rekeningRepository) FindByNoREK(noREK string) (model.Rekening, error) {
	utils.Log.WithFields(logrus.Fields{
		"no_rekening": noREK,
		"action":      "FindByNoREK",
		"layer":       "repository",
	}).Info("Mencari rekening berdasarkan noREK")
	var rekening model.Rekening
	err := r.db.Where("no_rekening = ?", noREK).First(&rekening).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.Log.WithFields(logrus.Fields{
			"no_rekening": noREK,
			"action":      "FindByNoREK",
			"layer":       "repository",
		}).Warn("No rekening tidak ditemukan")
		return model.Rekening{}, nil
	}
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"no_rekening": noREK,
			"error":       err,
			"action":      "FindByNoREK",
			"layer":       "repository",
		}).Error("Gagal mencari rekening")
		return model.Rekening{}, err
	}
	utils.Log.WithFields(logrus.Fields{
		"no_rekening": noREK,
		"action":      "FindByNoREK",
		"layer":       "repository",
	}).Info("No rekening ditemukan")
	return rekening, err
}

func NewRekeningRepository(db *gorm.DB) RekeningRepository {
	return &rekeningRepository{db}
}
