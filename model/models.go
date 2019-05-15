package model

type Concentrate struct {
	ID       string
	Name     string
	VendorID string
	Gravity  *float64
	URLIDs   []*string
}

type Vendor struct {
	ID     string
	Name   string
	Code   string
	URLIDs []*string
}
