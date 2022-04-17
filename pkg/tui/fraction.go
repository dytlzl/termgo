package tui

type Fraction struct {
	// Numerator
	Numer int
	// Denominator
	Denom int
}

func NewFraction(numer, denom int) *Fraction {
	return &Fraction{
		Numer: numer,
		Denom: denom,
	}
}
