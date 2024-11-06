package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

// AuthService struct to hold the JWT secret
type AuthService struct {
	JWTSecret []byte
}

// NewAuthService creates a new AuthService
func NewAuthService(secret string) *AuthService {
	return &AuthService{JWTSecret: []byte(secret)}
}

// GenerateToken generates a JWT token for a given user ID
func (auth *AuthService) GenerateToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString(auth.JWTSecret)
}

// ValidateToken parses and validates a JWT token
func (auth *AuthService) ValidateToken(tokenStr string) (*jwt.Token, error) {
	// Remove the "Bearer " prefix if it exists
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return auth.JWTSecret, nil
	})
}

// JWTInterceptor is a gRPC interceptor that validates JWT tokens, except for the Login method
func (auth *AuthService) JWTInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Bypass JWT validation for the Login method
	if info.FullMethod == "/company.CompanyService/Login" {
		return handler(ctx, req)
	}

	// Extract token from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no metadata in context")
	}

	// Retrieve the "authorization" header
	tokenStr := ""
	if authHeader, exists := md["authorization"]; exists && len(authHeader) > 0 {
		tokenStr = authHeader[0]
	}
	if tokenStr == "" {
		return nil, errors.New("authorization token is missing")
	}

	// Validate the token
	_, err := auth.ValidateToken(tokenStr)
	if err != nil {
		return nil, errors.New("invalid token: " + err.Error())
	}

	// Proceed with the handler if the token is valid
	return handler(ctx, req)
}
