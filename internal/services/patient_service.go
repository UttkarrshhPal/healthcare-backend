package services

import (
    "healthcare-portal/internal/models"
    "healthcare-portal/internal/repository"
)

type PatientService interface {
    CreatePatient(patient *models.Patient) error
    GetPatientByID(id uint) (*models.Patient, error)
    GetAllPatients(limit, offset int) ([]models.Patient, int64, error)
    UpdatePatient(patient *models.Patient) error
    DeletePatient(id uint) error
    SearchPatients(query string) ([]models.Patient, error)
}

type patientService struct {
    patientRepo repository.PatientRepository
}

func NewPatientService(patientRepo repository.PatientRepository) PatientService {
    return &patientService{patientRepo: patientRepo}
}

func (s *patientService) CreatePatient(patient *models.Patient) error {
    return s.patientRepo.Create(patient)
}

func (s *patientService) GetPatientByID(id uint) (*models.Patient, error) {
    return s.patientRepo.FindByID(id)
}

func (s *patientService) GetAllPatients(limit, offset int) ([]models.Patient, int64, error) {
    return s.patientRepo.FindAll(limit, offset)
}

func (s *patientService) UpdatePatient(patient *models.Patient) error {
    return s.patientRepo.Update(patient)
}

func (s *patientService) DeletePatient(id uint) error {
    return s.patientRepo.Delete(id)
}

func (s *patientService) SearchPatients(query string) ([]models.Patient, error) {
    return s.patientRepo.Search(query)
}