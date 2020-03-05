package models

import (
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jinzhu/gorm"
)

type Version struct {
	Id         int
	Name       string
	FileId     int
	VersionNum int
	Ext        string
	Mime       string
	Size       int
	Url        string
	History    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	io.Reader
	Versions []Version
}

func GetVersionsByFileId(limit int, fileId string, f *Version, d *gorm.DB) *gorm.DB {
	return d.Limit(limit).Order("id desc").Where("file_id = ?", fileId).Find(&f.Versions)
}

func GenerateVersionId(v *Version, d *gorm.DB) {
	d.Limit(1).Order("id desc").Where("file_id = ?", v.FileId).Last(&v)
	v.VersionNum = 1
}

func UpdateVersionNum(v *Version, d *gorm.DB) {
	var temp Version
	d.Limit(1).Order("id desc").Where("file_id = ?", v.FileId).Last(&temp)
	v.VersionNum = temp.VersionNum + 1
}

func StoreNewVersion(v *Version, d *gorm.DB) *gorm.DB {
	return d.Create(v)
}

func UpdateVersionHistory(v *Version, historyText string, d *gorm.DB) *gorm.DB {
	v.History = historyText
	return d.Update(v)
}

func (v *Version) ConvertVersionNumToString() string {
	versionStr := strconv.Itoa(v.VersionNum)
	return versionStr
}

func NewVersion(f multipart.File, h *multipart.FileHeader) (*Version, error) {
	buff := make([]byte, h.Size)
	if _, err := f.Read(buff); err != nil {
		return nil, err
	} else {
		if _, err = f.Seek(0, 0); err != nil {
			return nil, err
		}
	}
	ext := filepath.Ext(h.Filename)
	return &Version{
		Name:      strings.Replace(h.Filename, " ", "_", -1),
		Ext:       ext,
		Mime:      mimetype.Detect(buff).String(),
		Size:      int(h.Size),
		Reader:    f,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Versions:  nil,
	}, nil
}

// 不要
// func GetVersions(limit int, f *Version, d *gorm.DB) *gorm.DB {
// 	return d.Limit(limit).Order("id desc").Find(&f.Versions)
// }
