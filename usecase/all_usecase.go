package usecase

import (
	"errors"
	"math"

	"github.com/sferawann/go-bank-api/model"
	"github.com/sferawann/go-bank-api/repository"
	"github.com/sferawann/go-bank-api/utils"
	"github.com/sirupsen/logrus"
)

type AllUsecase interface {
	Create(newNasabah model.Nasabah) (model.Nasabah, error)
	FindByNasabahID(nasabahID int) (model.Rekening, error)
	FindByNoREK(noREK string) (model.Rekening, error)
	FindByRekeningID(rekeningID int) (model.Transaksi, error)
	Tarik(newTarik model.Transaksi) (model.Transaksi, error)
	Tabung(newTabung model.Transaksi) (model.Transaksi, error)
}

type allUsecase struct {
	NasabahRepository   repository.NasabahRepository
	RekeningRepository  repository.RekeningRepository
	TransaksiRepository repository.TransaksiRepository
}

func (u *allUsecase) Create(NewNasabah model.Nasabah) (model.Nasabah, error) {
	utils.Log.WithFields(logrus.Fields{
		"nama":   NewNasabah.Nama,
		"nik":    NewNasabah.NIK,
		"no_hp":  NewNasabah.NoHP,
		"action": "create",
		"layer":  "allUsecase",
	}).Info("menerima permintaan pembuatan nasabah")

	utils.Log.WithFields(logrus.Fields{
		"nik":    NewNasabah.NIK,
		"action": "FindByNIK",
		"layer":  "allUsecase",
	}).Info("Memeriksa ketersediaan nik")
	findNIK, err := u.NasabahRepository.FindByNIK(NewNasabah.NIK)
	if err != nil {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"nik":    NewNasabah.NIK,
			"action": "FindByNIK",
			"layer":  "allUsecase",
		}).Error("Gagal mencari nasabah berdasarkan nik")
		return model.Nasabah{}, err
	}
	if findNIK.ID != 0 {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"nik":    NewNasabah.NIK,
			"action": "validasi",
			"layer":  "allUsecase",
		}).Error("nik sudah digunakan")
		return model.Nasabah{}, errors.New("nik sudah digunakan")
	}

	utils.Log.WithFields(logrus.Fields{
		"nik":    NewNasabah.NoHP,
		"action": "FindByNoHP",
		"layer":  "allUsecase",
	}).Info("Memeriksa ketersediaan No HP")
	findNOHP, err := u.NasabahRepository.FindByNoHP(NewNasabah.NoHP)
	if err != nil {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"nik":    NewNasabah.NoHP,
			"action": "FindByNoHP",
			"layer":  "allUsecase",
		}).Error("Gagal mencari nasabah berdasarkan No HP")
		return model.Nasabah{}, err
	}
	if findNOHP.ID != 0 {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"nik":    NewNasabah.NoHP,
			"action": "validasi",
			"layer":  "allUsecase",
		}).Error("no hp sudah digunakan")
		return model.Nasabah{}, errors.New("no hp sudah digunakan")
	}

	utils.Log.WithFields(logrus.Fields{
		"nama":   NewNasabah.Nama,
		"nik":    NewNasabah.NIK,
		"no_hp":  NewNasabah.NoHP,
		"action": "create",
		"layer":  "allUsecase",
	}).Info("Membuat nasabah baru melalui repository")
	createdNasabah, err := u.NasabahRepository.Create(NewNasabah)
	if err != nil {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"nama":   NewNasabah.Nama,
			"nik":    NewNasabah.NIK,
			"no_hp":  NewNasabah.NoHP,
			"action": "create",
			"layer":  "allUsecase",
		}).Error("Gagal membuat nasabah melalui repository")
		return model.Nasabah{}, err
	}

	noRek := utils.GenerateNoRek()
	rekening := model.Rekening{
		NasabahID:  createdNasabah.ID,
		NoRekening: noRek,
	}

	utils.Log.WithFields(logrus.Fields{
		"nasabah_id":  createdNasabah.ID,
		"no_rekening": noRek,
		"action":      "create",
		"layer":       "allUsecase",
	}).Info("Membuat rekening untuk nasabah")
	_, err = u.RekeningRepository.Create(rekening)
	if err != nil {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"nasabah_id":  createdNasabah.ID,
			"no_rekening": noRek,
			"action":      "create",
			"layer":       "allUsecase",
		}).Error("Gagal membuat rekening")
		return model.Nasabah{}, err
	}

	utils.Log.WithFields(logrus.Fields{
		"id":          createdNasabah.ID,
		"nama":        createdNasabah.Nama,
		"nik":         createdNasabah.NIK,
		"no_hp":       createdNasabah.NoHP,
		"no_rekening": noRek,
		"action":      "create",
		"layer":       "allUsecase",
	}).Info("Pembuatan nasabah dan rekening berhasil")
	return createdNasabah, nil
}

func (u *allUsecase) FindByNasabahID(nasabahID int) (model.Rekening, error) {
	utils.Log.WithFields(logrus.Fields{
		"nasabah_id": nasabahID,
		"action":     "FindByNasabahID",
		"layer":      "allUsecase",
	}).Info("Mencari rekening berdasarkan Nasabah ID")
	FindNasabahID, err := u.RekeningRepository.FindByNasabahID(nasabahID)
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"nasabah_id": nasabahID,
			"error":      err,
			"action":     "FindByNasabahID",
			"layer":      "allUsecase",
		}).Error("Gagal mencari rekening berdasarkan Nasabah ID")
		return model.Rekening{}, err
	}
	utils.Log.WithFields(logrus.Fields{
		"nasabah_id": nasabahID,
		"action":     "FindByNasabahID",
		"layer":      "allUsecase",
	}).Info("Berhasil menemukan rekening berdasarkan Nasabah ID")
	return FindNasabahID, nil
}

func (u *allUsecase) FindByRekeningID(rekeningID int) (model.Transaksi, error) {
	utils.Log.WithFields(logrus.Fields{
		"rekening_id": rekeningID,
		"action":      "FindByRekeningID",
		"layer":       "allUsecase",
	}).Info("Mencari transaksi berdasarkan Rekening ID")
	FindNasabahID, err := u.TransaksiRepository.FindByRekeningID(rekeningID)
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"rekening_id": rekeningID,
			"error":       err,
			"action":      "FindByRekeningID",
			"layer":       "allUsecase",
		}).Error("Gagal mencari transaksi berdasarkan Rekening ID")
		return model.Transaksi{}, err
	}
	utils.Log.WithFields(logrus.Fields{
		"rekening_id": rekeningID,
		"action":      "FindByRekeningID",
		"layer":       "allUsecase",
	}).Info("Berhasil menemukan transaksi berdasarkan Rekening ID")
	return FindNasabahID, nil
}

func (u *allUsecase) FindByNoREK(noREK string) (model.Rekening, error) {
	utils.Log.WithFields(logrus.Fields{
		"no_rekening": noREK,
		"action":      "FindByNoREK",
		"layer":       "allUsecase",
	}).Info("Mencari rekening berdasarkan Nomor Rekening")
	FindNoREK, err := u.RekeningRepository.FindByNoREK(noREK)
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"no_rekening": noREK,
			"error":       err,
			"action":      "FindByNoREK",
			"layer":       "allUsecase",
		}).Error("Gagal mencari rekening berdasarkan Nomor Rekening")
		return model.Rekening{}, err
	}
	utils.Log.WithFields(logrus.Fields{
		"no_rekening": noREK,
		"action":      "FindByNoREK",
		"layer":       "allUsecase",
	}).Info("Berhasil menemukan rekening berdasarkan Nomor Rekening")
	return FindNoREK, nil
}

func (u *allUsecase) Tarik(newTarik model.Transaksi) (model.Transaksi, error) {
	utils.Log.WithFields(logrus.Fields{
		"no_rekening": newTarik.Rekening.NoRekening,
		"nominal":     newTarik.Nominal,
		"action":      "create transaksi tarik",
		"layer":       "allUsecase",
	}).Info("menerima permintaan pembuatan transaksi tarik")

	rekening, err := u.RekeningRepository.FindByNoREK(newTarik.Rekening.NoRekening)
	if err != nil || rekening.ID == 0 {
		utils.Log.WithFields(logrus.Fields{
			"no_rekening": newTarik.Rekening.NoRekening,
			"error":       err,
			"action":      "FindByNoREK",
			"layer":       "allUsecase",
		}).Warn("Rekening tidak ditemukan")
		return model.Transaksi{}, errors.New("rekening tidak ditemukan")
	}

	if newTarik.Nominal != math.Floor(newTarik.Nominal) {
		utils.Log.WithFields(logrus.Fields{
			"nominal": newTarik.Nominal,
			"action":  "validasi nominal bulat",
			"layer":   "allUsecase",
		}).Warn("Nominal tidak boleh desimal")
		return model.Transaksi{}, errors.New("nominal harus bilangan bulat")
	}

	if newTarik.Nominal <= 0 {
		utils.Log.WithFields(logrus.Fields{
			"nominal": newTarik.Nominal,
			"action":  "validasi nominal tarik",
			"layer":   "allUsecase",
		}).Warn("Nominal harus lebih dari 0")
		return model.Transaksi{}, errors.New("nominal harus lebih dari 0")
	}

	if rekening.Saldo < newTarik.Nominal {
		utils.Log.WithFields(logrus.Fields{
			"saldo":   rekening.Saldo,
			"nominal": newTarik.Nominal,
			"action":  "saldo kurang dari nominal",
			"layer":   "allUsecase",
		}).Warn("Saldo tidak mencukupi untuk melakukan transaksi tarik")
		return model.Transaksi{}, errors.New("saldo tidak mencukupi")
	}

	rekening.Saldo -= newTarik.Nominal

	_, err = u.RekeningRepository.UpdateSaldo(rekening)
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"rekening_id": rekening.ID,
			"error":       err,
			"action":      "update saldo",
			"layer":       "allUsecase",
		}).Error("Gagal update saldo setelah pengurangan")
		return model.Transaksi{}, err
	}
	newTransaksi := model.Transaksi{
		RekeningID:     rekening.ID,
		JenisTransaksi: "tarik",
		Nominal:        newTarik.Nominal,
	}

	transaksiTarik, err := u.TransaksiRepository.Tabung(newTransaksi)
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"no_rekening": newTarik.Rekening.NoRekening,
			"nominal":     newTarik.Nominal,
			"error":       err,
			"action":      "Tarik",
			"layer":       "allUsecase",
		}).Error("Gagal mencatat transaksi tarik, rollback saldo")
		rekening.Saldo += newTransaksi.Nominal
		u.RekeningRepository.UpdateSaldo(rekening)
		return model.Transaksi{}, err
	}

	utils.Log.WithFields(logrus.Fields{
		"transaksi_id": transaksiTarik.ID,
		"rekening_id":  rekening.ID,
		"nominal":      transaksiTarik.Nominal,
		"jenis":        transaksiTarik.JenisTransaksi,
		"action":       "Tarik",
		"layer":        "allUsecase",
	}).Info("Transaksi tarik berhasil dibuat")
	return transaksiTarik, nil
}

func (u *allUsecase) Tabung(newTabung model.Transaksi) (model.Transaksi, error) {
	utils.Log.WithFields(logrus.Fields{
		"no_rekening": newTabung.Rekening.NoRekening,
		"nominal":     newTabung.Nominal,
		"action":      "create transaksi tabung",
		"layer":       "allUsecase",
	}).Info("menerima permintaan pembuatan transaksi tabung")
	rekening, err := u.RekeningRepository.FindByNoREK(newTabung.Rekening.NoRekening)
	if err != nil || rekening.ID == 0 {
		utils.Log.WithFields(logrus.Fields{
			"no_rekening": newTabung.Rekening.NoRekening,
			"error":       err,
			"action":      "FindByNoREK",
			"layer":       "allUsecase",
		}).Warn("Rekening tidak ditemukan")
		return model.Transaksi{}, errors.New("rekening tidak ditemukan")
	}

	if newTabung.Nominal <= 0 {
		utils.Log.WithFields(logrus.Fields{
			"nominal": newTabung.Nominal,
			"action":  "validasi nominal tabung",
			"layer":   "allUsecase",
		}).Warn("Nominal harus lebih dari 0")
		return model.Transaksi{}, errors.New("nominal harus lebih dari 0")
	}

	rekening.Saldo += newTabung.Nominal

	_, err = u.RekeningRepository.UpdateSaldo(rekening)
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"rekening_id": rekening.ID,
			"error":       err,
			"action":      "update saldo",
			"layer":       "allUsecase",
		}).Error("Gagal update saldo setelah pengurangan")
		return model.Transaksi{}, err
	}
	newTransaksi := model.Transaksi{
		RekeningID:     rekening.ID,
		JenisTransaksi: "tabung",
		Nominal:        newTabung.Nominal,
	}

	transaksiTabung, err := u.TransaksiRepository.Tabung(newTransaksi)
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"no_rekening": newTabung.Rekening.NoRekening,
			"nominal":     newTabung.Nominal,
			"error":       err,
			"action":      "create tabung",
			"layer":       "allUsecase",
		}).Error("Gagal mencatat transaksi tabung, rollback saldo")
		rekening.Saldo -= newTransaksi.Nominal
		u.RekeningRepository.UpdateSaldo(rekening)
		return model.Transaksi{}, err
	}

	utils.Log.WithFields(logrus.Fields{
		"transaksi_id": transaksiTabung.ID,
		"rekening_id":  rekening.ID,
		"nominal":      transaksiTabung.Nominal,
		"jenis":        transaksiTabung.JenisTransaksi,
		"action":       "create tabung",
		"layer":        "allUsecase",
	}).Info("Transaksi tabung berhasil dibuat")
	return transaksiTabung, nil
}

func NewUsecase(nasabahRepository repository.NasabahRepository, rekeningRepository repository.RekeningRepository, transaksiRepository repository.TransaksiRepository) AllUsecase {
	return &allUsecase{
		NasabahRepository:   nasabahRepository,
		RekeningRepository:  rekeningRepository,
		TransaksiRepository: transaksiRepository,
	}
}
