package main

import (
	"_x3/sqldb/database"
	"_x3/sqldb/models"
	"_x3/sqldb/utils"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	// galleryURL  string
	initialize      bool
	scanFolder      string
	imageFile       string
	findFolder      string
	similar         int
	imagesToCompare []string
	fpExtractor     = &cobra.Command{
		Use:   "fap [string]",
		Short: "Download a gallery",
		Args:  cobra.ExactArgs(0),
		RunE:  Execute,
	}
)

func Execute(cmd *cobra.Command, args []string) error {
	var db *gorm.DB
	var err error
	if initialize {
		db = initDB("my_images.db")
		fmt.Println("Initializing database")
	} else {
		db, err = gorm.Open(sqlite.Open("file:my_images.db"), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
	}
	_db, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer _db.Close()
	fmt.Println("Connected to database. Executing command")

	if len(imagesToCompare) == 2 {
		img1, err := models.NewImageFileFromPath(imagesToCompare[0])
		if err != nil {
			return err
		}
		img2, err := models.NewImageFileFromPath(imagesToCompare[1])
		if err != nil {
			return err
		}
		distance, err := models.CompareImageOrb(img1, img2)
		if err != nil {
			return err
		}
		fmt.Printf("Distance between %v and %v is %v\n", img1.Name, img2.Name, distance)
		return nil

	}

	if similar > 0 && imageFile != "" {
		findSimilarInDB(db, imageFile)
	} else if similar >= 0 && findFolder == "" {
		findDuplicatesInDB(db)
	} else if scanFolder != "" {
		fmt.Println("Scanning folder")
		err := filepath.Walk(scanFolder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if utils.IsImage(filepath.Ext(path)) {
				fmt.Println("Found image", path)
				image, err := models.NewImageFileFromPath(path)
				if err != nil {
					return err
				}
				imageM := image.ToImageModel()
				err = imageM.Insert(db)
				if err != nil {
					if err.Error() != "UNIQUE constraint failed: image_models.filename" {
						fmt.Println(err)
						return err
					} else {
						return nil
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	} else if imageFile != "" {
		fmt.Println("Hashing image")
		myImage, err := models.NewImageFileFromPath(imageFile)
		if err != nil {
			return err
		}
		fmt.Printf("Hash: %s\n", myImage.ToImageModel().Hash)
		// myImage.ToImageModel().Insert(db)
		// hash, err := models.HashImagePerception(imageFile)
		// fmt.Println(hash)
	}
	return nil
}

func main() {
	// Execute the Cobra command
	fpExtractor.Flags().BoolVarP(&initialize, "initialize", "I", false, "Initialize the database")
	fpExtractor.Flags().StringVarP(&scanFolder, "scan", "S", "", "Scan a directory for images")
	fpExtractor.Flags().StringVarP(&imageFile, "image-file", "i", "", "Image file to hash")
	fpExtractor.Flags().StringVarP(&findFolder, "find-folder", "F", "", "Find similar images in a folder")
	fpExtractor.Flags().StringArrayVarP(&imagesToCompare, "compareImages", "C", []string{}, "Compare two images")
	fpExtractor.Flags().IntVar(&similar, "similar", 0, "Image file to hash")
	if err := fpExtractor.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func initDB(dbFile string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:my_images.db"), &gorm.Config{})
	// db, err := gorm.Open(sql.Open("sqlite3", dbFile), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	tx := db.Exec(`CREATE TABLE IF NOT EXISTS images (
		id INTEGER PRIMARY KEY,
		filename TEXT NOT NULL,
		hash TEXT NOT NULL,
		UNIQUE(filename)
	);`)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
	err = db.AutoMigrate(&models.ImageModel{})
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	db.WithContext(ctx)
	err = database.RegisterHammingDistanceFunc(db, ctx)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func findSimilarInDB(db *gorm.DB, imgPath string) error {
	fmt.Println("Finding similar images")
	myImage, err := models.NewImageFileFromPath(imgPath)
	if err != nil {
		log.Fatal(err)
	}
	imageM := myImage.ToImageModel()
	images, err := models.FindSimilarImages(db, imageM, similar)
	if err != nil {
		log.Fatal(err)
	}
	for _, image := range images {
		fmt.Println(image.Filename, image.HammingDistance)
	}
	return nil
}

func findDuplicatesInDB(db *gorm.DB) error {
	images, err := models.FindDuplicateGroups(db, similar)
	if err != nil {
		return err
	}
	for _, image := range images {
		fmt.Printf("Group %v\n", image)

	}
	return nil
}

// func start() {
// 	// Initialize the database
// 	db := initDB("my_images.db")
// 	defer db.Close()

// 	// Insert data
// 	_, err := db.Exec("INSERT INTO images (filename, hash) VALUES (?, ?)", "example.jpg", "1234556789")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Query data
// 	rows, err := db.Query("SELECT filename, hash FROM images")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	// Iterate over the results
// 	for rows.Next() {
// 		var filename, hash string
// 		err = rows.Scan(&filename, &hash)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("Filename: %s, Hash: %s\n", filename, hash)
// 	}

// 	// Check for errors
// 	if err = rows.Err(); err != nil {
// 		log.Fatal(err)

// 	}
// }
