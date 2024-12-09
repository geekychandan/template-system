package services

import (
	"errors"
	"fmt"
	"template-system/config"
	"template-system/models"
	"template-system/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// var jwtKey = []byte("pratik123")

func RegisterUser(input RegisterInput) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := utils.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func LoginUser(input LoginInput) (string, error) {
	var user models.User
	if err := utils.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("email or password is incorrect")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", errors.New("email or password is incorrect")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SendPasswordResetEmail(email string) error {
	var user models.User
	if err := utils.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("email not found")
		}
		return err
	}

	// Generate reset token
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWT_SECRET))
	if err != nil {
		return err
	}

	body := fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <head>
            <title>Password Reset</title>
        </head>
        <body>
            <p>Hello,</p>
            <p>We received a request to reset your password. Please use the following token to reset your password:</p>
            <p><strong>%s</strong></p>
            <p>Copy and paste this token in the password reset form.</p>
            <p>If you did not request a password reset, please ignore this email.</p>
            <p>Thanks,</p>
            <p>Your Company Team</p>
        </body>
        </html>
    `, tokenString)

	if err := utils.SendHTMLEmail(user.Email, "Password Reset", body); err != nil {
		return err
	}

	return nil
}

func ResetPassword(tokenString, newPassword string) error {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT_SECRET), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return fmt.Errorf("invalid token")
		}
		return fmt.Errorf("error parsing token")
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	// Find user
	var user models.User
	if err := utils.DB.Where("email = ?", claims.Email).First(&user).Error; err != nil {
		return err
	}

	// Update password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	if err := utils.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
