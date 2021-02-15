// Package grpc implements functions to expose project service endpoint using GRPC protocol.
package grpc

import (
	"context"

	"github.com/decentralized-cloud/project/models"
	"github.com/go-kit/kit/endpoint"
	"github.com/micro-business/go-core/jwt/grpc"
)

func (service *transportService) createAuthMiddleware(endpointName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			token, err := grpc.ParseAndVerifyToken(ctx, service.jwksURL, true)
			if err != nil {
				return nil, err
			}

			parsedToken := models.ParsedToken{Email: token.PrivateClaims()["email"].(string)}
			ctx = context.WithValue(ctx, models.ContextKeyParsedToken, parsedToken)

			return next(ctx, request)
		}
	}
}
