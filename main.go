package main

import (
  "log"
  "net/http"
  "os"

  "github.com/gin-gonic/gin"
  _ "github.com/heroku/x/hmetrics/onload"
  "database/sql"
  _ "github.com/lib/pq"
  "github.com/gin-contrib/cors"
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

//  router := gin.New()

  router := gin.Default()
  router.Use(cors.New(cors.Config{
//        AllowOrigins: []string{"http://localhost:8889"},
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
    AllowHeaders: []string{"*"},
  }))
  router.GET("/test", test)
  
  router.GET("/db", dbFunc(db))
  
  

  router.Run(":" + port)
}

func test(c *gin.Context){
  c.String(http.StatusOK, "test")
}

func dbFunc(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    test := Test{}
    db.QueryRow("SELECT content FROM test").Scan(&test.Content)
    c.JSON(http.StatusOK, test)
  }
}
