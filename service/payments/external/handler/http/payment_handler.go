package http_handler

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/service/upload/factory"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

var (
	ModuleID  int = 8
	AccountID int = 123456
)

func NewPaymentHTTPHandler(appConfig config.AppConfig) *handler.HTTPTreeGridHandler {
	dbConnString := appConfig.GetDBConnection()

	if _, err := sql_db.NewConnection(dbConnString); err != nil {
		log.Fatalln("new connection", err)
	}

	// uploadHandler := &http_handler.UploadHandler{ModuleID: ModuleID, AccountID: AccountID}
	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc: func(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
			conn := sql_db.Conn()
			// because need create each conn string per req
			uploadSvc, err := factory.NewUploadService(conn, ModuleID, AccountID)
			if err != nil {
				return &treegrid.PostResponse{
					IO: struct {
						Message string
						Result  int
					}{Message: "could not open connection",
						Result: http.StatusInternalServerError},
				}, nil
			}

			return uploadSvc.Handle(req)
		},
	}
	return handler
}
