package models

type DiscountVat struct {
	ID         int
	Value      float32
	Percentage int
}

func (d DiscountVat) Calculate(val float32) (float32, error) {
	if d.Percentage == 0 {
		return d.Value, nil
	}

	return val * d.Value / 100, nil
}

type Currency struct {
	ID           int
	ExchangeRate float32
}

