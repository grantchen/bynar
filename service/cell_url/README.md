For run this project you need to 
1.create database Name test
2.Dump db.sql file 
3.Run the project using go run main.go
4.Once project runs then use this url http://localhost:9003/static/test.html for front-end 
5.Do change as per your mysql credential
6.I added the env file you can do a changes as per your requirements


## Requirement
### 1. The developer should make the code structural, compactable and understandable. Exactly how the developer should do on real project environment.

I've reorganized the code into multi packages. The code structure contains layers:

- handler: receive a http request, parse param, call to service layer, check error code and return to client.
- service: receive a call from handler layer, get data from repository layer, do logics, and return the data.
- repository: on behalf of communicating between application and database.
- config: read config from env file, contain functions to return configs.
- main.go: Start point of application Init layers.

### 2. Remove ORM and use plain SQL/MySQL libraries, the same for Gorilla mux. Use standard libraries. No 3d party libraries but only native libraries for important services.

I've removed ORM and use database/sql library. Code is seperated in *repository/datagrid_repository.go*. In the old code, We had an sql injection. (concat query string). I resolved it by using prepareStatement.

### 3. Remove the EOF parameters conditions. In the example is just to show how it works so it will not be necessary.
### 4. In the tree gird API request use JSON instead of xml by adding in html BDO tags Cell_Format="JSON" and change the generated format in back-end to JSON. See the documentation.
I've followed the documentation, adding tag JSON, and rewrite initial config to test.json instead of test.xml.

### 5. In this example the html tags are generated in back-end. We need to remove the generation of html tags in back end and as response we need only JSON results. The developer can provide a solution that the function only responds the JSON results and the drop down to be generated in other method (front-end) and not in back-end.

Follow the documentation of tree grid, I've used function:

```
AjaxCall(url, param, function (code, data) {
	data = generateHTMLData(data);
	callback(code, data);
});
```
to pre-process response data from server, adding addition html tag and pass to the next step. 

In backend side, I return only JSON result (not html tag is generated).