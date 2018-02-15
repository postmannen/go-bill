## History

2. Tested person struct, and string to int conversion for person ID  
    Iterate the slice of struct  
3. Added invoice nr. to personStruct  
4. Added some menus with print and add person as options  
	Added getPersonNextNr function to look up the next available person nr  
5. Added auto next number for person  
	Added more variables describing person  
	Added /sp for show person info  
	Added /ap for add person  
	input only get added when "add" button is pushed  
6.	Add sqlite with add data functions  
	Removed some of the old not needed code  
7.	Added query functions towards sqlite db  
	Removed the not needed invoice variable from struct, and DB  
	Added dropdown list for /md (modify page)  
8.	Added templates.html  
	Added templates to the handlers  
	Nested the templates within the templates file, so only 1 call is needed for each web page, and not seperate calls for header, menu..etc...  
9.	Add modify person http, functions and database updates  
	 TODO: Use a temp struct instead of all the single variables in the modify http handler function Look into replacing the if to check update of fields with a switch/case  
10. Cleanup  
		Rename fakt.go to main.go  
		Update comments  
	Made the variable indexNR global to store the selected user in the modify form and function  
11. Cleanup  
		renamed where the name person(s) where used to user(s)  
		removed unused code  
12. Tested with html and CSS, but dont really understand how to align boxes, text etc.  
13.	Added orgnr to user table  
		ERROR : Does not update org nr in modify section  
			The problem is only with the modify function, add works ok  
	Changed the top menu to use links insted of input box's, and dropdown with css  
	showAllUsers : replaced the input boxes with a table  
14. Changing the css styling of the pages  
	Added delete user function  
	Fixed an Error for counting highest nr, and number of rows which were not visible until the delete function were implemented  
	Rewrote the function to get next index number to get last index number, and return highest user uid, and count of total uid's  
15. Rewrote the DB table names to use all small letters, and underscore to seperate words  
	 Changed all the code to reflect changes  
16. Wrote the first db template to use in "template-database-creation.sql"  
    Rewrote the addUser* functions to use type User (struct) instead of single variables of type int and string  
    Rewrote the modifyUser* functions to use type User (struct) instead of single variables of type int and string  
    Split the main package into main.go, web.go, and db.go.  
    Created README.md  
	Created the initiall bill page
	Added the choose user html part
17. Added bills table
    Added new bill functionality to create bill page
	Added new line to bill listing
18. Rewrite so all templates are drawn on bill create startup, and not after some action is taken
    Testing with the bill pages
	Renaming of a whole lot of variables
19. Making the 'data' variable used for webPages parsing of data global. Not sure if that is a smart
    idea, but time will show
20. Moved DB into package, and renamed the functions in db with better naming
21. Moved the content of db package into data package
	Fixed bill lines view
	Removed the global data variable, and replaced with using recievers on handler functions
	Split handlers into several files
22. Looking into minimalizing the code with switch statements etc.
    Splitting up the big functions into smaller functions
    Fixed the sum fields in bill view, also writes the data to db
    Split the functions used in handlers into separate to easyer get an overview
	Rewrote the css.
23. Added bill print page
24. Added Bank Account to DB, and changed all code to reflect the new table
	

## TODO
* Keep the number 0 in the deleted user row, incase the last user is deleted then a new used added will get that number
* Change the user pages to 1 page with add modify and delete
* Learn to pass values to web HandleFunc's, to make the exported user struct not exported
* Check if templates can be reused

## Ideas
* Make the primary keys uid and bill ID random numbers, so you can sync the database between different devices without getting a conflict.
* Sorting can be done on a dummy index value that don't have to be unique

