package models

type InboundFlow struct {
	ID               int
	ModuleID         int
	ItemID           int
	ParentID         int
	LocationID       int
	PostingDate      string
	Quantity         float32
	Value            float32
	OutboundQuantity float32
	OutboundValue    float32
	Status           int
}

type OutboundFlow struct {
	ID            int
	ModuleID      int
	LocationID    int
	ItemID        int
	ParentID      int
	TransactionNo int
	PostingDate   string
	Quantity      float32
	ValueAvco     float32
	ValueFifo     float32
}
