package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/google/uuid"
)

type UploadedFile struct {
	Id          string
	VersionName string
	UUID        string
	Hash        string
	Ext         string
	Mime        string
	Size        int
	Url         string
	CreatedAt   time.Time
	io.Reader
	UploadedFiles []UploadedFile
}

func GetUploadedFilesByVersionName(limit int, versionName string, f *UploadedFile, d *gorm.DB) *gorm.DB {
	return d.Limit(limit).Order("id desc").Where("version_name = ?", versionName).Find(&f.UploadedFiles)
}

func GetUploadedFiles(limit int, f *UploadedFile, d *gorm.DB) *gorm.DB {
	return d.Limit(limit).Order("id asc").Find(&f.UploadedFiles)
}

func StoreUploadedFile(f *UploadedFile, d *gorm.DB) *gorm.DB {
	return d.Create(f)
}

func NewFile(f multipart.File, h *multipart.FileHeader) (*UploadedFile, error) {
	hash, err := generateUUID()
	if err != nil {
		return nil, err
	}
	s := sha256.Sum256([]byte(h.Filename))
	buff := make([]byte, h.Size)
	if _, err := f.Read(buff); err != nil {
		return nil, err
	} else {
		if _, err = f.Seek(0, 0); err != nil {
			return nil, err
		}
	}
	fmt.Println(h.Filename)
	ext := filepath.Ext(h.Filename)
	return &UploadedFile{
		VersionName: strings.Replace(h.Filename, " ", "-", -1),
		Hash:        hex.EncodeToString(s[:]) + ext,
		UUID:        hash + ext,
		Ext:         ext,
		Mime:        http.DetectContentType(buff),
		Size:        int(h.Size),
		Reader:      f,
	}, nil
}

func generateUUID() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	str := u.String()
	return str, nil
}
