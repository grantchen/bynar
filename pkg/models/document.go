package models

type Document struct {
	ID                  int
	DocumentType        string
	DocumentAbbrevation string
	Workspace           string
	SeriesID            int
	Status              int
}

type DocumentSeries struct {
	Id          int
	Code        string
	Description string
	DefaultNos  int
	ManualNos   int
}

type DocumentSeriesItem struct {
	ID           int
	ParentID     int
	StartingNo   string
	IncrementNo  int
	LastDateUsed string
	LastNoUsed   string
	WarningNo    string
	EndingNo     string
	Open         int
}
