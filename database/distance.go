package database

import (
	"context"

	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

func RegisterHammingDistanceFunc(db *gorm.DB, ctx context.Context) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	conn, err := sqlDB.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.Raw(func(driverConn interface{}) error {
		sqliteConn := driverConn.(*sqlite3.SQLiteConn)
		return sqliteConn.RegisterFunc("hamming_distance", HammingDistance, true)
	})
}
func HammingDistance(hash1, hash2, name1, name2 string) int {
	if name1 == name2 {
		return 100
	}
	distance := 0
	for i := 0; i < len(hash1); i++ {
		if hash1[i] != hash2[i] {
			distance++
		}
	}
	return distance
}
