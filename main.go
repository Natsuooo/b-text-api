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
  
  router.POST("/signup_with_img", signupWithImg)

  router.Run(":" + port)
}

type User struct{
  Id int `json:"id"`
  Uid string `json:"uid"`
  Username string `json:"username"`
  University string `json:"university"`
  Profile_image string `json:"profile_image"`
  Sns_image string `json:"sns_image"`
  Is_signup_detail bool `json:"is_signup_detail"`
  Unread_messages []string `json:"unread_messages"`
  New_message string `json:"new_message"`
}

func dbFunc(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    test := Test{}
    db.QueryRow("SELECT content FROM test").Scan(&test.Content)
    c.JSON(http.StatusOK, test)
  }
}

func test(c *gin.Context){
  c.String(http.StatusOK, "test")
}

func dbFunc(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    stmt, err := db.Prepare("INSERT INTO users(uid, username, university, profile_image, sns_image, is_signup_detail) VALUES($1, $2, $3, $4, $5, $6) RETURNING uid")
    checkErr(err)
    uid := c.PostForm("uid")
    username := c.PostForm("username")
    university := c.PostForm("university")
    sns_image := ""
    file, _ := c.FormFile("profile_image")
    profile_image := uid+file.Filename
    c.SaveUploadedFile(file, "profile_images/"+profile_image)
    is_signup_detail := true
    stmt.Exec(uid, username, university, profile_image, sns_image, is_signup_detail)
    db.Close()
  }
}

func checkErr(err error) {
  if err != nil {
    panic(err)
  }
}

