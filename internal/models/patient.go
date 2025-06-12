package models

import (
    "time"
    "gorm.io/gorm"
)

type Patient struct {
    ID                uint           `json:"id" gorm:"primaryKey"`
    FirstName         string         `json:"first_name" gorm:"not null"`
    LastName          string         `json:"last_name" gorm:"not null"`
    Email             string         `json:"email" gorm:"uniqueIndex"`
    Phone             string         `json:"phone" gorm:"not null"`
    DateOfBirth       time.Time      `json:"date_of_birth"`
    Gender            string         `json:"gender" gorm:"type:varchar(20)"`
    Address           string         `json:"address"`
    MedicalHistory    string         `json:"medical_history" gorm:"type:text"`
    CurrentMedication string         `json:"current_medication" gorm:"type:text"`
    Allergies         string         `json:"allergies" gorm:"type:text"`
    EmergencyContact  string         `json:"emergency_contact"`
    BloodGroup        string         `json:"blood_group" gorm:"type:varchar(10)"`
    InsuranceNumber   string         `json:"insurance_number" gorm:"type:varchar(50)"`
    RegisteredBy      uint           `json:"registered_by"`
    LastUpdatedBy     uint           `json:"last_updated_by"`
    CreatedAt         time.Time      `json:"created_at"`
    UpdatedAt         time.Time      `json:"updated_at"`
    DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
    
    // Remove the relationship fields for now to avoid issues
    // RegisteredByUser  *User `json:"registered_by_user,omitempty" gorm:"foreignKey:RegisteredBy"`
    // LastUpdatedByUser *User `json:"last_updated_by_user,omitempty" gorm:"foreignKey:LastUpdatedBy"`
}