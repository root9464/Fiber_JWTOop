package auth

import (
	"root/services/consts"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func (method *DataAuthService) AccessTokenUpdate(c *fiber.Ctx) error {
    type RefreshRequest struct {
        RefreshToken string `json:"jwt_refresh_token"`
    }

    var request RefreshRequest
    if err := c.BodyParser(&request); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(consts.BadReuest + err.Error())
    }

    token, err := jwt.Parse(request.RefreshToken, func(token *jwt.Token) (interface{}, error) {
        return []byte("your_refresh_secret_key"), nil
    })
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorGenHash + err.Error())
    }

    if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).SendString("Invalid refresh token")
    }

    claims, _ := token.Claims.(jwt.MapClaims)
    userID := int(claims["user_id"].(float64))

    accessClaims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    }
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    access, err := accessToken.SignedString([]byte("your_access_secret_key"))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Ошибка обновления access token: " + err.Error())
    }

    if err := method.DB.Model(&Token{}).Where("user_id = ?", userID).Update("jwt_access_token", access).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Ошибка обновления access token: " + err.Error())
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "jwt_access_token": "Bearer " + access,
    })
}


func (method *DataAuthService) CreateTokens(userID int) (string, string, error) {
    refreshClaims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24 * 4).Unix(),
    }
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refresh, err := refreshToken.SignedString([]byte("your_refresh_secret_key"))
    if err != nil {
        return "", "", err
    }

    if err := method.DB.Model(&Token{}).Where("user_id = ?", userID).Update("jwt_refresh_token", refresh).Error; err != nil {
        return "", "", err
    }

    accessClaims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    }
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    access, err := accessToken.SignedString([]byte("your_access_secret_key"))
    if err != nil {
        return "", "", err
    }

    if err := method.DB.Model(&Token{}).Where("user_id = ?", userID).Update("jwt_access_token", access).Error; err != nil {
        return "", "", err
    }

    return refresh, access, nil
}

// зомби код
func (method *DataAuthService) CreateRefreshToken(userID int) (string, error) {
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 10 * 7).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refresh, err := refreshToken.SignedString([]byte("your_refresh_secret_key"))
	if err != nil {
		return "", err
	}

	if err := method.DB.Model(&Token{}).Where("user_id = ?", userID).Update("jwt_refresh_token", refresh).Error; err != nil {
		return "", err
	}

	return refresh, nil
}

func (method *DataAuthService) CreateAccessToken(userID int) (string, error) {
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 10).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	access, err := accessToken.SignedString([]byte("your_access_secret_key"))
	if err != nil {
		return "", err
	}

	if err := method.DB.Model(&Token{}).Where("user_id = ?", userID).Update("jwt_access_token", access).Error; err != nil {
		return "", err
	}

	return access, nil
}