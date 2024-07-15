package handlers

import (
	"net/http"
	"test-pp-back/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/pquerna/otp/totp"
)

func RegisterHandler(c *gin.Context) {
	var creds models.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if _, ok := models.Users[creds.Username]; ok {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	models.Users[creds.Username] = creds.Password

	// Generate a TOTP secret for the user
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "YourAppName",
		AccountName: creds.Username,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating TOTP secret"})
		return
	}

	// Store the secret securely
	models.TotpSecrets[creds.Username] = secret.Secret()

	// For this example, we'll send it back in the response (not secure)
	c.JSON(http.StatusOK, gin.H{"secret": secret.Secret()})
}

func LoginHandler(c *gin.Context) {
	var creds models.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if _, ok := models.Users[creds.Username]; !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expectedPassword, ok := models.Users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(models.JWTKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Generate OTP
	otp, err := totp.GenerateCode(models.TotpSecrets[creds.Username], time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating OTP"})
		return
	}

	c.SetCookie("token", tokenString, 300, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "otp": otp})
}

func VerifyHandler(c *gin.Context) {
	tokenStr, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
		return
	}
	otp := c.Query("otp")

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return models.JWTKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Retrieve the user's secret (in a real app, fetch this from your database)
	secret, ok := models.TotpSecrets[claims.Username]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving secret"})
		return
	}

	if totp.Validate(otp, secret) {
		c.JSON(http.StatusOK, gin.H{"message": "Verification successful", "username": claims.Username})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
	}
}

func CheckAuthHandler(c *gin.Context) {
	tokenStr, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
		return
	}

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return models.JWTKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Authenticated"})
}
