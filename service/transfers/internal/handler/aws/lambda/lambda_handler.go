package lambda_handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	treegrid_model "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/model/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/service"
	"github.com/aws/aws-lambda-go/events"
)

type LambdaHandler struct {
	transferService service.TransferService
}

func NewLambdaHandler(transferService service.TransferService) *LambdaHandler {
	return &LambdaHandler{transferService: transferService}
}

func (h *LambdaHandler) Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	// eventResp := events.APIGatewayProxyResponse{
	// 	StatusCode: 200,
	// 	MultiValueHeaders: map[string][]string{
	// 		"Access-Control-Allow-Headers": {"Origin", "X-Requested-With", "Content-Type", "Accept"},
	// 		"Content-type":                 {"application/json", "application/x-www-form-urlencoded"},
	// 	},
	// 	Headers: map[string]string{
	// 		"Access-Control-Allow-Origin": "*",
	// 	},
	// 	IsBase64Encoded: false,
	// }
	return h.routeRequest(ctx, request)
}

func (h *LambdaHandler) routeRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	if strings.Contains(request.RawPath, "/page") {
		return h.getTransferPage(ctx, request)
	}

	return h.getTransfersData(ctx, request)
}

func (h *LambdaHandler) getTransfersData(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	start := time.Now()

	reqBody, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		log.Fatal(err)
	}

	urlQuery, err := url.ParseQuery(string(reqBody))
	if err != nil {
		log.Fatalln(err)
	}
	if len(urlQuery["Data"]) == 0 {
		log.Fatalln(urlQuery)
	}

	req, err := treegrid_model.ParseRequest([]byte(urlQuery["Data"][0]))
	if err != nil {
		log.Fatalln(err)
	}

	trGrid, err := treegrid_model.NewTreegrid(req)
	if err != nil {
		log.Fatalln(err)
	}

	allPages := h.transferService.GetPagesCount(trGrid)

	// set this to allow Ajax requests from other origins
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Content-type":                 "application/json",
		"Tcalc":                        fmt.Sprintf("%d ms", time.Since(start).Milliseconds())}

	response, _ := json.Marshal((map[string]interface{}{
		"Body": []string{`#@@@` + fmt.Sprintf("%v", allPages)},
	}))

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         headers,
		Body:            string(response),
		IsBase64Encoded: false,
	}, nil
}

// get Transfers [parents or children of a row]
func (h *LambdaHandler) getTransferPage(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	start := time.Now()
	reqBody, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		log.Fatal(err)
	}

	urlQuery, err := url.ParseQuery(string(reqBody))
	if err != nil {
		log.Fatalln(err)
	}
	if len(urlQuery["Data"]) == 0 {
		log.Fatalln(urlQuery)
	}

	req, err := treegrid_model.ParseRequest([]byte(urlQuery["Data"][0]))
	if err != nil {
		log.Fatalln(err)
	}

	trGrid, err := treegrid_model.NewTreegrid(req)
	if err != nil {
		log.Fatalln(err)
	}

	var response = make([]map[string]string, 0, 100)

	response, err = h.transferService.GetTransfersPageData(trGrid)
	if err != nil {
		log.Println("err", err)
	}

	addData := [][]map[string]string{}
	addData = append(addData, response)

	// set this to allow Ajax requests from other origins
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Content-type":                 "application/json",
		"Tcalc":                        fmt.Sprintf("%d ms", time.Since(start).Milliseconds())}

	result, _ := json.Marshal(map[string][][]map[string]string{
		"Body": addData,
	})

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         headers,
		Body:            string(result),
		IsBase64Encoded: false,
	}, nil
}
