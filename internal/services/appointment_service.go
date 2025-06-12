package services

import (
    "errors"
    "time"
    
    "healthcare-portal/internal/models"
    "healthcare-portal/internal/repository"
)

type AppointmentService interface {
    CreateAppointment(appointment *models.Appointment) error
    GetAppointmentByID(id uint) (*models.Appointment, error)
    GetAllAppointments(limit, offset int) ([]models.Appointment, int64, error)
    GetAppointmentsByDate(date time.Time) ([]models.Appointment, error)
    GetPatientAppointments(patientID uint) ([]models.Appointment, error)
    GetDoctorAppointments(doctorID uint) ([]models.Appointment, error)
    UpdateAppointmentStatus(id uint, status models.AppointmentStatus) error
    DeleteAppointment(id uint) error
    CheckDoctorAvailability(doctorID uint, date time.Time, timeSlot string) (bool, error)
}

type appointmentService struct {
    appointmentRepo repository.AppointmentRepository
    patientRepo     repository.PatientRepository
    userRepo        repository.UserRepository
}

func NewAppointmentService(appointmentRepo repository.AppointmentRepository, patientRepo repository.PatientRepository, userRepo repository.UserRepository) AppointmentService {
    return &appointmentService{
        appointmentRepo: appointmentRepo,
        patientRepo:     patientRepo,
        userRepo:        userRepo,
    }
}

func (s *appointmentService) CreateAppointment(appointment *models.Appointment) error {
    // Check if doctor is available
    available, err := s.CheckDoctorAvailability(appointment.DoctorID, appointment.Date, appointment.Time)
    if err != nil {
        return err
    }
    if !available {
        return errors.New("doctor is not available at this time")
    }
    
    return s.appointmentRepo.Create(appointment)
}

func (s *appointmentService) GetAppointmentByID(id uint) (*models.Appointment, error) {
    return s.appointmentRepo.FindByID(id)
}

func (s *appointmentService) GetAllAppointments(limit, offset int) ([]models.Appointment, int64, error) {
    return s.appointmentRepo.FindAll(limit, offset)
}

func (s *appointmentService) GetAppointmentsByDate(date time.Time) ([]models.Appointment, error) {
    return s.appointmentRepo.FindByDate(date)
}

func (s *appointmentService) GetPatientAppointments(patientID uint) ([]models.Appointment, error) {
    return s.appointmentRepo.FindByPatientID(patientID)
}

func (s *appointmentService) GetDoctorAppointments(doctorID uint) ([]models.Appointment, error) {
    return s.appointmentRepo.FindByDoctorID(doctorID)
}

func (s *appointmentService) UpdateAppointmentStatus(id uint, status models.AppointmentStatus) error {
    return s.appointmentRepo.UpdateStatus(id, status)
}

func (s *appointmentService) DeleteAppointment(id uint) error {
    return s.appointmentRepo.Delete(id)
}

func (s *appointmentService) CheckDoctorAvailability(doctorID uint, date time.Time, timeSlot string) (bool, error) {
    appointments, err := s.appointmentRepo.FindByDate(date)
    if err != nil {
        return false, err
    }
    
    for _, appointment := range appointments {
        if appointment.DoctorID == doctorID && appointment.Time == timeSlot && appointment.Status != models.StatusCancelled {
            return false, nil
        }
    }
    
    return true, nil
}