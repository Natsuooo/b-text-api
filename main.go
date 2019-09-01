package main

import (
  "log"
  "net/http"
  "os"

  "github.com/gin-gonic/gin"
  _ "github.com/heroku/x/hmetrics/onload"
  "database/sql"
  _ "github.com/lib/pq"
  "controller"
)

type Test struct{
  Id int `json:"id"`
  Content string `json:"content"`
}

func main() {
  port := os.Getenv("PORT")

  if port == "" {
      log.Fatal("$PORT must be set")
  }
  
  db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
  if err != nil {
    log.Fatalf("Error opening database: %q", err)
  }

  router := gin.New()
  router.Use(gin.Logger())
  router.LoadHTMLGlob("templates/*.tmpl.html")
  router.Static("/static", "static")

  router.GET("/", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.tmpl.html", nil)
  })

//  router.GET("/test", test)
  
//  router.GET("/db", dbFunc(db))
  
//  controller := test.Controller{}
  
  router.GET("/db_controller", dbFunc(db))

  router.Run(":" + port)
}

//func test(c *gin.Context){
//  c.String(http.StatusOK, "test")
//}

//func dbFunc(db *sql.DB) gin.HandlerFunc{
//  return func (c *gin.Context){
//    test := Test{}
//    db.QueryRow("SELECT content FROM test").Scan(&test.Content)
//    c.JSON(http.StatusOK, test)
//  }
//}
