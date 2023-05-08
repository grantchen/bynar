package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/simplecurd/internal/model"
)

type languageRepostory struct {
	db            *sql.DB
	languageTable string
}

var (
	ErrMissingRequiredParams = errors.New("missing required params")
	ErrAlreadyExist          = errors.New("already exist")
)

// AddNewLanguage implements LanguageRepository
func (l *languageRepostory) AddNewLanguage(newLang *model.Language) (*model.Language, error) {
	// sql query
	stmtStr := `
	INSERT INTO  languages (country, language, two_letters, three_letters, number) 
	VALUES (?,?,?,?,?)
	`
	stmt, err := l.db.Prepare(stmtStr)
	if err != nil {
		return nil, fmt.Errorf("prepare statement: [%w]", err)
	}
	defer stmt.Close()

	rs, err := stmt.Exec(newLang.Country, newLang.Language, newLang.Two_letters, newLang.Three_letters, newLang.Number)
	if err != nil {
		return nil, fmt.Errorf("exec statement: [%w]", err)
	}

	id64, _ := rs.LastInsertId()
	newLang.Id = int(id64)
	return newLang, nil
}

// DeleteLanguage implements LanguageRepository
func (l *languageRepostory) DeleteLanguage(ID int64) error {
	stmt, err := l.db.Prepare("DELETE from  languages where id=?")
	if err != nil {
		return fmt.Errorf("prepare statement: [%w]", err)
	}
	defer stmt.Close()

	id64 := int64(ID)
	if _, err := stmt.Exec(id64); err != nil {
		return fmt.Errorf("exec statement: [%w]", err)
	}

	return nil
}

// UpdateLanguage implements LanguageRepository
func (l *languageRepostory) UpdateLanguage(lang *model.Language) (*model.Language, error) {
	var (
		colNames    []string
		colVals     []interface{}
		inInterface map[string]interface{}
	)

	// fields that need to be updated in db
	excludeFields := map[string]bool{
		"id":            true,
		"Added":         true,
		"Deleted":       true,
		"Changed":       true,
		"country":       false,
		"language":      false,
		"two_letters":   false,
		"three_letters": false,
		"number":        false,
	}

	// json.Marshal and json.Unmarshal - just a mapping struct to map
	inrec, err := json.Marshal(lang)
	if err != nil {
		return nil, fmt.Errorf("marshal lang: [%w]", err)
	}

	if err := json.Unmarshal(inrec, &inInterface); err != nil {
		return nil, fmt.Errorf("unmarshal lang: [%w]", err)
	}

	// fullfilling column names and column values for db query
	for field, val := range inInterface {
		if !excludeFields[field] && val != "" /*&& strings.Contains(jsonString,strings.ToLower(field))*/ {
			colNames = append(colNames, strings.ToLower(field)+" = ?")
			colVals = append(colVals, val)
		}
	}

	colVals = append(colVals, lang.Id)

	// if len=0 then some value empty that indicates bad request
	if len(colNames) == 0 {
		return nil, ErrMissingRequiredParams
	}

	stmt, err := l.db.Prepare("UPDATE  languages SET " + strings.Join(colNames, ", ") + " WHERE id=?")
	if err != nil {
		return nil, fmt.Errorf("prepare statement: [%w]", err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(colVals...); err != nil {
		return nil, fmt.Errorf("exec statement: [%w]", err)
	}

	return lang, nil
}

// GetAllLanguage implements LanguageRepository
func (l *languageRepostory) GetAllLanguage() []*model.Language {
	var languages []*model.Language

	// Query to get all language records
	rows, err := l.db.Query("SELECT * FROM languages")
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var lng *model.Language
		if err := rows.Scan(&lng.Id, &lng.Country, &lng.Language, &lng.Two_letters, &lng.Three_letters, &lng.Number); err != nil {
			log.Fatalln(err)
		}
		languages = append(languages, lng)
	}
	if err := rows.Err(); err != nil {
		log.Fatalln(err)
	}
	return languages
}

// GetCountryAndNumber implements LanguageRepository
func (l *languageRepostory) GetCountryAndNumber(id int) (country string, number int64, err error) {
	query := `
	SELECT country, number FROM  languages 
	WHERE id = ?
	`
	stmt, err := l.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("prepare statement: [%w]", err)

		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&country, &number)
	if err != nil {
		err = fmt.Errorf("prepare statement: [%w]", err)

		return
	}

	return
}

// ValidateOnIntegrity implements LanguageRepository
func (l *languageRepostory) ValidateOnIntegrity(id int, country string, number int64) (bool, error) {
	query := `
	SELECT count(*) from  languages 
	WHERE country = ? AND number = ? AND id != ?
	`
	stmt, err := l.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("prepare statement: [%w]", err)
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(country, number, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("query row: [%w]", err)
	}

	return count == 0, nil
}

func NewLanguageRepository(db *sql.DB, languageTable string) LanguageRepository {
	return &languageRepostory{db: db, languageTable: languageTable}
}
