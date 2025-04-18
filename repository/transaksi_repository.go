package repository

import (
	"errors"

	"github.com/sferawann/go-bank-api/model"
	"github.com/sferawann/go-bank-api/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransaksiRepository interface {
	Tarik(newTarik model.Transaksi) (model.Transaksi, error)
	Tabung(newTabung model.Transaksi) (model.Transaksi, error)
	FindByRekeningID(rekeningID int) (model.Transaksi, error)
}

type transaksiRepository struct {
	db *gorm.DB
}

func (r *transaksiRepository) Tarik(newTarik model.Transaksi) (model.Transaksi, error) {
	utils.Log.WithFields(logrus.Fields{
		"no_rekening": newTarik.Rekening.NoRekening,
		"nominal":     newTarik.Nominal,
		"action":      "tarik",
		"layer":       "repository",
	}).Info("Mencoba membuat transaksi tarik baru")
	result := r.db.Create(&newTarik)
	if result.Error != nil {
		utils.Log.WithError(result.Error).WithFields(logrus.Fields{
			"no_rekening": newTarik.Rekening.NoRekening,
			"nominal":     newTarik.Nominal,
			"action":      "tarik",
			"layer":       "repository",
		}).Error("Gagal membuat transaksi tarik baru")
		return model.Transaksi{}, result.Error
	}
	utils.Log.WithFields(logrus.Fields{
		"id":              newTarik.ID,
		"no_rekening":     newTarik.Rekening.NoRekening,
		"nominal":         newTarik.Nominal,
		"jenis_transaksi": newTarik.JenisTransaksi,
		"action":          "tarik",
		"layer":           "repository",
	}).Info("berhasil membuat transaksi tarik baru")
	return newTarik, nil
}

func (r *transaksiRepository) Tabung(newTabung model.Transaksi) (model.Transaksi, error) {
	utils.Log.WithFields(logrus.Fields{
		"no_rekening": newTabung.Rekening.NoRekening,
		"nominal":     newTabung.Nominal,
		"action":      "tabung",
		"layer":       "repository",
	}).Info("Mencoba membuat transaksi tabung baru")
	result := r.db.Create(&newTabung)
	if result.Error != nil {
		utils.Log.WithError(result.Error).WithFields(logrus.Fields{
			"no_rekening": newTabung.Rekening.NoRekening,
			"nominal":     newTabung.Nominal,
			"action":      "tabung",
			"layer":       "repository",
		}).Error("Gagal membuat transaksi tabung baru")
		return model.Transaksi{}, result.Error
	}
	utils.Log.WithFields(logrus.Fields{
		"id":              newTabung.ID,
		"no_rekening":     newTabung.Rekening.NoRekening,
		"nominal":         newTabung.Nominal,
		"jenis_transaksi": newTabung.JenisTransaksi,
		"action":          "tabung",
		"layer":           "repository",
	}).Info("berhasil membuat transaksi tabung baru")
	return newTabung, nil
}

func (r *transaksiRepository) FindByRekeningID(rekeningID int) (model.Transaksi, error) {
	utils.Log.WithFields(logrus.Fields{
		"rekening_id": rekeningID,
		"action":      "FindByRekeningID",
		"layer":       "repository",
	}).Info("Mencari rekening berdasarkan rekeningID")
	var transaksi model.Transaksi
	err := r.db.Preload("Rekening").Where("rekening_id = ?", rekeningID).Last(&transaksi).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.Log.WithFields(logrus.Fields{
			"rekening_id": rekeningID,
			"action":      "FindByRekeningID",
			"layer":       "repository",
		}).Warn("Rekening tidak ditemukan berdasarkan rekeningID")
		return model.Transaksi{}, nil
	}
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"rekening_id": rekeningID,
			"error":       err,
			"action":      "FindByRekeningID",
			"layer":       "repository",
		}).Error("Gagal mencari rekening berdasarkan rekeningID")
		return model.Transaksi{}, err
	}
	utils.Log.WithFields(logrus.Fields{
		"id":          transaksi.ID,
		"rekening_id": rekeningID,
		"nominal":     transaksi.Nominal,
		"jenis":       transaksi.JenisTransaksi,
		"action":      "FindByRekeningID",
		"layer":       "repository",
	}).Info("Transaksi ditemukan")
	return transaksi, nil
}

func NewTransaksiRepository(db *gorm.DB) TransaksiRepository {
	return &transaksiRepository{db}
}
