package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/postmannen/fakt/data"
)

//separateStrNumForButton , takes *http.Request as input, and returns string and int
func separateStrNumForButton(r *http.Request) (string, int) {
	var buttonNumbers string
	var buttonValue string
	var num int
	var err error
	for k, v := range r.Form {
		//fmt.Println("---VERDIER ---- ", k, " : ", v)
		re := regexp.MustCompile("[a-zA-Z]+")
		buttonLetters := re.FindString(k)
		re = regexp.MustCompile("[0-9]+")
		if buttonLetters == "billLineAddButton" {
			buttonValue = v[0]
			buttonNumbers = re.FindString(k)
			num, err = strconv.Atoi(buttonNumbers)
			if err != nil {
				log.Printf("ERROR: strconv.Atoi : %v", err)
			}
		}
		if buttonLetters == "billLineDeleteButton" {
			buttonValue = v[0]
			buttonNumbers = re.FindString(k)
			num, err = strconv.Atoi(buttonNumbers)
			if err != nil {
				log.Printf("ERROR: strconv.Atoi : %v", err)
			}
		}
		if buttonLetters == "billLineModifyButton" {
			buttonValue = v[0] //value is a slice of strings, get the first value
			buttonNumbers = re.FindString(k)
			num, err = strconv.Atoi(buttonNumbers)
			if err != nil {
				log.Printf("ERROR: strconv.Atoi : %v", err)
			}
		}
	}

	return buttonValue, num
}

//sortBills sorts the bills so the last bill_id is first in the slice
func sortBills(bills []data.Bill) []data.Bill {
	for i := 0; i < len(bills); i++ {
		for ii := 0; ii < len(bills); ii++ {
			if bills[i].BillID > bills[ii].BillID {
				tmp := bills[ii]
				bills[ii] = bills[i]
				bills[i] = tmp
			}
		}
	}
	return bills
}

//Finds all the numbers used in html names in form
//used to get all the unique bill lines
func findBillLineNumbersInForm(r *http.Request) (numbers []int) {
	for k, v := range r.Form {
		fmt.Printf("--- k = %v of type %T , and v = %v of type %T\n", k, k, v, v)
		reLetters := regexp.MustCompile("[a-zA-Z]+")
		reNum := regexp.MustCompile("[0-9]+")
		letterPart := reLetters.FindString(k)
		numberStr := reNum.FindString(k)
		numberPart, _ := strconv.Atoi(numberStr)
		log.Printf("-----letterPart = %v, and numberPart = %v\n", letterPart, numberPart)

		found := false
		//check if numberPart is allready in the numbers slice, if NOT.....add it
		for _, vv := range numbers {
			//fmt.Printf("***trying to compare vv=%v and numberPart=%v \n", vv, numberPart)
			if numberPart == vv {
				found = true
				fmt.Println("The numbers are equal")
			}
		}
		if !found {
			numbers = append(numbers, numberPart)
		}
	}
	return numbers
}

//range billLines to compare storedLines.X with formBillLines.X to see if any values have changed. Return changed = true if changed
func checkIfBillLineChanged(lineNRs []int, storedLines []data.BillLines, formLines []data.BillLines) (changed bool) {
	for _, num := range lineNRs {
		for _, line := range storedLines {
			if line.LineID == num {
				for _, line2 := range formLines {
					if line2.LineID == num {
						if line.LineID != line2.LineID {
							changed = true
						}
						if line.Description != line2.Description {
							changed = true
						}
						if line.Quantity != line2.Quantity {
							changed = true
						}
						if line.DiscountPercentage != line2.DiscountPercentage {
							changed = true
						}
						if line.VatUsed != line2.VatUsed {
							changed = true
						}
						if line.PriceExVat != line2.PriceExVat {
							changed = true
						}
					}
				}
			}
		}
	}
	return changed
}

//getBillLineFormValues parses all the data in the form, compares them with the current billID, and returns
//a slice with all the values entered in the form.
//all fields and buttons in the form have name values postfixed with the {{.LineID}}, so this function
//separates the first part of the name and the {{.LineID}} to know what fields to update
func getBillLineFormValues(lineNumbers []int, r *http.Request, billID int) (formBillLines []data.BillLines) {
	var tempLines data.BillLines

	for _, num := range lineNumbers {
		fmt.Println("-$- Outerloop, num of lineNumbers = ", num)
		//iterate all the data in form
		for k, v := range r.Form {

			//split out the letter and number part of button name
			reLetters := regexp.MustCompile("[a-zA-Z]+")
			reNum := regexp.MustCompile("[0-9]+")
			letterPart := reLetters.FindString(k)    //get name "billLineModifyButton"
			numberStr := reNum.FindString(k)         //get the line nr.that the button belonged to. Nr is postfixed in the name
			numberPart, _ := strconv.Atoi(numberStr) //convert the nr got`en from form to int, so it can be used later

			if num == numberPart {
				tempLines.BillID = billID
				if letterPart == "billLineID" {
					myVal, err := strconv.Atoi(v[0])
					if err != nil {
						log.Println("ERROR: strconv billLineID : ", err)
					}
					tempLines.LineID = myVal
					fmt.Printf("--- templLines.BillID er satt til %v\n", tempLines.BillID)
				}
				if letterPart == "billLineDescription" {

					tempLines.Description = v[0]
					fmt.Printf("--- templLines.Description er satt til %v\n", v[0])
				}
				if letterPart == "billLineQuantity" {
					myVal, err := strconv.Atoi(v[0])
					if err != nil {
						log.Println("ERROR: strconv billLineQuantity : ", err)
					}
					tempLines.Quantity = myVal
					fmt.Printf("--- templLines.Quantity er satt til %v\n", tempLines.Quantity)
				}
				if letterPart == "billLineDiscountPercentage" {
					myVal, err := strconv.Atoi(v[0])
					if err != nil {
						log.Println("ERROR: strconv billLineDiscountPercentage : ", err)
					}
					tempLines.DiscountPercentage = myVal
					fmt.Printf("--- templLines.DiscountPercentage er satt til %v\n", tempLines.DiscountPercentage)
				}
				if letterPart == "billLineVatUsed" {
					myVal, err := strconv.Atoi(v[0])
					if err != nil {
						log.Println("ERROR: strconv billLineVatUsed : ", err)
					}
					tempLines.VatUsed = myVal
					fmt.Printf("--- templLines.VatUsed er satt til %v\n", tempLines.VatUsed)
				}
				if letterPart == "billLinePriceExVat" {
					myVal, err := strconv.ParseFloat(v[0], 64)
					if err != nil {
						log.Println("ERROR: strconv billLinePriceExVat : ", err)
					}
					tempLines.PriceExVat = myVal
				}
			}
		}
		formBillLines = append(formBillLines, tempLines)
	}
	return
}

//updateBillTotalExVat updates the bill field total price ex vat,
//also writes the update info to correct field in db
func updateBillTotalExVat(bill *data.Bill, billID int, billLines []data.BillLines) {
	//TODO: Fix so the total values are made from line total to avoid doing the same calculations in several functions
	bill.TotalExVat = 0
	for _, v := range billLines {
		v.PriceExVat *= float64(v.Quantity)
		v.PriceExVat -= v.PriceExVat / 100 * float64(v.DiscountPercentage)
		bill.TotalExVat += v.PriceExVat
	}

	/* NOTE : Removing db db write, will put in a write for all fields when all data is processed
	//add the TotalExVat to db here
	if bill.TotalExVat != 0 {
		data.UpdateBillPriceExVat(pDB, bill.TotalExVat, bill.BillID)
	}
	*/
}

//updateBillTotalIncVat updates the bill field total price ex vat,
//also writes the update info to correct field in db
func updateBillTotalIncVat(bill *data.Bill, billLines []data.BillLines) {
	//TODO: Fix so the total values are made from line total to avoid doing the same calculations in several functions
	var lineIncVat float64
	bill.TotalIncVat = 0
	for _, v := range billLines {
		if v.VatUsed == 0 {
			lineIncVat = v.PriceExVat
			lineIncVat -= lineIncVat / 100 * float64(v.DiscountPercentage)
		} else {
			lineIncVat = v.PriceExVat + (v.PriceExVat / 100.0 * float64(v.VatUsed))
			lineIncVat -= lineIncVat / 100 * float64(v.DiscountPercentage)
		}
		lineIncVat *= float64(v.Quantity)
		bill.TotalIncVat += lineIncVat

	}
}

//updateLineExVatTotal updates the struct data for the total ex vat pr. line
func updateLineExVatTotal(b []data.BillLines) {
	for i := 0; i < len(b); i++ {
		sum := b[i].PriceExVat * float64(b[i].Quantity)
		sum = sum - (sum / 100 * float64(b[i].DiscountPercentage))
		b[i].PriceExVatTotal = sum
	}
}
