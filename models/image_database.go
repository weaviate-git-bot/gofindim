package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type NullImageModel struct {
	ImageModel *ImageModel
}

func (ni *NullImageModel) Scan(value interface{}) error {
	if value == nil {
		ni.ImageModel = nil
		return nil
	}
	var hash, filename string
	err := json.Unmarshal(value.([]byte), &struct {
		Hash     string
		Filename string
	}{hash, filename})
	if err != nil {
		return err
	}
	ni.ImageModel = &ImageModel{Hash: hash, Filename: filename}
	return nil
}

type GroupedImageWithDistance struct {
	Image1Id        uint `gorm:"column:image1_id"`
	Image2Id        uint `gorm:"column:image2_id"`
	Image1Filename  string
	Image2Filename  string
	HammingDistance int `gorm:"-"`
}

func FindDuplicateGroups(db *gorm.DB, threshold int) ([]GroupedImageWithDistance, error) {
	var groups []GroupedImageWithDistance

	err := db.Raw(`SELECT im1.id as image1_id, im2.id as image2_id, im1.filename as image1_filename, im2.filename as image2_filename, hamming_distance(im1.Hash, im2.Hash, im1.Filename, im2.Filename) as hamming_distance
		FROM image_models im1, image_models im2
		WHERE im1.id != im2.id AND hamming_distance(im1.Hash, im2.Hash, im1.Filename, im2.Filename) <= ?`,
		threshold).Scan(&groups).Error

	return groups, err
}

// type DuplicateGroup struct {
// 	Image1          ImageModel
// 	Image2          ImageModel
// 	HammingDistance int
// }

// func FindDuplicateGroups(db *gorm.DB, threshold int) ([]DuplicateGroup, error) {
// 	var duplicates []DuplicateGroup
// 	err := db.Raw(`
// 		SELECT a.*, b.*, hamming_distance(a.hash, b.hash, a.filename, b.filename) as hamming_distance
// 		FROM image_models AS a
// 		JOIN image_models AS b ON a.id != b.id
// 		WHERE hamming_distance(a.hash, b.hash, a.filename, b.filename) <= ?`,
// 		threshold).Scan(&duplicates).Error
// 	return duplicates, err
// }

type ImageWithDistance struct {
	ImageModel
	HammingDistance int `gorm:"-"`
}

func FindSimilarImages(db *gorm.DB, image *ImageModel, threshold int) ([]ImageWithDistance, error) {
	var images []ImageWithDistance
	err := db.Raw(`SELECT *, hamming_distance(Hash, ?, Filename, ?) as hamming_distance FROM image_models
		WHERE hamming_distance(Hash, ?, Filename, ?) <= ?`,
		image.Hash, image.Filename, image.Hash, image.Filename, threshold).Scan(&images).Error
	return images, err
}
