package main

import (
    "log"
    "time"

    "github.com/joho/godotenv"
    "healthcare-portal/internal/database"
    "healthcare-portal/internal/models"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Initialize database
    database.Initialize()
    db := database.GetDB()

    log.Println("Starting database seeding...")

    // Create default users
    users := []struct {
        Email    string
        Password string
        Name     string
        Role     models.UserRole
    }{
        {
            Email:    "receptionist@healthcare.com",
            Password: "receptionist123",
            Name:     "Sarah Johnson",
            Role:     models.RoleReceptionist,
        },
        {
            Email:    "doctor@healthcare.com",
            Password: "doctor123",
            Name:     "Dr. Michael Smith",
            Role:     models.RoleDoctor,
        },
        {
            Email:    "admin.receptionist@healthcare.com",
            Password: "admin123",
            Name:     "Emily Davis",
            Role:     models.RoleReceptionist,
        },
        {
            Email:    "senior.doctor@healthcare.com",
            Password: "senior123",
            Name:     "Dr. Robert Brown",
            Role:     models.RoleDoctor,
        },
    }

    // Create users
    for _, userData := range users {
        // Check if user already exists
        var existingUser models.User
        result := db.Where("email = ?", userData.Email).First(&existingUser)
        
        if result.Error == nil {
            log.Printf("User %s already exists, skipping...", userData.Email)
            continue
        }

        // Hash password
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
        if err != nil {
            log.Printf("Failed to hash password for %s: %v", userData.Email, err)
            continue
        }

        // Create new user
        user := models.User{
            Email:    userData.Email,
            Password: string(hashedPassword),
            Name:     userData.Name,
            Role:     userData.Role,
            IsActive: true,
        }

        if err := db.Create(&user).Error; err != nil {
            log.Printf("Failed to create user %s: %v", userData.Email, err)
        } else {
            log.Printf("✓ Created user: %s (%s) - Password: %s", userData.Name, userData.Email, userData.Password)
        }
    }

    // Create sample patients (optional)
    log.Println("\nCreating sample patients...")
    patients := []models.Patient{
        {
            FirstName:         "Alice",
            LastName:          "Johnson",
            Email:             "alice.johnson@example.com",
            Phone:             "+1234567890",
            DateOfBirth:       time.Now().AddDate(-30, 0, 0),
            Gender:            "Female",
            Address:           "123 Main St, New York, NY 10001",
            BloodGroup:        "A+",
            InsuranceNumber:   "INS-001234",
            MedicalHistory:    "No significant medical history",
            Allergies:         "None",
            RegisteredBy:      1,
            LastUpdatedBy:     1,
        },
        {
            FirstName:         "Bob",
            LastName:          "Williams",
            Email:             "bob.williams@example.com",
            Phone:             "+0987654321",
            DateOfBirth:       time.Now().AddDate(-45, 0, 0),
            Gender:            "Male",
            Address:           "456 Oak Ave, Los Angeles, CA 90001",
            BloodGroup:        "O+",
            InsuranceNumber:   "INS-005678",
            MedicalHistory:    "Hypertension, Diabetes Type 2",
            CurrentMedication: "Metformin 500mg twice daily, Lisinopril 10mg once daily",
            Allergies:         "Penicillin",
            EmergencyContact:  "+1122334455",
            RegisteredBy:      1,
            LastUpdatedBy:     1,
        },
        {
            FirstName:         "Carol",
            LastName:          "Davis",
            Email:             "carol.davis@example.com",
            Phone:             "+1112223333",
            DateOfBirth:       time.Now().AddDate(-25, 0, 0),
            Gender:            "Female",
            Address:           "789 Pine St, Chicago, IL 60601",
            BloodGroup:        "B+",
            InsuranceNumber:   "INS-009012",
            MedicalHistory:    "Asthma",
            CurrentMedication: "Albuterol inhaler as needed",
            Allergies:         "Dust, Pollen",
            RegisteredBy:      1,
            LastUpdatedBy:     1,
        },
    }

    for _, patient := range patients {
        // Check if patient already exists
        var existingPatient models.Patient
        result := db.Where("email = ?", patient.Email).First(&existingPatient)
        
        if result.Error == nil {
            log.Printf("Patient %s %s already exists, skipping...", patient.FirstName, patient.LastName)
            continue
        }

        if err := db.Create(&patient).Error; err != nil {
            log.Printf("Failed to create patient %s %s: %v", patient.FirstName, patient.LastName, err)
        } else {
            log.Printf("✓ Created patient: %s %s", patient.FirstName, patient.LastName)
        }
    }

    log.Println("\n=== Seeding completed successfully! ===")
    log.Println("\nYou can now login with these credentials:")
    log.Println("\nReceptionist Accounts:")
    log.Println("  Email: receptionist@healthcare.com")
    log.Println("  Password: receptionist123")
    log.Println("")
    log.Println("  Email: admin.receptionist@healthcare.com")
    log.Println("  Password: admin123")
    log.Println("\nDoctor Accounts:")
    log.Println("  Email: doctor@healthcare.com")
    log.Println("  Password: doctor123")
    log.Println("")
    log.Println("  Email: senior.doctor@healthcare.com")
    log.Println("  Password: senior123")
}