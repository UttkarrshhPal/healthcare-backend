package repository

import (
    "time"
    "healthcare-portal/internal/models"
    "gorm.io/gorm"
)

type AppointmentRepository interface {
    Create(appointment *models.Appointment) error
    FindAll(limit, offset int) ([]models.Appointment, int64, error)
    FindByID(id uint) (*models.Appointment, error)
    FindByDate(date time.Time) ([]models.Appointment, error)
    FindByPatientID(patientID uint) ([]models.Appointment, error)
    FindByDoctorID(doctorID uint) ([]models.Appointment, error)
    Update(appointment *models.Appointment) error
    Delete(id uint) error
    UpdateStatus(id uint, status models.AppointmentStatus) error
}

type appointmentRepository struct {
    db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
    return &appointmentRepository{db: db}
}

func (r *appointmentRepository) Create(appointment *models.Appointment) error {
    return r.db.Create(appointment).Error
}

func (r *appointmentRepository) FindAll(limit, offset int) ([]models.Appointment, int64, error) {
    var appointments []models.Appointment
    var total int64

    err := r.db.Model(&models.Appointment{}).Count(&total).Error
    if err != nil {
        return nil, 0, err
    }

    err = r.db.Preload("Patient").Preload("Doctor").Preload("CreatedByUser").
        Limit(limit).Offset(offset).
        Order("date DESC, time DESC").
        Find(&appointments).Error
    
    return appointments, total, err
}

func (r *appointmentRepository) FindByID(id uint) (*models.Appointment, error) {
    var appointment models.Appointment
    err := r.db.Preload("Patient").Preload("Doctor").Preload("CreatedByUser").
        First(&appointment, id).Error
    if err != nil {
        return nil, err
    }
    return &appointment, nil
}

func (r *appointmentRepository) FindByDate(date time.Time) ([]models.Appointment, error) {
    var appointments []models.Appointment
    startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
    endOfDay := startOfDay.Add(24 * time.Hour)
    
    err := r.db.Preload("Patient").Preload("Doctor").Preload("CreatedByUser").
        Where("date >= ? AND date < ?", startOfDay, endOfDay).
        Order("time ASC").
        Find(&appointments).Error
    return appointments, err
}

func (r *appointmentRepository) FindByPatientID(patientID uint) ([]models.Appointment, error) {
    var appointments []models.Appointment
    err := r.db.Preload("Patient").Preload("Doctor").Preload("CreatedByUser").
        Where("patient_id = ?", patientID).
        Order("date DESC, time DESC").
        Find(&appointments).Error
    return appointments, err
}

func (r *appointmentRepository) FindByDoctorID(doctorID uint) ([]models.Appointment, error) {
    var appointments []models.Appointment
    err := r.db.Preload("Patient").Preload("Doctor").Preload("CreatedByUser").
        Where("doctor_id = ?", doctorID).
        Order("date DESC, time DESC").
        Find(&appointments).Error
    return appointments, err
}

func (r *appointmentRepository) Update(appointment *models.Appointment) error {
    return r.db.Save(appointment).Error
}

func (r *appointmentRepository) Delete(id uint) error {
    return r.db.Delete(&models.Appointment{}, id).Error
}

func (r *appointmentRepository) UpdateStatus(id uint, status models.AppointmentStatus) error {
    return r.db.Model(&models.Appointment{}).Where("id = ?", id).Update("status", status).Error
}