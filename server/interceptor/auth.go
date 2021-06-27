package interceptor

import (
	"context"

	"github.com/MonsterYNH/athena/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func checkAuth(ctx context.Context, jwtKey string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", grpc.Errorf(codes.Unauthenticated, "no auth info")
	}

	tokenArray, ok := md["athena_token"]
	if !ok {
		return "", grpc.Errorf(codes.Unauthenticated, "no auth info")
	}

	if len(tokenArray) == 0 {
		return "", grpc.Errorf(codes.Unauthenticated, "no auth info")
	}

	userID, err := util.ParseToken(tokenArray[0], jwtKey)
	if err != nil {
		return "", grpc.Errorf(codes.Unauthenticated, err.Error())
	}

	return userID, nil
}

func AuthInterceptor(jwtKey string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		userID, err := checkAuth(ctx, jwtKey)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, "user_id", userID)

		return handler(ctx, req)
	}
}
