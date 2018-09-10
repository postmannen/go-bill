package main

//User is used for all customers and users
type User struct { //some
	Number         int
	FirstName      string
	LastName       string
	Mail           string
	Address        string
	PostNrAndPlace string
	PhoneNr        string
	OrgNr          string
	CountryID      string
	Selected       string
	BankAccount    string
}

type ape struct {
	navn string
}

type hest struct {
	ape
}

//Bill struct specifications
type Bill struct {
	BillID      int
	UserID      int
	CreatedDate string
	DueDate     string
	Comment     string
	TotalExVat  float64
	TotalIncVat float64
	TotalVat    float64
	Paid        int
}

//BillLines struct. Fields must be export (starting Capital letter) to be passed to template
type BillLines struct {
	BillID             int
	LineID             int
	ItemID             int
	Description        string
	Quantity           int
	DiscountPercentage int
	VatUsed            int
	PriceExVat         float64
	PriceExVatTotal    float64
	//just create some linenumbers for testing
}
