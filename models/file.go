package models

import (
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type File struct {
	Id        int
	Title     string
	Detail    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Files     []File
}

func NewFile() *File {
	return &File{}
}

func NewFileWithTitle(title, detail string) *File {
	return &File{
		Title:  title,
		Detail: detail,
	}
}

func GetFiles(limit int, f *File, d *gorm.DB) *gorm.DB {
	return d.Limit(limit).Order("id desc").Find(&f.Files)
}

func GetFileById(id string, f *File, d *gorm.DB) *gorm.DB {
	return d.Limit(1).Order("id desc").Where("id = ?", id).Find(&f)
}

func (f *File)ConvertFileIdToStoring() string {
	fileIdStr := strconv.Itoa(f.Id)
	return fileIdStr
}

func StoreFile(f *File, d *gorm.DB) *gorm.DB {
	return d.Create(f)
}

func UpdateFile(f *File, d *gorm.DB) *gorm.DB {
	return d.Save(f)
}