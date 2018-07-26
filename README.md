# Create Bill's using Go

Data is stored in local sqlite3 DB.
Access to the DB is done via the data package.

## TODO: Using the WebSocket from HTML pages

The `webs.ocket.html` file contains the JavaScript to be executed in browsers.
Functions to be called either from other HTML elements or via the Websocket :

* `send(stringValue)` will send the given string to Go websocketHandler.
* `deleteElement(elementID)` will delete the element with the given ID.

`<button onclick="send('modifyUserSelection')">show modify user selection</button>`

## TODO: other

* Keep the number 0 in the deleted user row, increase the last user is deleted then a new used added will get that number
* Change the user pages to 1 page with add modify and delete
* Make the primary keys uid and bill ID random numbers, so you can sync the database between different devices without getting a conflict.