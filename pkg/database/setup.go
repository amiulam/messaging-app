package database

import (
	"fmt"
	"log"
	"os"

	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/pkg/env"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", env.GetEnv("DB_USER", ""), env.GetEnv("DB_PASSWORD", ""), env.GetEnv("DB_HOST", "127.0.0.1"), env.GetEnv("DB_PORT", ""), env.GetEnv("DB_NAME", ""))
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database \n", err.Error())
		os.Exit(1)
	}

	DB.Logger = logger.Default.LogMode(logger.Info)
	err = DB.AutoMigrate(&models.User{}, &models.UserSession{})

	if err != nil {
		log.Fatal("Failed to migrate database ", err)
	}

	fmt.Println("Database migrated successfully")
}

func SetupMongoDB() {
	uri := env.GetEnv("MONGODB_URI", "")
	client, err := mongo.Connect(options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	coll := client.Database("message").Collection("message_history")
	MongoDB = coll

	fmt.Println("successfully connected to mongodb")
}
