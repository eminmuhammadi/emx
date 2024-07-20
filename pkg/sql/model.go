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
package sql

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type BaseModel struct {
	ID        int64                 `gorm:"primaryKey,autoIncrement,unique,not null" json:"id"`
	CreatedAt time.Time             `gorm:"not null,index" json:"created_at"`
	UpdatedAt time.Time             `gorm:"not null,index" json:"updated_at"`
	IsDeleted soft_delete.DeletedAt `gorm:"softDelete:flag" json:"-"`
}
