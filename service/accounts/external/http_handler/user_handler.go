package http_handler

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
)

func NewUserHTTPHandler() *handler.HTTPTreeGridHandler {

	userSvc := service.NewUserService()

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:   userSvc.Handle,
		CallbackGetPageDataFunc:  userSvc.GetPageData,
		CallbackGetPageCountFunc: userSvc.GetPageCount,
	}
	return handler

}
