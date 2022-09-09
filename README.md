# First Execise

By this exercise is intended to help deepen our understanding of GoLang, REST APIs, microservices, and Echo framework.

### CRYPTO_API Microservice
Create a local database using MySQL. 
Create a  table Coin in the database with the following attributes: symbol, priceChange, priceChangePercentage
Create a REST API using Echo web framework. 
Your API should communicate with this public API where you are supposed to get details of different crypto coins:
https://api2.binance.com/api/v3/ticker/24hr 

#### This API should have the following endpoints
* GET Request **/api/coins:**  Get a list of all the coins that are saved in the database 
* POST Request  **/api/coins/:coin_symbol:** This should get data from the public API for this coin and save it in the database. You should return an error message in case this coin is not returned as a valid result from the public URL.
* DELETE Request **/api/coins/:coin_symbol:** This should delete all rows in the database that have this coin