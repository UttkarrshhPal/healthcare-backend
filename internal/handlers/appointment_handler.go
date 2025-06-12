package handlers

import (
	"net/http"
	"strconv"
	"time"

	"healthcare-portal/internal/models"
	"healthcare-portal/internal/services"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	appointmentService services.AppointmentService
}

func NewAppointmentHandler(appointmentService services.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}

type CreateAppointmentRequest struct {
	PatientID uint   `json:"patient_id" binding:"required"`
	DoctorID  uint   `json:"doctor_id" binding:"required"`
	Date      string `json:"date" binding:"required"`
	Time      string `json:"time" binding:"required"`
	Notes     string `json:"notes"`
}

// @Summary Create Appointment
// @Description Create a new appointment (Receptionist only)
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAppointmentRequest true "Appointment details"
// @Success 201 {object} models.Appointment
// @Router /api/appointments [post]
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// The service layer will handle the validation of patient and doctor
	// We don't need to check here since the service already does these checks

	appointment := &models.Appointment{
		PatientID: req.PatientID,
		DoctorID:  req.DoctorID,
		Date:      date,
		Time:      req.Time,
		Notes:     req.Notes,
		Status:    models.StatusScheduled,
		CreatedBy: userID.(uint),
	}

	if err := h.appointmentService.CreateAppointment(appointment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

// @Summary Get All Appointments
// @Description Get all appointments with pagination
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /api/appointments [get]
func (h *AppointmentHandler) GetAllAppointments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	appointments, total, err := h.appointmentService.GetAllAppointments(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"appointments": appointments,
		"total":        total,
		"page":         page,
		"limit":        limit,
	})
}

// @Summary Get Appointments by Date
// @Description Get appointments for a specific date
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {array} models.Appointment
// @Router /api/appointments/date [get]
func (h *AppointmentHandler) GetAppointmentsByDate(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date parameter required"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	appointments, err := h.appointmentService.GetAppointmentsByDate(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

// @Summary Get Appointment by ID
// @Description Get an appointment by ID
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Appointment ID"
// @Success 200 {object} models.Appointment
// @Router /api/appointments/{id} [get]
func (h *AppointmentHandler) GetAppointmentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	appointment, err := h.appointmentService.GetAppointmentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	c.JSON(http.StatusOK, appointment)
}

// @Summary Update Appointment Status
// @Description Update the status of an appointment
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Appointment ID"
// @Param status body object{status=string} true "New status"
// @Success 200 {object} map[string]string
// @Router /api/appointments/{id}/status [patch]
func (h *AppointmentHandler) UpdateAppointmentStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := models.AppointmentStatus(req.Status)
	if status != models.StatusScheduled && status != models.StatusCompleted && status != models.StatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	if err := h.appointmentService.UpdateAppointmentStatus(uint(id), status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

// @Summary Delete Appointment
// @Description Delete an appointment (Receptionist only)
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Appointment ID"
// @Success 200 {object} map[string]string
// @Router /api/appointments/{id} [delete]
func (h *AppointmentHandler) DeleteAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	if err := h.appointmentService.DeleteAppointment(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment deleted successfully"})
}

// @Summary Get Patient Appointments
// @Description Get all appointments for a specific patient
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param patientId path int true "Patient ID"
// @Success 200 {array} models.Appointment
// @Router /api/appointments/patient/{patientId} [get]
func (h *AppointmentHandler) GetPatientAppointments(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("patientId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	appointments, err := h.appointmentService.GetPatientAppointments(uint(patientID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

// @Summary Get Doctor Appointments
// @Description Get all appointments for a specific doctor
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param doctorId path int true "Doctor ID"
// @Success 200 {array} models.Appointment
// @Router /api/appointments/doctor/{doctorId} [get]
func (h *AppointmentHandler) GetDoctorAppointments(c *gin.Context) {
	doctorID, err := strconv.ParseUint(c.Param("doctorId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	appointments, err := h.appointmentService.GetDoctorAppointments(uint(doctorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}
