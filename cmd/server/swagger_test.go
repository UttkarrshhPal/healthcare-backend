package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    
    _ "healthcare-portal/docs"
)

func TestSwaggerRoute(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/swagger/index.html", nil)
    router.ServeHTTP(w, req)
    
    t.Logf("Status: %d", w.Code)
    t.Logf("Body: %s", w.Body.String())
}