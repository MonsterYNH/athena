package services

import (
	"athena/api/v1/auth2"
	"athena/orm"
	"athena/util"
	"context"
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Auth2Service struct {
	auth2.UnimplementedAuth2SerivceServer
}

func (services *Auth2Service) Auth(ctx context.Context, request *auth2.Auth2AuthRequest) (*auth2.Auth2AuthResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no token info")
	}

	values := md.Get("athena_token")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "empty token info")
	}
	userID, err := ParseToken(values[0], "jwt_key")

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	response := &auth2.Auth2AuthResponse{Service: request.Service}
	tokenStr, err := util.GenerateToken(userID, "jwt_key", time.Hour*12)
	if err == nil {
		response.Status = true
	}

	if err := grpc.SendHeader(ctx, metadata.Pairs("athena_token", tokenStr)); err != nil {
		return nil, err
	}

	return response, nil
}

func (services *Auth2Service) Login(ctx context.Context, request *auth2.Auth2LoginRequest) (*auth2.Auth2LoginResponse, error) {
	if len(request.Account) == 0 || len(request.Password) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid account or password")
	}
	db := orm.GetDB()

	user := &orm.User{}

	if err := db.Debug().Where("account = ? and password = ?", request.Account, request.Password).First(user).Error; err != nil {
		return nil, err
	}

	tokenStr, err := GenerateToken(*user, "jwt_key", time.Hour*12)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if err := grpc.SendHeader(ctx, metadata.Pairs("athena_token", tokenStr)); err != nil {
		return nil, err
	}

	return &auth2.Auth2LoginResponse{
		Token:   tokenStr,
		Account: user.Account,
		Name:    user.Name,
		Status:  true,
	}, nil
}

func (services *Auth2Service) Regist(ctx context.Context, request *auth2.Auth2RegistRequest) (*auth2.Auth2RegistResponse, error) {
	db := orm.GetDB()

	user := &orm.User{
		Account:  request.Account,
		Password: request.Password,
	}
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	tokenStr, err := GenerateToken(*user, "jwt_key", time.Hour*12)
	if err != nil {
		return nil, err
	}

	if err := grpc.SendHeader(ctx, metadata.Pairs("athena_token", tokenStr)); err != nil {
		return nil, err
	}

	return &auth2.Auth2RegistResponse{
		Status: true,
	}, nil
}

func GenerateToken(user orm.User, jwtKey string, expireDuration time.Duration) (string, error) {
	expire := time.Now().Add(expireDuration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.UnixNano(),
		},
	})

	return token.SignedString([]byte(jwtKey))
}

func ParseToken(tokenStr, jwtKey string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(tk *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("Unauthorized")
	}

	return claims.ID, nil
}

func ParseTokenWithUser(tokenStr, jwtKey string) (*orm.User, error) {
	id, err := ParseToken(tokenStr, jwtKey)
	if err != nil {
		return nil, err
	}

	db := orm.GetDB()

	user := &orm.User{}
	return user, db.First(user, "id = ?", id).Error
}

type Claims struct {
	ID string
	jwt.StandardClaims
}
