# Create Bill's using Go

Data is stored in local sqlite3 DB.
Access to the DB is done via the data package.

## TODO: WebSocket

* `send(stringValue)` will send the given string to Go websocketHandler.
* `deleteElement(elementID)` will delete the element with the given ID.

`<button onclick="send('modifyUserSelection')">show modify user selection</button>`

Dynamic content on/in a page shall be controlled by the button of field who changes a value,
or request something new drawn into the page.
A bill page need div's like this:

* DivMain
  * DivUser
  * DivBillInfo
  * DivBillLines
    * DivLine
    * DIVLine
            ...
With a layout like this we can delete and redraw single div elements to update content.
We need a html template for each type of div, wich is drawn, or redrawn based on the
action chosen by the Admin.

If we wanted to draw a complete bill page by choosing for example a BillID and not starting with a UserID
as above, we would need a complete template for that function outlined like above. That template would
be buildt up and calling all the other templates to draw the complete page.

## TODO: other

* Keep the number 0 in the deleted user row, increase the last user is deleted then a new used added will get that number
* Change the user pages to 1 page with add modify and delete
* Make the primary keys uid and bill ID random numbers, so you can sync the database between different devices without getting a conflict.