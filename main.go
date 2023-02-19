package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title  string `gorm:"not null"`
	IsDone bool   `gorm:"not null"`
}

var Db *gorm.DB

func init() {
	dsn := "root:root@tcp(127.0.0.1:3306)/todo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Db = db
}

func findById(c echo.Context) error {
	// todo1 := Todo{}
	// result1 := Db.First(&todo1)
	// if errors.Is(result1.Error, gorm.ErrRecordNotFound) {
	// 	log.Fatal(result1.Error)
	// }

	// todo2 := Todo{}
	// result2 := Db.Take(&todo2)
	// if errors.Is(result2.Error, gorm.ErrRecordNotFound) {
	// 	log.Fatal(result2.Error)
	// }

	// todo3 := Todo{}
	// result3 := Db.Last(&todo3)
	// if errors.Is(result3.Error, gorm.ErrRecordNotFound) {
	// 	log.Fatal(result3.Error)
	// }

	todo4 := Todo{}
	// result4 := Db.First(&todo4, 1)
	result4 := Db.First(&todo4, "id = ?", 1)
	if errors.Is(result4.Error, gorm.ErrRecordNotFound) {
		log.Fatal(result4.Error)
	}

	return c.JSON(http.StatusOK, todo4)
}

func findAll(c echo.Context) error {
	todos := []Todo{}
	result := Db.Find(&todos)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	return c.JSON(http.StatusOK, todos)
}

func create(c echo.Context) error {
	todo := Todo{
		Title:  "タイトル",
		IsDone: false,
	}
	result := Db.Create(&todo)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	return c.JSON(http.StatusCreated, todo)
}

func creates(c echo.Context) error {
	todos := []Todo{
		{
			Title:  "タイトル1",
			IsDone: false,
		},
		{
			Title:  "タイトル2",
			IsDone: true,
		},
		{
			Title:  "タイトル3",
			IsDone: false,
		},
	}
	result := Db.Create(&todos)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	fmt.Println("count:", result.RowsAffected)

	return c.JSON(http.StatusCreated, todos)
}

func save(c echo.Context) error {
	todo1 := Todo{}
	todo1.Title = "タイトル10"
	result1 := Db.Save(&todo1)
	if result1.Error != nil {
		log.Fatal(result1.Error)
	}

	todo2 := Todo{}
	Db.First(&todo2)

	todo2.Title = "タイトル2"
	result2 := Db.Save(&todo2)
	if result2.Error != nil {
		log.Fatal(result2.Error)
	}

	return c.JSON(http.StatusOK, todo1)
}

func update(c echo.Context) error {
	result := Db.Model(&Todo{}).Where("id = 1").Update("title", "更新されたタイトル")
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	todo := Todo{}
	Db.Where("id = 1").Take(&todo)
	return c.JSON(http.StatusOK, todo)
}

func updates(c echo.Context) error {
	result := Db.Model(&Todo{}).Where("title = 'タイトル1'").Updates(Todo{Title: "一斉更新", IsDone: true})
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	todo := Todo{}
	Db.Where("title = '一斉更新'").Take(&todo)
	return c.JSON(http.StatusOK, todo)
}

func delete(c echo.Context) error {
	Db.Where("id = 1").Delete(&Todo{})

	return c.String(http.StatusOK, "削除できました｡")
}

func main() {
	Db.AutoMigrate(&Todo{})

	e := echo.New()
	e.GET("/api/v1/todo", findById)
	e.GET("/api/v1/todos", findAll)
	e.POST("/api/v1/todo", create)
	e.POST("/api/v1/todos", creates)
	e.PUT("/api/v1/todo/save", save)
	e.PUT("/api/v1/todo/update", update)
	e.PUT("/api/v1/todo/updates", updates)
	e.DELETE("/api/v1/todo/delete", delete)

	e.Logger.Fatal(e.Start(":1323"))
}
