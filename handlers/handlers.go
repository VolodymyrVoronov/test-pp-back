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
	// For this example, we'll send it back in the response (not secure)
	c.JSON(http.StatusOK, gin.H{"secret": secret.Secret()})
}

func LoginHandler(c *gin.Context) {
	var creds models.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	expectedPassword, ok := models.Users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
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

	c.SetCookie("token", tokenString, 300, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
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
	secret := "" // replace with the actual secret

	if totp.Validate(otp, secret) {
		c.JSON(http.StatusOK, gin.H{"message": "Verification successful"})
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

	c.JSON(http.StatusOK, gin.H{"message": "Authenticated", "username": claims.Username})
}
