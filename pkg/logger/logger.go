/*
Copyright (c) 2024 Emin Muhammadi

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
package logger

import (
	"log"
	"os"
	"time"

	httpLoggerMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm/logger"
)

var DatabaseLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.Ldate|log.LUTC|log.Lmicroseconds|log.Llongfile),
	logger.Config{
		SlowThreshold: 400 * time.Millisecond,
		LogLevel: func() logger.LogLevel {
			if os.Getenv("SQL_VERBOSE") == "true" {
				return logger.Info
			}

			return logger.Silent
		}(),
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      false,
		Colorful:                  true,
	},
)

var HttpLogger = httpLoggerMiddleware.New(httpLoggerMiddleware.Config{
	Format:     "\n${time} ${method} ${path} ${status} ${ip}\r\n",
	TimeFormat: "2006/01/02 15:04:05.437903",
	TimeZone:   "UTC",
	Output:     os.Stdout,
})

var Log = log.New(os.Stdout, "\r\n", log.Ldate|log.LUTC|log.Lmicroseconds|log.Llongfile)
