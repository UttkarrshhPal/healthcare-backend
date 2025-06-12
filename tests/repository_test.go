package tests

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "healthcare-portal/internal/models"
    "healthcare-portal/internal/repository"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    assert.NoError(t, err)

    err = db.AutoMigrate(&models.User{}, &models.Patient{})
    assert.NoError(t, err)

    return db
}

func TestUserRepository(t *testing.T) {
    db := setupTestDB(t)
    userRepo := repository.NewUserRepository(db)

    t.Run("Create User", func(t *testing.T) {
        user := &models.User{
            Email:    "test@example.com",
            Password: "hashedpassword",
            Name:     "Test User",
            Role:     models.RoleReceptionist,
        }

        err := userRepo.Create(user)
        assert.NoError(t, err)
        assert.NotZero(t, user.ID)
    })

    t.Run("Find User by Email", func(t *testing.T) {
        user, err := userRepo.FindByEmail("test@example.com")
        assert.NoError(t, err)
        assert.NotNil(t, user)
        assert.Equal(t, "test@example.com", user.Email)
    })

    t.Run("Find User by ID", func(t *testing.T) {
        user, err := userRepo.FindByID(1)
        assert.NoError(t, err)
        assert.NotNil(t, user)
        assert.Equal(t, uint(1), user.ID)
    })
}

func TestPatientRepository(t *testing.T) {
    db := setupTestDB(t)
    patientRepo := repository.NewPatientRepository(db)

    t.Run("Create Patient", func(t *testing.T) {
        patient := &models.Patient{
            FirstName:    "John",
            LastName:     "Doe",
            Email:        "john.doe@example.com",
            Phone:        "1234567890",
            DateOfBirth:  time.Now().AddDate(-30, 0, 0),
            Gender:       "Male",
            RegisteredBy: 1,
        }

        err := patientRepo.Create(patient)
        assert.NoError(t, err)
        assert.NotZero(t, patient.ID)
    })

    t.Run("Find All Patients", func(t *testing.T) {
        patients, total, err := patientRepo.FindAll(10, 0)
        assert.NoError(t, err)
        assert.Equal(t, int64(1), total)
        assert.Len(t, patients, 1)
    })

    t.Run("Search Patients", func(t *testing.T) {
        patients, err := patientRepo.Search("John")
        assert.NoError(t, err)
        assert.Len(t, patients, 1)
        assert.Equal(t, "John", patients[0].FirstName)
    })
}