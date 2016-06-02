package models

type Play struct {
	Id          int64
	Title       string
	Image       string
	Type        string
	Possession  int
	Premium     int
	Productid   string
	Producttype int
	Price       float64
	Canflip     int
	Version     int
	Filesize    float64
	Description string
}
