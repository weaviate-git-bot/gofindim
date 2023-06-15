package data

import (
	"gorm.io/gorm"
)

type ImageModel struct {
	ID       uint   `gorm:"primaryKey"`
	Filename string `gorm:"unique"`
	Feature  []byte
	Hash     string
	Rating   int8
	Tags     []Tag `gorm:"many2many:image_tags;"`
}
type Tag struct {
	ID     uint `gorm:"primaryKey"`
	Name   string
	Images []ImageModel `gorm:"many2many:image_tags"`
}
type ImageTag struct {
	ImageID uint `gorm:"primaryKey"`
	TagID   uint `gorm:"primaryKey"`
}

func (m *ImageModel) Insert(db *gorm.DB) error {
	return db.Create(m).Error

}
