/*
Copyright 2024 Emin Muhammadi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package sql

import (
	"os"
	"time"

	"github.com/eminmuhammadi/emx/pkg/logger"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var maxOpenConns = 1
var maxIdleConns = 1
var maxLifetimeInSeconds = 300

// Creates a global in-memory database connection
var Database = func() *gorm.DB {
	db, err := Connect()
	if err != nil {
		logger.Log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}()

func getDSN() string {
	dsn := os.Getenv("SQLITE_DSN")
	if dsn == "" {
		dsn = ":memory:?_pragma=foreign_keys(1)"
	}

	return dsn
}

func TimeNow() time.Time {
	return time.Now().UTC()
}

var gormConf = &gorm.Config{
	PrepareStmt:            true,
	SkipDefaultTransaction: true,
	DryRun:                 false,
	NowFunc: func() time.Time {
		return TimeNow()
	},
	Logger: logger.DatabaseLogger,
}

func ApplyConfig(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	logger.Log.Println("Setting up database connection pool")

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetimeInSeconds) * time.Second)

	logger.Log.Println("Database connection pool set up")

	return nil
}

func Connect() (*gorm.DB, error) {
	logger.Log.Println("Connecting to database")

	db, err := gorm.Open(sqlite.Open(getDSN()), gormConf)
	if err != nil {
		return nil, err
	}

	logger.Log.Println("Connected to database")

	if err := ApplyConfig(db); err != nil {
		return nil, err
	}

	return db, nil
}

func CloseConnection(db *gorm.DB) error {
	logger.Log.Println("Closing database connection")

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	logger.Log.Println("Database connection closed")

	return sqlDB.Close()
}
