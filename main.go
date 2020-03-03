package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type App struct {
	db *gorm.DB
}

type User struct {
	Id          int    `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"id"`
	Firstname   string `json:"firstName" binding:"required" sql:"null"`
	Lastname    string `json:"lastName" binding:"required" sql:"null"`
	Email       string `json:"email" binding:"required" sql:"null"`
	Phone       string `json:"phone" binding:"required" sql:"null"`
	Address     string `json:"address" binding:"required" sql:"null"`
	Description string `json:"description" binding:"required" sql:"null"`
}

func (i *User) TableName() string {
	return "users"
}

func main() {

	router, db := initializeAPI()
	defer db.Close()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func initializeAPI() (*gin.Engine, *gorm.DB) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=template1 sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	var app = App{db}

	r := gin.Default()
	app.initializeRoutes(r)

	r.Run()

	return r, db
}

func (app *App) initializeRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
	r.GET("/users", app.returnGoods)
	r.NoRoute(func(c *gin.Context) {
		render(c, gin.H{"payload": "not found"})
	})

}

func (app *App) returnGoods(c *gin.Context) {
	users := []*User{}
	app.db.Find(&users)
	render(c, gin.H{"payload": users})
}

func render(c *gin.Context, data gin.H) {

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

	c.JSON(http.StatusOK, data["payload"])
}
