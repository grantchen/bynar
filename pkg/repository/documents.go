package repository

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/models"
)

type documentRepository struct {
	conn            *sql.DB
	updateTableName string
}

func NewDocuments(conn *sql.DB, updateTableName string) DocumentRepository {
	return &documentRepository{conn: conn, updateTableName: updateTableName}
}

func (d *documentRepository) GetDocument(docID int) (m models.Document, err error) {
	query := `
	SELECT document_type,document_abbrevation,series_id,status	
	FROM documents
	WHERE id = ?
	`

	err = d.conn.QueryRow(query, docID).Scan(&m.DocumentType, &m.DocumentAbbrevation, &m.SeriesID, &m.Status)

	return
}

func (d *documentRepository) GetDocumentSeries(seriesID int) (m models.DocumentSeries, err error) {
	logger.Debug("get document series", seriesID)

	query := `
	SELECT code, description, default_nos, manual_nos
	FROM series
	WHERE id = ?
	`

	err = d.conn.QueryRow(query, seriesID).Scan(&m.Code, &m.Description, &m.DefaultNos, &m.ManualNos)

	return
}

func (d *documentRepository) GetDocumentSeriesItem(seriesID int) (m models.DocumentSeriesItem, err error) {
	logger.Debug("get document series item", seriesID)

	query := `
	SELECT id, parent_id, starting_no, increment_no, last_date_used, last_no_used, ending_no, open 	
	FROM series_items
	WHERE parent_id = ? AND open = 1
	`

	err = d.conn.QueryRow(query, seriesID).Scan(&m.ID, &m.ParentID, &m.StartingNo, &m.IncrementNo, &m.LastDateUsed, &m.LastNoUsed, &m.EndingNo, &m.Open)

	return
}

func (d *documentRepository) UpdateDocumentSeriesItem(tx *sql.Tx, item models.DocumentSeriesItem) (err error) {
	logger.Debug("UpdateDocumentSeriesItem", item)

	query := `
	UPDATE series_items
	SET last_date_used = ?, last_no_used = ?,  open = ? 	
	WHERE id = ?
	`
	_, err = tx.Exec(query, item.LastDateUsed, item.LastNoUsed, item.Open, item.ID)

	return
}

func (d *documentRepository) UpdateDocNumber(tx *sql.Tx, id int, docNumber string) (err error) {
	logger.Debug("UpdateDocumentSeriesItem", id, docNumber)

	query := `
	UPDATE ` + d.updateTableName + `
	SET document_no = ?
	WHERE id = ?
	`

	_, err = tx.Exec(query, docNumber, id)

	return
}
