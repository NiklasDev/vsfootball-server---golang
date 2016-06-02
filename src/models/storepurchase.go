package models

type Storepurchase struct {
	Id        int64
	Userid    string
	Productid string
	Playid    int64
	Receipt   string
	Os        string
}
