package config

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Database struct {
	Conn *gorm.DB
}

var (
	instance *Database
	once     sync.Once
)

// GetDBInstance returns a singleton instance of Database
func GetDBInstance(config DBConfig) *Database {
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.User, config.Password, config.Host, config.Port, config.DBName)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		instance = &Database{Conn: db}
		log.Println("Database connection established")
	})

	return instance
}

// Close closes the database connection
func (db *Database) Close() {
	if db.Conn != nil {
		sqlDB, err := db.Conn.DB()
		if err != nil {
			log.Printf("Error getting SQL DB: %v", err)
			return
		}
		sqlDB.Close()
		log.Println("Database connection closed")
	}
}
