package models

type Unit struct {
	ID                       int
	Code                     string
	Description              string
	BaseUnit                 string
	OperationValue           float32
	UnitValue                float32
	TransactionCode          string
	ResponsibilityCenterUuid int
	Value                    float32
}
