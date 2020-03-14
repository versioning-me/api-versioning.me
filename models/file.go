package models

import (
	"github.com/jinzhu/gorm"
	"sort"
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

type ById []File
func (a ById) Len() int { return len(a) }
func (a ById) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ById) Less(i, j int) bool { return a[i].Id < a[j].Id }

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
	r := d.Limit(limit).Order("id desc").Find(&f.Files)
	//slack形式のソート
	sort.Sort(ById(f.Files))
	return r
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