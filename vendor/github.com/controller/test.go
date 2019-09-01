package controller

import(
  "github.com/gin-gonic/gin"
  _ "github.com/heroku/x/hmetrics/onload"
  "database/sql"
  _ "github.com/lib/pq"
)

func dbFunc(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    test := Test{}
    db.QueryRow("SELECT content FROM test").Scan(&test.Content)
    c.JSON(http.StatusOK, test)
  }
}