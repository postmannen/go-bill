package main

import (
	"database/sql"
	"fmt"
	"log"
)

/****************************
*	DATABASE FUNCTIONS		*
*							*
****************************/

//************************ USER SECTION ************************

//Query the database for all users, and return a slice of struct with all users
func queryDBForAllUserInfo(pDB *sql.DB) []User {
	//get total rows in database
	lastUserID, countLines := queryDBForLastCustomerUID(pDB)
	p := []User{}
	fmt.Println("queryDBForAllUserInfo : queryDBForAllUserInfo highestNR ER = ", lastUserID)
	fmt.Println("queryDBForAllUserInfo : queryDBForAllUserInfo countlines = ", countLines)

	for i := 0; i <= lastUserID; i++ {
		//append the row to slice
		pTemp := queryDBForSingleUserInfo(pDB, i)
		if pTemp.Number != 999999999 {
			p = append(p, queryDBForSingleUserInfo(pDB, i))
		}
	}
	return p
}

//Query the database for the info of a single user. Takes user ID of type int as input, returns struct of single user
func queryDBForSingleUserInfo(db *sql.DB, uid int) User {

	rows, err := db.Query("select * from user where user_id=?", uid)
	checkErr(err)

	var pid int
	//variables to store the rows.Scan below
	var firstname, lastname, mail, address, postnrandplace, phonenr, orgnr, countryID string
	//Next prepares the next result row for reading with the Scan method. It returns true on success,
	//or false if there is no next result row or an error happened while preparing it.
	//Err should be consulted to distinguish between the two cases.
	for rows.Next() {
		//Scan copies the columns in the current row into the values pointed at by dest.
		//The number of values in dest must be the same as the number of columns in Rows of database.
		rows.Scan(&pid, &firstname, &lastname, &mail, &address, &postnrandplace, &phonenr, &orgnr, &countryID)

	}
	m := User{}
	m.Number = pid
	m.FirstName = firstname
	m.LastName = lastname
	m.Mail = mail
	m.Address = address
	m.PostNrAndPlace = postnrandplace
	m.PhoneNr = phonenr
	m.OrgNr = orgnr
	m.CountryID = countryID

	defer rows.Close()
	return m
}

//input *sql.DB and returns the highest uid number, and line count of rows in DB
func queryDBForLastCustomerUID(db *sql.DB) (int, int) {
	rows, err := db.Query("select user_id from user")
	checkErr(err)
	defer rows.Close()
	//Prepare the slice to store numbers read from DB
	var num []int
	for rows.Next() {
		var pid int
		//The number of values below must be the same amount
		//as the number of rows in the DB
		err := rows.Scan(&pid) //puts data into the address of the variable
		checkErr(err)
		num = append(num, pid)
	}

	highestNr := 0
	countLines := 0
	//iterate the slice, and find the highest number, and number of lines.
	for i := range num {
		if highestNr < num[i] {
			highestNr = num[i]
			//log.Println("queryDBForLastCustomerUID : highestNr = ", highestNr)
			countLines++
		}
	}
	log.Println("queryDBForLastCustomerUID: highestNr = ", highestNr)
	log.Println("queryDBForLastCustomerUID: countLines = ", countLines)
	return highestNr, countLines
}

//Update user in Database, takes pointer to db and type User struct as input
func updateUserInDB(db *sql.DB, u User) {
	tx, err := db.Begin()
	checkErr(err)

	log.Println("The org nr. sendt to updateUserDB function = ", u.OrgNr)
	stmt, err := tx.Prepare("UPDATE user SET user_id=?,first_name=?,last_name=?,mail=?,address=?,post_nr_place=?,phone_nr=?,org_nr=?,country_id=? WHERE user_id=?")
	checkErr(err)
	defer stmt.Close()
	log.Println("updateUserInDB : Number in updateUserInDB function = ", u.Number)
	log.Println("************", u.Number, u.FirstName, u.LastName, u.Mail, u.Address, u.PostNrAndPlace, u.PhoneNr, u.OrgNr, u.CountryID, u.Number, "*************")
	//number is passed an extra time at the end of DB statement to fill the variable for the Query, which is done by number of user
	_, err = stmt.Exec(u.Number, u.FirstName, u.LastName, u.Mail, u.Address, u.PostNrAndPlace, u.PhoneNr, u.OrgNr, u.CountryID, u.Number)

	tx.Commit()
	checkErr(err)

}

//Adds user to Database. takes pointer to DB, and type User struct as input
func addUserToDB(db *sql.DB, u User) {
	//start db session
	tx, err := db.Begin()
	checkErr(err)

	//create statement to insert values to DB
	stmt, err := tx.Prepare("insert into user(user_id,first_name,last_name,mail,address,post_nr_place,phone_nr,org_nr,country_id) values(?,?,?,?,?,?,?,?,?)")
	checkErr(err)
	//At the end of function close the DB
	defer stmt.Close()

	//execute the statement on the DB
	_, err = stmt.Exec(u.Number, u.FirstName, u.LastName, u.Mail, u.Address, u.PostNrAndPlace, u.PhoneNr, u.OrgNr, u.CountryID)
	//commit to DB
	tx.Commit()
	checkErr(err)

}

//************************** BILL SECTION ***********************************

//Query the database all the bill lines for a specific bill nr. Takes bill_id of type int as input,
//returns a slice of struct type BillLines
func queryDBForBillLinesInfo(db *sql.DB, billID int) []BillLines {

	rows, err := db.Query("select * from bill_lines where bill_id=?", billID)
	checkErr(err)

	//variables to store the rows.Scan below

	var indx int
	/*var indx, BillID, LineID, ItemID, Quantity, DiscountPercentage, VatUsed int
	var Description string
	var PriceExVat float64*/

	//used to store the single row values read from DB
	mm := BillLines{}

	//used to store a slice of all the values from mm
	m := []BillLines{}

	//Next prepares the next result row for reading with the Scan method. It returns true on success,
	//or false if there is no next result row or an error happened while preparing it.
	//Err should be consulted to distinguish between the two cases.
	for rows.Next() {
		//Scan copies the columns in the current row into the values pointed at by dest.
		//The number of values in dest must be the same as the number of columns in Rows of database.
		rows.Scan(&indx, &mm.BillID, &mm.LineID, &mm.ItemID, &mm.Description, &mm.Quantity, &mm.DiscountPercentage, &mm.VatUsed, &mm.PriceExVat)

		//DO THIS WORK ?????????????? trying to append a struct into a slice of structs of the same type
		m = append(m, mm)
		log.Println("queryDBForBillLinesInfo: Content of mm : ", mm)
	}
	log.Println("queryDBForBillLinesInfo: Content of m : ", mm)
	defer rows.Close()
	return m
}

//------------
//Query the database all the bills for a specific user_id. Takes user_id of type int as input,
//returns a slice of struct type Bill
func queryDBForBillsForUser(db *sql.DB, userID int) []Bill {

	rows, err := db.Query("select * from bills where user_id=?", userID)
	if err != nil {
		log.Println("ERROR: queryDBForBillsForUser: Query error = ", err)
	}

	//used to store the single row values read from DB
	mm := Bill{}

	//used to store a slice of all the values from mm
	m := []Bill{}

	//Next prepares the next result row for reading with the Scan method. It returns true on success,
	//or false if there is no next result row or an error happened while preparing it.
	//Err should be consulted to distinguish between the two cases.
	for rows.Next() {
		//Scan copies the columns in the current row into the values pointed at by dest.
		//The number of values in dest must be the same as the number of columns in Rows of database.
		rows.Scan(&mm.BillID, &mm.UserID, &mm.CreatedDate, &mm.DueDate, &mm.Comment, &mm.TotalExVat, &mm.TotalIncVat, &mm.Paid)

		//DO THIS WORK ?????????????? trying to append a struct into a slice of structs of the same type
		m = append(m, mm)
		log.Println("queryDBForBillsForUser: Content of mm : ", mm)
	}
	log.Println("queryDBForBillsForUser: Content of m : ", mm)
	defer rows.Close()

	fmt.Println("INFO: queryDBForBillsForUser : The slice of bills for user looks like = ", m)
	return m
}

//------------

//input *sql.DB and returns the highest bill number, and line count of rows in DB
func queryDBForLastBillID(db *sql.DB) (int, int) {
	rows, err := db.Query("SELECT bill_id FROM bills")
	checkErr(err)
	defer rows.Close()

	//Prepare the slice to store numbers read from DB
	var num []int

	for rows.Next() {
		var readValue int
		//The number of values below must be the same amount
		//as the number of rows in the DB
		err := rows.Scan(&readValue) //reads data from db and puts it into the address of the variable
		checkErr(err)
		num = append(num, readValue)
	}

	highestNr := 0
	countLines := 0
	//iterate the slice, and find the highest number, and number of lines.
	for i := range num {
		if highestNr < num[i] {
			highestNr = num[i]
			//log.Println("queryDBForLastBillID : highestNr = ", highestNr)
			countLines++
		}
	}
	log.Println("queryDBForLastBillID: highestNr = ", highestNr)
	log.Println("queryDBForLastBillID: countLines = ", countLines)
	return highestNr, countLines
}

//Adds new bill to Database. takes pointer to DB, and type bill struct as input. Returns bill ID of type int
func addBillToDB(db *sql.DB, b Bill) int {
	//start db session
	tx, err := db.Begin()
	if err != nil {
		log.Println("ERROR: addBillToDB: ", err)
	}

	//create statement to insert values to DB
	stmt, err := tx.Prepare("insert into bills(bill_id,user_id,create_date,due_date,comment,total_ex_vat,total_inc_vat,paid) values(?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Println("ERROR: addBillToDB: statement problem = ", err)
	}
	//At the end of function close the DB
	defer stmt.Close()

	log.Println("addBillToDB: The struct b of type Bill contains = ", b)

	//execute the statement on the DB
	_, err = stmt.Exec(b.BillID, b.UserID, b.CreatedDate, b.DueDate, b.Comment, b.TotalExVat, b.TotalIncVat, b.Paid)
	//commit to DB
	err = tx.Commit()
	if err != nil {
		log.Println("ERROR: addBillToDB: stmt.Exec problem = ", err)
	}

	return b.BillID
}

//WORKING HERE !!!
//Adds new bill line to Database. takes pointer to DB, and type BillLines struct as input
//Create a function to keep track of the next available indx number in database
func addBillLineToDB(db *sql.DB, b BillLines) {
	//start db session
	tx, err := db.Begin()
	if err != nil {
		log.Println("ERROR: addBillLineToDB: Error starting db.begin = ", err)
	}

	//create statement to insert values to DB
	stmt, err := tx.Prepare("insert into bill_lines(indx,bill_id,line_id,item_id,description,quantity,discount_percentage,vat_used,price_ex_vat) values(?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Println("ERROR: addBillLineToDB: stmt error = ", err)
	}
	//At the end of function close the DB
	defer stmt.Close()

	//get last used index number in indx row,
	// and increment it by one to prepare for the next record
	indx, _ := queryDBForLastBillLineIndx(pDB)
	indx++

	//execute the statement on the DB
	_, err = stmt.Exec(indx, b.BillID, b.LineID, b.ItemID, b.Description, b.Quantity, b.DiscountPercentage, b.VatUsed, b.PriceExVat)
	if err != nil {
		log.Println("ERROR: addBillLineToDB: stmt.Execute error = ", err)
	}
	//commit to DB
	err = tx.Commit()
	if err != nil {
		log.Println("ERROR: addBillLineToDB: tx.Commit error = ", err)
	}

}

//query db for the last used Bill Line Index. Returns last used indx, and lineCount
func queryDBForLastBillLineIndx(db *sql.DB) (int, int) {
	rows, err := db.Query("SELECT indx FROM bill_lines")
	checkErr(err)
	defer rows.Close()

	//Prepare the slice to store numbers read from DB
	var num []int

	for rows.Next() {
		var readValue int
		//The number of values below must be the same amount
		//as the number of rows in the DB
		err := rows.Scan(&readValue) //reads data from db and puts it into the address of the variable
		checkErr(err)
		num = append(num, readValue)
	}

	highestNr := 0
	countLines := 0
	//iterate the slice, and find the highest number, and number of lines.
	for i := range num {
		if highestNr < num[i] {
			highestNr = num[i]
			//log.Println("queryDBForLastBillLineIndx : highestNr = ", highestNr)
			countLines++
		}
	}
	log.Println("queryDBForLastBillLineIndx: highestNr = ", highestNr)
	log.Println("queryDBForLastBillLineIndx: countLines = ", countLines)
	return highestNr, countLines
}

//query db for the last used Bill Line for specific bill.
//Input: *sql.DBReturns, and billID. Returns: last used bill line, and lineCount
func queryDBForLastBillLine(db *sql.DB, billID int) (int, int) {
	rows, err := db.Query("SELECT line_id FROM bill_lines WHERE bill_id=?", billID)
	checkErr(err)
	defer rows.Close()

	//Prepare the slice to store numbers read from DB
	var num []int

	for rows.Next() {
		var readValue int
		//The number of values below must be the same amount
		//as the number of rows in the DB
		err := rows.Scan(&readValue) //reads data from db and puts it into the address of the variable
		checkErr(err)
		num = append(num, readValue)
	}

	highestNr := 0
	countLines := 0
	//iterate the slice, and find the highest number, and number of lines.
	for i := range num {
		if highestNr < num[i] {
			highestNr = num[i]
			//log.Println("queryDBForLastBillLineIndx : highestNr = ", highestNr)
			countLines++
		}
	}
	log.Println("queryDBForLastBillLineIndx: highestNr = ", highestNr)
	log.Println("queryDBForLastBillLineIndx: countLines = ", countLines)
	return highestNr, countLines
}

//**************************  creates the database  ********************************
func createDB() *sql.DB {
	//1. Open connection

	db, err := sql.Open("sqlite3", "./fakt.db") //return types = *DB, error
	checkErr(err)
	//2. fail-fast if can't connect to DB
	checkErr(db.Ping())
	//3. create table
	_, err = db.Exec(`create table if not exists user (
						pid integer PRIMARY KEY, 
						firstname string not null,
						lastname string,
						mail string,
						address string,
						postnrandplace string,
						phonenr string,
						orgnr string)
					;`)
	checkErr(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bill_lines (
						indx int PRIMARY KEY,
						bill_id int,
						line_id int,
						item_id int,
						description string,
						quantity int,
						discount_percentage int,
						vat_used int,
						price_ex_vat real)
					;`)
	checkErr(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bills (
						bill_id int PRIMARY KEY,
						user_id int,
						create_date text,
						due_date text,
						comment string,
						total_ex_vat real,
						total_inc_vat real,
						paid integer)
					;`)

	return db
}

//if error !=nil, print error message to web page
func checkErr(err error, args ...string) {
	if err != nil {
		log.Printf("ERROR : %q: %s\n", err, args)
	}
}

//Delete a row in user DB, takes pointer to db, and index number uid which corresponds to column 1 in DB for input
func deleteUserInDB(db *sql.DB, number int) {
	tx, err := db.Begin()
	checkErr(err)
	log.Println("deleteUserInDB: The index number of the person to delete is = ", number)

	//Make the sql statement to execute
	stmt, err := tx.Prepare("DELETE FROM user WHERE user_id=?")
	checkErr(err)

	defer stmt.Close()
	//prepare the statement with a value for the "?"
	_, err = stmt.Exec(number)
	tx.Commit()
	checkErr(err)

}
