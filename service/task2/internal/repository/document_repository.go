package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type documentRepository struct {
}

func NewDocumentRepository(db *sql.DB) DocumentRepository {
	return &documentRepository{}
}

func (dr *documentRepository) IsAuto(tx *sql.Tx, tr *treegrid.MainRow) (bool, error) {
	return false, nil
}

func (dr *documentRepository) Generate(tx *sql.Tx, tr *treegrid.MainRow) (string, error) {
	return "", nil
}

func (dr *documentRepository) Save(tx *sql.Tx, tr *treegrid.MainRow) error {
	query := `
	SELECT s.id, s.default_nos, s.manual_nos
	FROM transfer t
	INNER JOIN documents d ON t.document_id = d.id
	INNER JOIN series s ON s.id = d.series_id
	WHERE t.id = ?
	`
	var seriesID, defNos, manNos int
	err := tx.QueryRow(query, tr.Fields.GetID()).Scan(&seriesID, &defNos, &manNos)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return fmt.Errorf("query row: [%w], query: %s", err, query)
	}

	// without generation
	if defNos == 0 {
		return nil
	}

	//
	if manNos == 1 {
		var doc_number string

		if err := tx.QueryRow(`
		SELECT document_no
		FROM transfer
		WHERE id = ?
		`, tr.Fields.GetID()).Scan(&doc_number); err != nil {
			return fmt.Errorf("query row: [%w], query: %s", err, query)
		}

		if doc_number != "" {
			return nil
		}
	}

	return generate(tx, seriesID, tr.Fields.GetID())
}

type DocSeries struct {
	ID           int
	StartingNo   string
	IncrementNo  int
	LastDateUsed string
	LastUsedNo   string
	EndingNo     string
	Open         int
}

func generate(tx *sql.Tx, seriesID int, transferID interface{}) error {
	query := `
	SELECT id, starting_no, increment_no, last_date_used, last_used_no, ending_no, open
	FROM series_items
	WHERE parent_id = ?
	`

	var docSeries DocSeries

	if err := tx.QueryRow(query, transferID).
		Scan(&docSeries.ID, &docSeries.StartingNo, &docSeries.IncrementNo, &docSeries.LastDateUsed,
			&docSeries.LastUsedNo, &docSeries.EndingNo, &docSeries.Open); err != nil {
		return fmt.Errorf("rows scan: [%w], query: %s", err, query)
	}

	if docSeries.Open != 1 {
		return nil
	}

	if docSeries.LastUsedNo == "" {
		docSeries.LastUsedNo = docSeries.StartingNo

		return updateLastDocument(tx, docSeries)
	}

	noAbr := ""
	noIDStr := ""
	if len(docSeries.LastUsedNo) > 2 {
		noAbr = docSeries.LastUsedNo[:2]
		noIDStr = docSeries.LastUsedNo[2:]
	}

	if noAbr == "" || noIDStr == "" {
		return fmt.Errorf("invalid doc string: %s", docSeries.LastUsedNo)
	}

	noIDInt, err := strconv.Atoi(noIDStr)
	if err != nil {
		return fmt.Errorf("convert str to number: [%w]", err)
	}

	noIDInt += docSeries.IncrementNo
	newIDStr := strconv.Itoa(noIDInt)

	addNulsCount := 4 - len(newIDStr)
	for i := 0; i < addNulsCount; i++ {
		newIDStr = "0" + newIDStr
	}

	docSeries.LastUsedNo = noAbr + newIDStr

	return updateLastDocument(tx, docSeries)
}

func updateLastDocument(tx *sql.Tx, d DocSeries) error {
	query := `
	UPDATE series_items
	SET last_no_used = ?
	WHERE id = ?
	`
	_, err := tx.Exec(query, d.LastUsedNo, d.ID)
	if err != nil {
		return fmt.Errorf("exec: [%w], query: %s", err, query)
	}

	return nil
}
