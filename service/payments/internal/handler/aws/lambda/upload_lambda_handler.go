package lambda_handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/scope"
	svc_upload "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/service/upload"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/service/upload/factory"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/aws/secretsmanager"
	sql_connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"github.com/aws/aws-lambda-go/events"
)

type LambdaHandler struct {
	secretsManager      secretsmanager.SecretsManager
	connectionPool      ConnectionResolver
	connectionStringKey string
	uploadSvc           svc_upload.UploadService
}

var (
	// ModuleID is hardcoded as provided in the specification.
	ModuleID = 8

	connectionStringKey = "bynar"
	awsRegion           = "eu-central-1"

	authorizationHeader = "Authorization"
)

func NewLambdaHandler(secretsManager secretsmanager.SecretsManager, connectionPool ConnectionResolver) *LambdaHandler {
	return &LambdaHandler{secretsManager: secretsManager, connectionPool: connectionPool}
}

func (h *LambdaHandler) Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	eventResp := events.APIGatewayProxyResponse{
		StatusCode: 200,
		MultiValueHeaders: map[string][]string{
			"Access-Control-Allow-Headers": {"Origin", "X-Requested-With", "Content-Type", "Accept"},
			"Content-type":                 {"application/json", "application/x-www-form-urlencoded"},
		},
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		IsBase64Encoded: false,
	}

	secret, err := h.secretsManager.GetString(connectionStringKey)
	if err != nil {
		log.Println(err)

		eventResp.Body = err.Error()

		return eventResp, nil
	}

	token, exists := request.Headers[authorizationHeader]
	if !exists {
		eventResp.StatusCode = http.StatusUnauthorized
		eventResp.Body = "unauthorized: token is not provided"

		return eventResp, nil
	}

	requestScope, err := scope.ResolveFromToken(token)
	if err != nil {
		log.Println(err)

		eventResp.StatusCode = http.StatusUnauthorized
		eventResp.Body = "unauthorized: invalid scopes provided"

		return eventResp, nil
	}

	conn, err := h.connectionPool.Get(sql_connection.ChangeDatabaseConnectionSchema(secret, strconv.Itoa(requestScope.OrganizationID)))
	if err != nil {
		log.Println(err)

		eventResp.StatusCode = http.StatusInternalServerError
		eventResp.Body = "could not open connection"

		return eventResp, nil
	}

	urlQuery, err := url.ParseQuery(request.Body)
	if err != nil {
		log.Println(err)

		eventResp.Body = err.Error()

		return eventResp, nil
	}

	if len(urlQuery["Data"]) == 0 {
		eventResp.Body = errors.New("empty request").Error()

		return eventResp, nil
	}

	postData := &treegrid.PostRequest{
		Changes: make([]map[string]interface{}, 10),
	}
	if unmarshalErr := json.Unmarshal([]byte(urlQuery["Data"][0]), &postData); unmarshalErr != nil {
		log.Println(unmarshalErr)

		eventResp.Body = errors.New("empty request").Error()

		return eventResp, nil
	}

	uploadSvc, err := factory.NewUploadService(conn, ModuleID, requestScope.AccountID)
	if err != nil {
		log.Println(err)

		eventResp.StatusCode = http.StatusInternalServerError
		eventResp.Body = "Problem while initializing uploader service. Please check logs."

		return eventResp, err
	}

	resp, err := uploadSvc.Handle(postData)
	if err != nil {
		log.Println(err)

		eventResp.Body = err.Error()

		return eventResp, err
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		eventResp.Body = err.Error()

		return eventResp, nil
	}

	eventResp.Body = string(respBytes)

	return eventResp, nil
}
