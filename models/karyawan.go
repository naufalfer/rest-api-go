package models

import (
	"time"

	"gorm.io/gorm"
)

type Karyawan struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Nama         string         `json:"nama"`
	NIP          string         `json:"nip"`
	TempatLahir  string         `json:"tempat_lahir"`
	TanggalLahir time.Time      `json:"tanggal_lahir"`
	Umur         int            `json:"umur"`
	Alamat       string         `json:"alamat"`
	Agama        string         `json:"agama"`
	JenisKelamin string         `json:"jenis_kelamin"`
	NoHandphone  string         `json:"no_handphone"`
	Email        string         `json:"email"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
