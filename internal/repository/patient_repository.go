package repository

import (
    "healthcare-portal/internal/models"
    "gorm.io/gorm"
)

type PatientRepository interface {
    Create(patient *models.Patient) error
    FindAll(limit, offset int) ([]models.Patient, int64, error)
    FindByID(id uint) (*models.Patient, error)
    Update(patient *models.Patient) error
    Delete(id uint) error
    Search(query string) ([]models.Patient, error)
}

type patientRepository struct {
    db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
    return &patientRepository{db: db}
}

func (r *patientRepository) Create(patient *models.Patient) error {
    return r.db.Create(patient).Error
}

func (r *patientRepository) FindAll(limit, offset int) ([]models.Patient, int64, error) {
    var patients []models.Patient
    var total int64

    err := r.db.Model(&models.Patient{}).Count(&total).Error
    if err != nil {
        return nil, 0, err
    }

    // Remove Preload for now
    err = r.db.Limit(limit).Offset(offset).
        Order("created_at DESC").
        Find(&patients).Error
    
    return patients, total, err
}

func (r *patientRepository) FindByID(id uint) (*models.Patient, error) {
    var patient models.Patient
    err := r.db.First(&patient, id).Error
    if err != nil {
        return nil, err
    }
    return &patient, nil
}

func (r *patientRepository) Update(patient *models.Patient) error {
    return r.db.Save(patient).Error
}

func (r *patientRepository) Delete(id uint) error {
    return r.db.Delete(&models.Patient{}, id).Error
}

func (r *patientRepository) Search(query string) ([]models.Patient, error) {
    var patients []models.Patient
    searchQuery := "%" + query + "%"
    err := r.db.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR phone ILIKE ?",
        searchQuery, searchQuery, searchQuery, searchQuery).
        Find(&patients).Error
    return patients, err
}