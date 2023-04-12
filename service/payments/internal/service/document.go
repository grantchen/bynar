package svc

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
)

type (
	DocumentStorage interface {
		GetDocument(docID int) (models.Document, error)
		GetDocumentSeries(seriesID int) (models.DocumentSeries, error)
		GetDocumentSeriesItem(seriesID int) (models.DocumentSeriesItem, error)
		UpdateDocumentSeriesItem(tx *sql.Tx, item models.DocumentSeriesItem) error
		UpdateDocNumber(tx *sql.Tx, id int, docNumber string) error
	}

	documentService struct {
		store DocumentStorage
	}
)

func NewDocumentService(store DocumentStorage) DocumentService {
	return &documentService{store: store}
}

func (d *documentService) Handle(tx *sql.Tx, modelID, docID int, docNo string) error {
	doc, err := d.store.GetDocument(docID)
	if err != nil {
		return fmt.Errorf("get document: [%w]", err)
	}

	if doc.Status == 0 {
		return fmt.Errorf("doc status is 0, doc closed")
	}

	docSeries, err := d.store.GetDocumentSeries(doc.SeriesID)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("no document series with id ", doc.SeriesID)

			return nil
		}

		return fmt.Errorf("get document series: [%w]", err)
	}

	// no auto generating
	if docSeries.DefaultNos == 0 {
		logger.Debug("auto default no is off")

		return nil
	}

	if docSeries.ManualNos == 1 {
		if docNo != "" {
			logger.Debug("manual doc number and setted by user", docNo)

			return nil
		}

		return d.generate(tx, modelID, doc.SeriesID, false)
	}

	return d.generate(tx, modelID, doc.SeriesID, true)
}

func (d *documentService) generate(tx *sql.Tx, transferID int, seriesID int, update bool) error {
	seriesItem, err := d.store.GetDocumentSeriesItem(seriesID)
	if err != nil {
		return fmt.Errorf("get document series item: [%w]", err)
	}

	if seriesItem.Open == 0 {
		return fmt.Errorf("document series is full: series item id : %d", seriesItem.ID)
	}

	seriesItem = d.getNewDocumentNo(seriesItem)
	logger.Debug("new series", seriesItem)

	logger.Debug("update transfer doc number", seriesItem.LastNoUsed)

	if err := d.store.UpdateDocNumber(tx, transferID, seriesItem.LastNoUsed); err != nil {
		return fmt.Errorf("update transfer document_no: [%w]", err)
	}

	if update {
		logger.Debug("update series item")

		return d.store.UpdateDocumentSeriesItem(tx, seriesItem)
	}

	return nil
}

func (d *documentService) getNewDocumentNo(seriesItem models.DocumentSeriesItem) models.DocumentSeriesItem {
	const (
		numberLength = 4
	)

	logger.Debug(seriesItem.LastNoUsed)

	if seriesItem.LastNoUsed == "" {
		seriesItem.LastNoUsed = seriesItem.StartingNo
		seriesItem.LastDateUsed = time.Now().Format("2006-01-02") //2022-09-20

		return seriesItem
	}

	abbr, serNumer := getAbbrNumber(seriesItem.LastNoUsed)
	serNumer += seriesItem.IncrementNo

	_, maxNumer := getAbbrNumber(seriesItem.EndingNo)
	logger.Debug("new number vals", abbr, serNumer, "max number", maxNumer)

	if serNumer > maxNumer {
		logger.Debug("new series number lager than max, close series", serNumer, maxNumer)

		seriesItem.Open = 0
	}

	serNumStr := strconv.Itoa(serNumer)
	serNumPref := strings.Repeat("0", numberLength-len(serNumStr))

	seriesItem.LastNoUsed = abbr + serNumPref + serNumStr
	seriesItem.LastDateUsed = time.Now().Format("2006-01-02") //2022-09-20

	return seriesItem
}

func getAbbrNumber(str string) (abbr string, numb int) {
	var (
		abbrInd int
		err     error
	)
	for k, v := range str {
		if v >= '0' && v <= '9' {
			abbrInd = k
			break
		}
	}

	abbr = str[:abbrInd]
	numStr := str[abbrInd:]

	numb, err = strconv.Atoi(numStr)
	if err != nil {
		logger.Debug("invalid doc number", abbr, numStr, numb)
	}

	return
}
