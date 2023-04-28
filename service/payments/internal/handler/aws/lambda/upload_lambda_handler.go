package lambda_handler

import (
	"log"
	"net/http"
	"strconv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/service/upload/factory"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/aws/scope"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/aws/secretsmanager"
	sql_connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type LambdaHandler struct {
	secretsManager secretsmanager.SecretsManager
	connectionPool ConnectionResolver
}

var (
	// ModuleID is hardcoded as provided in the specification.
	ModuleID = 8

	connectionStringKey = "bynar"
)

func NewLambdaHandler(secretsManager secretsmanager.SecretsManager, connectionPool ConnectionResolver) *LambdaHandler {
	return &LambdaHandler{secretsManager: secretsManager, connectionPool: connectionPool}
}

func (h *LambdaHandler) Handle(requestScope *scope.RequestScope, postData *treegrid.PostRequest) (
	*treegrid.PostResponse, error) {
	resp := &treegrid.PostResponse{}
	resp.IO.Result = 200

	secret, err := h.secretsManager.GetString(connectionStringKey)
	if err != nil {
		log.Println(err)
		resp.IO.Message = err.Error()
		return resp, err
	}

	conn, err := h.connectionPool.Get(sql_connection.ChangeDatabaseConnectionSchema(secret, strconv.Itoa(requestScope.OrganizationID)))
	if err != nil {
		log.Println(err)

		resp.IO.Result = http.StatusInternalServerError
		resp.IO.Message = "could not open connection"

		return resp, err
	}

	uploadSvc, err := factory.NewUploadService(conn, ModuleID, requestScope.AccountID)
	if err != nil {
		log.Println(err)

		resp.IO.Result = http.StatusInternalServerError
		resp.IO.Message = "Problem while initializing uploader service. Please check logs."

		return resp, err
	}

	resp, err = uploadSvc.Handle(postData)
	return resp, err
}
