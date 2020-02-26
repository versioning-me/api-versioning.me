package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Version struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Versions  []Version
}

func NewVersion(f *UploadedFile) *Version {
	return &Version{
		Name: f.VersionName,
	}
}

func StoreVersion(v *Version, d *gorm.DB) *gorm.DB {
	return d.Create(v)
}

func GetVersions(limit int, v *Version, d *gorm.DB) *gorm.DB {
	return d.Limit(limit).Order("id desc").Find(&v.Versions)
}
