package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sferawann/go-bank-api/model"
	"github.com/sferawann/go-bank-api/usecase"
	"github.com/sferawann/go-bank-api/utils"
	"github.com/sirupsen/logrus"
)

type AllController interface {
	Create(ctx echo.Context) error
	Tabung(ctx echo.Context) error
	Tarik(ctx echo.Context) error
	GetSaldo(ctx echo.Context) error
}

type allController struct {
	AllUsecase usecase.AllUsecase
}

func (c *allController) Create(ctx echo.Context) error {
	utils.Log.WithFields(logrus.Fields{
		"action": "bind data create nasabah",
		"layer":  "allController",
	}).Info("Mencoba memproses data req pembuatan nasabah")

	var newNasabah model.Nasabah

	if err := ctx.Bind(&newNasabah); err != nil {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"action": "bind data create nasabah",
			"layer":  "allController",
		}).Error("Format data req tidak valid")
		return ctx.JSON(http.StatusBadRequest, map[string]string{"remark": "Format Data Tidak Valid!"})
	}

	if newNasabah.NIK == "" || newNasabah.Nama == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"remark": "Field nik dan nama wajib diisi",
		})
	}

	utils.Log.WithFields(logrus.Fields{
		"nama":   newNasabah.Nama,
		"nik":    newNasabah.NIK,
		"no_hp":  newNasabah.NoHP,
		"action": "create nasabah",
		"layer":  "allController",
	}).Info("Meneruskan permintaan pembuatan nasabah ke usecase")
	createdNasabah, err := c.AllUsecase.Create(newNasabah)
	if err != nil {
		// fmt.Println("aaaaERROR:", err.Error())
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"nama":   newNasabah.Nama,
			"nik":    newNasabah.NIK,
			"no_hp":  newNasabah.NoHP,
			"action": "create",
			"layer":  "allController",
		}).Error("Gagal membuat nasabah melalui usecase")
		if err.Error() == "nik sudah digunakan" || err.Error() == "no hp sudah digunakan" {
			utils.Log.WithFields(logrus.Fields{
				"error":  err.Error(),
				"action": "validasi",
				"layer":  "allController",
			}).Warn("validasi data gagal")
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"remark": err.Error(),
			})
		}

		utils.Log.WithError(err).WithFields(logrus.Fields{
			"action": "internal_server_error",
			"layer":  "allController",
		}).Error("Terjadi kesalahan pada server saat membuat nasabah")
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"remark": "Terjadi kesalahan pada server",
		})
	}

	utils.Log.WithFields(logrus.Fields{
		"nasabah_id": createdNasabah.ID,
		"action":     "FindByNasabahID",
		"layer":      "allController",
	}).Info("Mencari data rekening berdasarkan ID nasabah")
	rekening, err := c.AllUsecase.FindByNasabahID(createdNasabah.ID)
	if err != nil {
		// fmt.Println("ERROR FIND NASABAH ID:", err.Error())
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"id":     createdNasabah.ID,
			"action": "FindByNasabahID",
			"layer":  "allController",
		}).Error("Gagal mengambil data rekening")
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"remark": "Gagal mengambil data rekening",
		})
	}

	utils.Log.WithFields(logrus.Fields{
		"no_rekening": rekening.NoRekening,
		"action":      "response",
		"layer":       "allController",
	}).Info("Berhasil membuat nasabah dan mengambil nomor rekening")
	return ctx.JSON(http.StatusOK, map[string]string{
		"no_rekening": rekening.NoRekening,
	})
}

func (c *allController) Tabung(ctx echo.Context) error {
	var newTabung model.Transaksi

	utils.Log.WithFields(logrus.Fields{
		"action": "bind data tabung",
		"layer":  "allController",
	}).Info("Mencoba memproses data req pembuatan transaksi tabung")
	if err := ctx.Bind(&newTabung); err != nil {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"action": "bind data tabung",
			"layer":  "allController",
		}).Error("Format data req tidak valid")
		return ctx.JSON(http.StatusBadRequest, map[string]string{"remark": "Format Data Tidak Valid!"})
	}

	utils.Log.WithFields(logrus.Fields{
		"no_rekening": newTabung.Rekening.NoRekening,
		"action":      "find rekening tabung",
		"layer":       "allController",
	}).Info("Mengecek apakah rekening tersedia")
	_, err := c.AllUsecase.FindByNoREK(newTabung.Rekening.NoRekening)
	if err != nil {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"no_rekening": newTabung.Rekening.NoRekening,
			"action":      "find rekening tabung",
			"layer":       "allController",
		}).Warn("Rekening tidak ditemukan")
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"remark": "Rekening tidak ditemukan",
		})
	}

	createdTabung, err := c.AllUsecase.Tabung(newTabung)
	if err != nil {
		// fmt.Println("TABUNG ERORR:", err.Error())
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"no_rekening": newTabung.Rekening.NoRekening,
			"nominal":     newTabung.Nominal,
			"action":      "create transaksi tabung",
			"layer":       "allController",
		}).Error("Gagal melakukan transaksi tabung")
		if err.Error() == "rekening tidak ditemukan" {
			utils.Log.WithError(err).WithFields(logrus.Fields{
				"no_rekening": newTabung.Rekening.NoRekening,
				"action":      "validasi",
				"layer":       "allController",
			}).Error("rekening tidak ditemukan")
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"remark": err.Error(),
			})
		}
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"action": "internal_server_error",
			"layer":  "allController",
		}).Error("Terjadi kesalahan pada server saat membuat transaksi tabung")
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"remark": "Terjadi kesalahan pada server",
		})
	}

	utils.Log.WithFields(logrus.Fields{
		"no_rekening": createdTabung.RekeningID,
		"action":      "FindByRekeningID",
		"layer":       "allController",
	}).Info("Mengecek apakah rekening tersedia berdasarkan rekeningg id")
	rekening, err := c.AllUsecase.FindByRekeningID(createdTabung.RekeningID)
	if err != nil {
		// fmt.Println("ERROR FIND REKENING ID:", err.Error())
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"rekening_id": createdTabung.RekeningID,
			"action":      "FindByRekeningID",
			"layer":       "allController",
		}).Error("Gagal mengambil data rekening setelah tabung")
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"remark": "Gagal mengambil data rekening",
		})
	}

	utils.Log.WithFields(logrus.Fields{
		"no_rekening": rekening.Rekening.NoRekening,
		"saldo":       rekening.Rekening.Saldo,
		"action":      "create transaksi tabung",
		"layer":       "allController",
	}).Info("Berhasil melakukan transaksi tabung")
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"saldo": rekening.Rekening.Saldo,
	})

}

func (c *allController) Tarik(ctx echo.Context) error {
	var newTarik model.Transaksi

	utils.Log.WithFields(logrus.Fields{
		"action": "bind data tarik",
		"layer":  "allController",
	}).Info("Mencoba memproses data req pembuatan transaksi tarik")

	if err := ctx.Bind(&newTarik); err != nil {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"action": "bind data tarik",
			"layer":  "allController",
		}).Error("Format data req tidak valid")
		return ctx.JSON(http.StatusBadRequest, map[string]string{"remark": "Format Data Tidak Valid!"})
	}

	_, err := c.AllUsecase.FindByNoREK(newTarik.Rekening.NoRekening)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"remark": "Rekening tidak ditemukan",
		})
	}

	createdTarik, err := c.AllUsecase.Tarik(newTarik)
	if err != nil {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"no_rekening": newTarik.Rekening.NoRekening,
			"action":      "tarik saldo",
			"layer":       "allController",
		}).Error("Gagal melakukan penarikan saldo")
		// fmt.Println("TARIK ERORR:", err.Error())
		if err.Error() == "rekening tidak ditemukan" || err.Error() == "saldo tidak mencukupi" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"remark": err.Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"remark": "Terjadi kesalahan pada server",
		})
	}
	// return ctx.JSON(http.StatusOK, createdTabung)

	rekening, err := c.AllUsecase.FindByRekeningID(createdTarik.RekeningID)
	if err != nil {
		fmt.Println("ERROR FIND REKENING ID:", err.Error())
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"remark": "Gagal mengambil data rekening",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"saldo": rekening.Rekening.Saldo,
	})

}

func (c *allController) GetSaldo(ctx echo.Context) error {
	noREK := ctx.Param("no_rekening")

	utils.Log.WithFields(logrus.Fields{
		"no_rekening": noREK,
		"action":      "GetSaldo",
		"layer":       "allController",
	}).Info("Menerima permintaan cek saldo")

	rekening, err := c.AllUsecase.FindByNoREK(noREK)
	if err != nil || rekening.ID == 0 {
		utils.Log.WithError(err).WithFields(logrus.Fields{
			"no_rekening": noREK,
			"action":      "GetSaldo",
			"layer":       "allController",
		}).Warn("Rekening tidak ditemukan saat cek saldo")
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"remark": "rekening tidak ditemukan",
		})
	}

	utils.Log.WithFields(logrus.Fields{
		"no_rekening": noREK,
		"saldo":       rekening.Saldo,
		"action":      "GetSaldo",
		"layer":       "allController",
	}).Info("Berhasil mengambil saldo rekening")
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"saldo": rekening.Saldo,
	})
}

func NewController(AllUsecase usecase.AllUsecase) AllController {
	return &allController{AllUsecase}
}
