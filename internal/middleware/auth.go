package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "healthcare-portal/internal/utils"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        claims, err := utils.ValidateJWT(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Set("userID", claims.UserID)
        c.Set("email", claims.Email)
        c.Set("role", claims.Role)
        c.Next()
    }
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists {
                        c.JSON(http.StatusForbidden, gin.H{"error": "Role not found"})
            c.Abort()
            return
        }

        roleStr := userRole.(string)
        authorized := false
        for _, role := range roles {
            if roleStr == role {
                authorized = true
                break
            }
        }

        if !authorized {
            c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }

        c.Next()
    }
}