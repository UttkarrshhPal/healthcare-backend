package utils

import (
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

type PasswordResetClaims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func GenerateJWT(userID uint, email string, role string) (string, error) {
    claims := &Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateJWT(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

func GeneratePasswordResetToken(userID uint, email string) (string, error) {
    claims := &PasswordResetClaims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // 1 hour expiry
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    // Use a different secret for password reset tokens for added security
    resetSecret := os.Getenv("JWT_SECRET") + "-reset"
    return token.SignedString([]byte(resetSecret))
}

func ValidatePasswordResetToken(tokenString string) (*PasswordResetClaims, error) {
    resetSecret := os.Getenv("JWT_SECRET") + "-reset"
    token, err := jwt.ParseWithClaims(tokenString, &PasswordResetClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(resetSecret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*PasswordResetClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid reset token")
}