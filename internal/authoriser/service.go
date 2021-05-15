package authoriser

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/rank-analyse/internal/auth"
)

type Service struct {
	authClient *auth.Client
}

// NewService instantiates a new service
func NewService(authClient *auth.Client) *Service {
	return &Service{
		authClient: authClient,
	}
}

func (s *Service) Authorise(ctx context.Context, request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	authToken := request.AuthorizationToken
	accessToken := strings.Split(authToken, " ")[1]

	_, err := s.authClient.GetClaims(accessToken)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("Unauthorized")

	}

	context := map[string]interface{}{}
	return generatePolicy("user", "Allow", request.MethodArn, context), nil
}

func generatePolicy(principalID, effect string, resource string, context map[string]interface{}) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	authResponse.Context = context
	return authResponse
}
