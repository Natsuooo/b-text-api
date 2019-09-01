package test

import(
  "github.com/gin-gonic/gin"
  _ "github.com/heroku/x/hmetrics/onload"
  "database/sql"
  _ "github.com/lib/pq"
  "net/http"
  "models"
)

func dbFunc(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    test := models.Test{}
    db.QueryRow("SELECT content FROM test").Scan(&test.Content)
    c.JSON(http.StatusOK, test)
  }
}