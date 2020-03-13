package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type User struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/dummy_db")
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

	// http://localhost:1323
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, this was displayed using Echo!")
	})

	// http://localhost:1323/users
	e.POST("/users", func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, u)
	})

	// http://localhost:1323/users/Joe
	e.GET("/users/:id", getUser)

	// http://localhost:1323/show?team=x-men&member=wolverine
	e.GET("/show", show)

	e.POST("/save", func(c echo.Context) error {
		user := new(User)
		if err := c.Bind(user); err != nil {
			return err
		}
		//
		sql := "INSERT INTO user(name, email) VALUES( ?, ?)"
		stmt, err := db.Prepare(sql)

		if err != nil {
			fmt.Print(err.Error())
		}
		defer stmt.Close()
		result, err2 := stmt.Exec(user.Name, user.Email)

		// Exit if we get an error
		if err2 != nil {
			panic(err2)
		}
		fmt.Println(result.LastInsertId())

		return c.JSON(http.StatusCreated, user.Name)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, "Hello, "+id)
}

func show(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:"+team+", member:"+member)
}
