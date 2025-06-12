package models

import (
    "time"
    "gorm.io/gorm"
)

type AppointmentStatus string

const (
    StatusScheduled AppointmentStatus = "scheduled"
    StatusCompleted AppointmentStatus = "completed"
    StatusCancelled AppointmentStatus = "cancelled"
)

type Appointment struct {
    ID            uint              `json:"id" gorm:"primaryKey"`
    PatientID     uint              `json:"patient_id" gorm:"not null"`
    DoctorID      uint              `json:"doctor_id" gorm:"not null"`
    Date          time.Time         `json:"date"`
    Time          string            `json:"time" gorm:"type:varchar(10)"`
    Status        AppointmentStatus `json:"status" gorm:"type:varchar(20);default:'scheduled'"`
    Notes         string            `json:"notes" gorm:"type:text"`
    CreatedBy     uint              `json:"created_by"`
    CreatedAt     time.Time         `json:"created_at"`
    UpdatedAt     time.Time         `json:"updated_at"`
    DeletedAt     gorm.DeletedAt    `json:"-" gorm:"index"`
}