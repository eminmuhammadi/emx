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

// Copyright (C) 1992-2020 Free Software Foundation, Inc.
// This file is part of the GNU C Library.
//
// The GNU C Library is free software; you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 2.1 of the License, or (at your option) any later version.
//
// The GNU C Library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with the GNU C Library; if not, see
// <https://www.gnu.org/licenses/>.
//
// SPDX-License-Identifier: GPL-2.0+ WITH Linux-syscall-note
// netinet/in/in_linux_ppc64le.go
// signal/signal_linux_ppc64le.go
// sys/socket/socket_linux_ppc64le.go
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version
// 2 of the License, or (at your option) any later version.
//
// Copyright (C) 1992-2021 Free Software Foundation, Inc.
//
// This file is part of GCC.
//
// GCC is free software; you can redistribute it and/or modify it under
// the terms of the GNU General Public License as published by the Free
// Software Foundation; either version 3, or (at your option) any later
// version.
//
// GCC is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License
// for more details.
//
// Under Section 7 of GPL version 3, you are granted additional
// permissions described in the GCC Runtime Library Exception, version
// 3.1, as published by the Free Software Foundation.
//
// You should have received a copy of the GNU General Public License and
// a copy of the GCC Runtime Library Exception along with this program;
// see the files COPYING3 and COPYING.RUNTIME respectively.  If not, see
// <http://www.gnu.org/licenses/>.
//
// SPDX-License-Identifier: GPL-2.0 WITH Linux-syscall-note
//  arch/arm/include/asm/posix_types.h
//
//  Copyright (C) 1996-1998 Russell King.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 2 as
// published by the Free Software Foundation.
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
