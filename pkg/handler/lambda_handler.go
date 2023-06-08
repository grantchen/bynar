package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/aws/scope"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"github.com/aws/aws-lambda-go/events"
)

type LambdaTreeGridPaths struct {
	PathPageCount string
	PathPageData  string
	PathUpload    string
	PathCell      string
}

type LambdaTreeGridHandler struct {
	CallbackGetPageCountFunc CallBackGetPageCount
	CallbackGetPageDataFunc  CallBackGetPageData
	CallbackUploadDataFunc   CallBackLambdaUploadData
	CallBackGetCellDataFunc  CallBackGetCellData
	LambdaPaths              *LambdaTreeGridPaths
}

var authorizationHeader = "Authorization"

func (l *LambdaTreeGridHandler) Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	return l.routeRequest(ctx, request)
}

func (l *LambdaTreeGridHandler) routeRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	switch request.RawPath {
	case l.LambdaPaths.PathPageData:
		return l.getPageData(ctx, request)
	case l.LambdaPaths.PathPageCount:
		return l.getPageCount(ctx, request)
	case l.LambdaPaths.PathUpload:
		return l.handleUpload(ctx, request)
	case l.LambdaPaths.PathCell:
		return l.getCellData(ctx, request)

	}
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Content-type":                 "application/json",
	}
	return events.APIGatewayProxyResponse{
		StatusCode:      400,
		Headers:         headers,
		Body:            "bad request!",
		IsBase64Encoded: false,
	}, nil

}

func (l *LambdaTreeGridHandler) getPageData(ctx context.Context, request events.APIGatewayV2HTTPRequest) (
	events.APIGatewayProxyResponse, error) {
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

	req, err := treegrid.ParseRequest([]byte(urlQuery["Data"][0]))
	if err != nil {
		log.Fatalln(err)
	}

	trGrid, err := treegrid.NewTreegrid(req)
	if err != nil {
		log.Fatalln(err)
	}

	var response = make([]map[string]string, 0, 100)

	response, err = l.CallbackGetPageDataFunc(trGrid)
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

func (l *LambdaTreeGridHandler) getPageCount(ctx context.Context, request events.APIGatewayV2HTTPRequest) (
	events.APIGatewayProxyResponse, error) {
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

	req, err := treegrid.ParseRequest([]byte(urlQuery["Data"][0]))
	if err != nil {
		log.Fatalln(err)
	}

	trGrid, err := treegrid.NewTreegrid(req)
	if err != nil {
		log.Fatalln(err)
	}

	allPages, _ := l.CallbackGetPageCountFunc(trGrid)

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

func (l *LambdaTreeGridHandler) handleUpload(ctx context.Context, request events.APIGatewayV2HTTPRequest) (
	events.APIGatewayProxyResponse, error) {
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

	token, exists := request.Headers[authorizationHeader]
	if !exists {
		eventResp.StatusCode = http.StatusUnauthorized
		eventResp.Body = "unauthorized: token is not provided"

		return eventResp, nil
	}

	// requestScope, err := scope.ResolveFromToken(token)
	requestScope, err := scope.ResolveFromToken(token)
	if err != nil {
		log.Println(err)

		eventResp.StatusCode = http.StatusUnauthorized
		eventResp.Body = "unauthorized: invalid scopes provided"

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

	resp, err := l.CallbackUploadDataFunc(&requestScope, postData)
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

func (l *LambdaTreeGridHandler) getCellData(ctx context.Context, request events.APIGatewayV2HTTPRequest) (
	events.APIGatewayProxyResponse, error) {
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

	req, err := treegrid.ParseRequest([]byte(urlQuery["Data"][0]))
	if err != nil {
		log.Fatalln(err)
	}

	trGrid, err := treegrid.NewTreegrid(req)
	if err != nil {
		log.Fatalln(err)
	}

	response, err := l.CallBackGetCellDataFunc(ctx, trGrid)
	if err != nil {
		log.Println("err", err)
	}

	// set this to allow Ajax requests from other origins
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		"Content-type":                 "application/json",
		"Tcalc":                        fmt.Sprintf("%d ms", time.Since(start).Milliseconds())}

	respBytes, err := json.Marshal(response)
	if err != nil {
		log.Println("err", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         headers,
		Body:            string(respBytes),
		IsBase64Encoded: false,
	}, nil
}
