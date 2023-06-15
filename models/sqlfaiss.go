package models

import (
	_ "github.com/mattn/go-sqlite3" // import sqlite3 driver
	"gocv.io/x/gocv"
	"gorm.io/gorm"
)

func insertImage(db *gorm.DB, imgF *ImageFile, feature []float32) {
	img := imgF.ToImageModel()
	orb := gocv.NewORBWithParams(1000, 1.2, 8, 31, 0, 2, gocv.ORBScoreTypeHarris, 31, 20)
	defer orb.Close()
	_, desc := imgF.toKeypointsDescriptors(&orb)

	defer desc.Close()
	// featureBytes := desc.ToBytes()
	img.Insert(db)

	//		panic(err)
	//		if err != nil {
	//	}
}
