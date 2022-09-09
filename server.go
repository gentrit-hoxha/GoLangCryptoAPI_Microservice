package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	type Coin struct {
		Symbol             string `json:"symbol"`
		PriceChange        string `json:"priceChange"`
		PriceChangePercent string `json:"priceChangePercent"`
	}
	type Coins struct {
		Coins []Coin `json:"coin"`
	}

	db, err := sql.Open("mysql", "root:Gentrit-2022@tcp(127.0.0.1:3306)/CryptoCoin_API")

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("db is connected")
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}

	//BY this endpoint we will pass a parameter and delete that record in the database
	e.DELETE("/coins/:coin_symbol", func(c echo.Context) error {
		requested_id := c.Param("coin_symbol")
		sql := "Delete from Coin where symbol = ?"
		stmt, err := db.Prepare(sql)
		if err != nil {
			fmt.Println(err)
		}
		result, err2 := stmt.Exec(requested_id)
		if err2 != nil {
			panic(err2)
		}
		fmt.Println(result.RowsAffected())
		return c.JSON(http.StatusOK, "Deleted")
	})

	// BY this endpoint we will get all the Crypto Coin values in the database
	e.GET("/coins", func(ctx echo.Context) error {
		sql := "Select * from Coin"

		if err != nil {
			fmt.Println(err)
		}

		rows, err2 := db.Query(sql)
		fmt.Println(rows)
		fmt.Println(err2)

		if err2 != nil {
			panic(err2)
		}

		defer rows.Close()
		result := Coins{}

		for rows.Next() {
			coin := Coin{}
			err3 := rows.Scan(&coin.Symbol, &coin.PriceChange, &coin.PriceChangePercent)
			// Exit if we get an error
			if err3 != nil {
				fmt.Print(err3)
			}
			result.Coins = append(result.Coins, coin)
		}

		return ctx.JSON(http.StatusOK, result)

	})

	// This endpoint will create a Crypto Coin record in the database passing by just the symbol
	// and fetching data from an external API
	e.POST("/coins/:coin_symbol", func(ctx echo.Context) error {
		// Build the request

		requestedSymbol := ctx.Param("coin_symbol")

		req, err := http.NewRequest("GET", "https://api2.binance.com/api/v3/ticker/24hr?symbol="+requestedSymbol+"", nil)
		if err != nil {
			fmt.Println("Error is req: ", err)
		}

		// create a Client
		client := &http.Client{}

		// Do sends an HTTP request and
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error in send req: ", err)
		}

		// Fill the data with the data from the JSON
		var data Coin

		// Use json.Decode for reading streams of JSON data
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			fmt.Println(err)
		}

		if data.Symbol == "" {
			fmt.Println("No data for this coin you are looking for")
			return ctx.JSON(http.StatusBadRequest, "No data for the coin you are looking for")
		} else {
			fmt.Println(data.Symbol)
			fmt.Println(data.PriceChange)
			fmt.Println(data.PriceChangePercent)

			sql := "INSERT INTO Coin(symbol, priceChange, priceChangePercent) VALUES( ?, ?, ?)"
			stmt, err := db.Prepare(sql)

			if err != nil {
				fmt.Print(err.Error())
			}
			defer stmt.Close()
			result, err2 := stmt.Exec(data.Symbol, data.PriceChange, data.PriceChangePercent)

			// Exit if we get an error
			if err2 != nil {
				panic(err2)
			}
			fmt.Println(result.LastInsertId())
			return ctx.JSON(http.StatusCreated, data.Symbol)
		}

	})

	//This will create a new COIN in the database from all the parameters that come from us
	e.POST("/coin", func(c echo.Context) error {
		emp := new(Coin)
		if err := c.Bind(emp); err != nil {
			return err
		}
		//
		sql := "INSERT INTO Coin(symbol, priceChange, priceChangePercent) VALUES( ?, ?, ?)"
		stmt, err := db.Prepare(sql)

		if err != nil {
			fmt.Print(err.Error())
		}
		defer stmt.Close()
		result, err2 := stmt.Exec(emp.Symbol, emp.PriceChange, emp.PriceChangePercent)

		// Exit if we get an error
		if err2 != nil {
			panic(err2)
		}
		fmt.Println(result.LastInsertId())

		return c.JSON(http.StatusCreated, emp.Symbol)
	})

	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
