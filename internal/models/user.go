package models

import (
    "time"
    "gorm.io/gorm"
)

type UserRole string

const (
    RoleReceptionist UserRole = "receptionist"
    RoleDoctor       UserRole = "doctor"
)

type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Email     string         `json:"email" gorm:"uniqueIndex;not null"`
    Password  string         `json:"-" gorm:"not null"`
    Name      string         `json:"name" gorm:"not null"`
    Role      UserRole       `json:"role" gorm:"type:varchar(50);not null"`
    IsActive  bool           `json:"is_active" gorm:"default:true"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}