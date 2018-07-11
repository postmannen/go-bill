package data

//User is used for all customers and users
type User struct { //some
	Number         int
	FirstName      string `schema:"firstname"`
	LastName       string `schema:"lastname"`
	Mail           string `schema:"mail"`
	Address        string `schema:"address"`
	PostNrAndPlace string `schema:"poAddr"`
	PhoneNr        string `schema:"phone"`
	OrgNr          string `schema:"orgNr"`
	CountryID      string `schema:"countryId"`
	Selected       string
	BankAccount    string `schema:"bankAccount"`
}

type ape struct {
	navn string
}

//type hest struct {
//	ape
//}

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
