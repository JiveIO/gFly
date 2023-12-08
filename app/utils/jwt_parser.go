package utils

import (
	"app/app/domain/permissions"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
)

// TokenMetadata struct to describe metadata in JWT.
type TokenMetadata struct {
	UserID      uuid.UUID
	Credentials map[string]bool
	Expires     int64
}

// ExtractTokenMetadata func to extract metadata from JWT.
func ExtractTokenMetadata(tokenString string) (*TokenMetadata, error) {
	token, err := verifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Setting and checking token and credentials.
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// User ID.
		userID, err := uuid.Parse(claims["id"].(string))
		if err != nil {
			return nil, err
		}

		// Expires time.
		expires := int64(claims["expires"].(float64))

		// User credentials.
		credentials := map[string]bool{
			permissions.ProfileViewPermission:    claims[permissions.ProfileViewPermission].(bool),
			permissions.UpdatePasswordPermission: claims[permissions.UpdatePasswordPermission].(bool),
		}

		return &TokenMetadata{
			UserID:      userID,
			Credentials: credentials,
			Expires:     expires,
		}, nil
	}

	return nil, err
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
