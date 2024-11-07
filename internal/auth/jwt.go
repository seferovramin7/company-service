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

type AuthService struct {
	JWTSecret []byte
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{JWTSecret: []byte(secret)}
}

func (auth *AuthService) GenerateToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString(auth.JWTSecret)
}

func (auth *AuthService) ValidateToken(tokenStr string) (*jwt.Token, error) {

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return auth.JWTSecret, nil
	})
}

func (auth *AuthService) JWTInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	if info.FullMethod == "/company.CompanyService/Login" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no metadata in context")
	}

	tokenStr := ""
	if authHeader, exists := md["authorization"]; exists && len(authHeader) > 0 {
		tokenStr = authHeader[0]
	}
	if tokenStr == "" {
		return nil, errors.New("authorization token is missing")
	}

	_, err := auth.ValidateToken(tokenStr)
	if err != nil {
		return nil, errors.New("invalid token: " + err.Error())
	}

	return handler(ctx, req)
}
