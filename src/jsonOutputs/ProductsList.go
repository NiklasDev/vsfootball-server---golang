package jsonOutputs

// used for record the each product info for user.
type ProductInList struct {
	ProductName string  //Play name
	ProductId string // eg. VsFPlay17
	Price string
	Type string  //
	AppleId string
	Purchased int  //0 - no  1- yes
}

type ProductsListOutput struct {
	Success    string
	Message    string
	ProductInfo   []ProductInList
	Statuscode int
}