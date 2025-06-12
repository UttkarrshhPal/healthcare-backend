package handlers

import (
	"net/http"
	"strconv"
	"time"

	"healthcare-portal/internal/models"
	"healthcare-portal/internal/services"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	patientService services.PatientService
}

func NewPatientHandler(patientService services.PatientService) *PatientHandler {
	return &PatientHandler{patientService: patientService}
}

type CreatePatientRequest struct {
	FirstName         string `json:"first_name" binding:"required"`
	LastName          string `json:"last_name" binding:"required"`
	Email             string `json:"email" binding:"email"`
	Phone             string `json:"phone" binding:"required"`
	DateOfBirth       string `json:"date_of_birth" binding:"required"`
	Gender            string `json:"gender" binding:"required"`
	Address           string `json:"address"`
	MedicalHistory    string `json:"medical_history"`
	CurrentMedication string `json:"current_medication"`
	Allergies         string `json:"allergies"`
	EmergencyContact  string `json:"emergency_contact"`
	BloodGroup        string `json:"blood_group"`
	InsuranceNumber   string `json:"insurance_number"`
}

// @Summary Create Patient
// @Description Create a new patient (Receptionist only)
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePatientRequest true "Patient details"
// @Success 201 {object} models.Patient
// @Router /api/patients [post]
func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var req CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	patient := &models.Patient{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Email:             req.Email,
		Phone:             req.Phone,
		Gender:            req.Gender,
		Address:           req.Address,
		MedicalHistory:    req.MedicalHistory,
		CurrentMedication: req.CurrentMedication,
		Allergies:         req.Allergies,
		EmergencyContact:  req.EmergencyContact,
		BloodGroup:        req.BloodGroup,
		InsuranceNumber:   req.InsuranceNumber,
		RegisteredBy:      userID.(uint),
		LastUpdatedBy:     userID.(uint),
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
	patient.DateOfBirth = dob

	if err := h.patientService.CreatePatient(patient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, patient)
}

// @Summary Get All Patients
// @Description Get all patients with pagination
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /api/patients [get]
func (h *PatientHandler) GetAllPatients(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	patients, total, err := h.patientService.GetAllPatients(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"patients": patients,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// @Summary Get Patient by ID
// @Description Get a patient by ID
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Patient ID"
// @Success 200 {object} models.Patient
// @Router /api/patients/{id} [get]
func (h *PatientHandler) GetPatientByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	patient, err := h.patientService.GetPatientByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

// @Summary Update Patient
// @Description Update a patient's information
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Patient ID"
// @Param request body CreatePatientRequest true "Updated patient details"
// @Success 200 {object} models.Patient
// @Router /api/patients/{id} [put]
func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	var req CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	patient, err := h.patientService.GetPatientByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	userID, _ := c.Get("userID")

	// Update patient fields
	patient.FirstName = req.FirstName
	patient.LastName = req.LastName
	patient.Email = req.Email
	patient.Phone = req.Phone
	patient.Gender = req.Gender
	patient.Address = req.Address
	patient.MedicalHistory = req.MedicalHistory
	patient.CurrentMedication = req.CurrentMedication
	patient.Allergies = req.Allergies
	patient.EmergencyContact = req.EmergencyContact
	patient.BloodGroup = req.BloodGroup
	patient.InsuranceNumber = req.InsuranceNumber
	patient.LastUpdatedBy = userID.(uint)

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
	patient.DateOfBirth = dob

	if err := h.patientService.UpdatePatient(patient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, patient)
}

// @Summary Delete Patient
// @Description Delete a patient (Receptionist only)
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Patient ID"
// @Success 200 {object} map[string]string
// @Router /api/patients/{id} [delete]
func (h *PatientHandler) DeletePatient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	if err := h.patientService.DeletePatient(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Patient deleted successfully"})
}

// @Summary Search Patients
// @Description Search patients by name, email, or phone
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Success 200 {array} models.Patient
// @Router /api/patients/search [get]
func (h *PatientHandler) SearchPatients(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}

	patients, err := h.patientService.SearchPatients(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, patients)
}
