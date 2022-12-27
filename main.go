package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// go run main.go
func main() {
	//open database
	db, err := gorm.Open(mysql.Open("test:12345678@tcp(203.154.71.142:3306)/exam?charset=utf8mb4&parseTime=True&loc=Local"),
		&gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db

	//start echo
	e := echo.New()
	//path param
	e.GET("/users/:id", getUser, testMiddleware)
	//query param
	e.GET("/show", getShow)
	//json
	e.POST("/users", getUser2)
	//form file
	e.POST("/save", save)
	//quiz use gorm
	e.GET("/bookings", getBooking)

	e.Start(":1325")
}

// BookingHistory struct for store data from database
type BookingHistory struct {
	ID            int        `json:"id"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	MaximumPerson int        `json:"maximum_person"`
	SumGrade      string     `json:"sum_grade"`
}

// TableName reference struct BookingHistory to table bookings in database
func (BookingHistory) TableName() string {
	return "bookings"
}

// use gorm
func getBooking(c echo.Context) error {
	var bookingHistories []BookingHistory
	DB.Select([]string{
		"bookings.id",
		"first_name",
		"last_name",
		"start_date",
		"end_date",
		"maximum_person",
		`CASE WHEN sum_grade = 'A'THEN 'ดีมาก' 
				WHEN sum_grade = 'B' THEN 'ดี' 
				WHEN sum_grade = 'C' THEN 'พอใช้' 
				WHEN sum_grade = 'D' THEN 'ปรับปรุง' 
				ELSE 'ยังไม่มีเกรด' 
			END AS sum_grade`}).
		Joins("LEFT JOIN users ON users.id = bookings.users_id").
		Joins("LEFT JOIN rooms ON rooms.id = bookings.rooms_id").
		Find(&bookingHistories)
	return c.JSON(200, bookingHistories)
}

// request by path param
func getUser(c echo.Context) error {
	id := c.Param("id")
	return c.String(200, id)
}

// request by query param
type TeamMember struct {
	Team   string `query:"team"`
	Member int    `query:"member"`
}

// request by query param
func getShow(c echo.Context) error {
	//t := c.QueryParam("team")
	//m := c.QueryParam("member")
	var tm TeamMember

	err := c.Bind(&tm)
	if err != nil {
		return err
	}
	return c.JSON(200, tm)
}

// request by json
type Users struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// request by json
func getUser2(c echo.Context) error {
	var u Users
	//u := Users{}
	//u := new(Users)
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	return c.JSON(200, u)
}

// request by form-data
func save(c echo.Context) error {
	//n := c.FormValue("name")
	image, err := c.FormFile("file")
	if err != nil {
		return c.JSON(400, err.Error())
	}

	src, err := image.Open()
	defer src.Close()

	path := fmt.Sprintf("picture/12312342vcxz13423.png")
	dst, err := os.Create(path)
	defer dst.Close()

	io.Copy(dst, src)

	return c.JSON(200, image)
}

// middleware
func testMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Request().Header.Get("key")
		if key != "test" {
			return c.JSON(401, "no key")
		}
		return next(c)
	}
}
