# todo_rest_service
simple rest application in golang with jwt authentication for apis and mysql as db  
To run the application goto folder with main.go and run cmd :  
```go run main.go todo.go todoHandlers.go tokenHandler.go dbHandler.go handleNull.go  ```  
The application opens at HOSTNAME:8080.(tested at localhost as HOSTNAME)  
## List of apis 
* **GET** ``HOSTNAME:8080/`` displays welcome message(for checking connectivity)  
* **GET** ``` HOSTNAME:8080/todo?id=&name=&description=&priority=&due=&completed=&completion_date=``` GET todo based on filters in query parameters, if no param in given all records returned  
* **POST** ``HOSTNAME:8080/todo`` CREATE new todo record
* **GET** ``HOSTNAME:8080/todo/{id}`` READ todo record based on {id}
* **PUT** ``HOSTNAME:8080/todo/{id}`` UPDATE todo values for {id}
* **DELETE** ``HOSTNAME:8080/todo/{id}`` DELETE todo based on {id}
* **GET** ``HOSTNAME:8080//getToken`` GET JWT token to be used for all other apis
## X-Security-Token is mandatory as header for all the api except /getToken