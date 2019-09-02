package main

import (
  "log"
  "net/http"
  "os"
  "github.com/gin-gonic/gin"
  _ "github.com/heroku/x/hmetrics/onload"
  "database/sql"
  "github.com/lib/pq"
  "github.com/gin-contrib/cors"
  "time"
)

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
  
  router.POST("/signup_with_img", signupWithImg(db))
  router.POST("/signup", signup(db))
  router.GET("/user", user(db))
  router.GET("/get_user", getUser(db))
  router.POST("/sell", sell(db))
  router.GET("/mybooks", mybooks(db))
  router.GET("/book_images/:original_image", bookImages)
  router.POST("/stop_selling", stopSelling(db))
  router.POST("/restart_selling", restartSelling(db))
  router.GET("/books", books(db))
  router.GET("/books/detail", bookDetail(db))
  router.POST("/likes/register", registerLike(db))
  router.POST("/likes/delete", deleteLike(db))
  router.GET("/likes", likes(db))
  router.GET("/likes/books", likedBooks(db))
  router.GET("/messages", messages(db))
  router.POST("/messages/send", send(db))
  router.GET("/messages/mybooks", messagesMyBooks(db))
  router.GET("/messages/users", messagesUsers(db))
  router.POST("/messages/read", readMessages(db))
  router.POST("/users/update_with_img", updateUsersWithImg(db))
  router.POST("/users/update", updateUsers(db))
  router.GET("/users/:profile_image", userImage)
  router.POST("/rates/create", createRate(db))
  router.POST("/rates/update", updateRate(db))
  router.GET("/rates", getRate(db))
  router.GET("/rates/my", getMyRates(db))
  router.GET("/messages/buy", getBuyMessages(db))
  router.GET("/messages/unread", getUnreadMessages(db))

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

type Book struct{
  Id int `json:"id"`
  User_id int `json:"user_id"`
  Google_image string `json:"google_image"`
  Original_image string `json:"original_image"`
  Title string `json:"title"`
  State string `json:"state"`
  Price int `json:"price"`
  Note string `json:"note"`
  University string `json:"university"`
  Liked int `json:"liked"`
  Is_public bool `json:"is_public"`
  Updated_at time.Time `json:"updated_at"`
  Created_at time.Time `json:"created_at"`
  Messages_id []string `json:"messages_id"`
  Unread_messages []string `json:"unread_messages"`
}

type Message struct{
  Id int `json:"id"`
  Book_id int `json:"book_id"`
  From_user_id int `json:"from_user_id"`
  To_user_id int `json:"to_user_id"`
  Content string `json:"content"`
  Is_read bool `json:"is_read"`
  Created_at time.Time `json:"created_at"`
  Count int `json:"count"`
}

type Like struct{
  Id int `json:"id"`
  User_id int `json:"user_id"`
  Book_id int `json:"book_id"`
}

type Rate struct{
  Id int `json:"id"`
  Rating int `json:"rating"`
  From_user_id int `json:"from_user_id"`
  To_user_id int `json:"to_user_id"`
}

func getUnreadMessages(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    to_user_id := c.Query("to_user_id")
    message:=Message{}
    db.QueryRow("SELECT COUNT(id) FROM messages WHERE to_user_id=$1 AND is_read=false;", to_user_id).Scan(&message.Count)
    db.Close()
    c.JSON(http.StatusOK, message)
  }
}

func getBuyMessages(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    from_user_id := c.Query("from_user_id")
    rows, _ := db.Query("SELECT DISTINCT books.id, books.user_id, google_image, original_image, title, state, price, note, liked, is_public, books.created_at, ARRAY(SELECT messages.id FROM messages WHERE messages.book_id=books.id AND messages.to_user_id=$1 AND messages.is_read=false) AS messages_id FROM books INNER JOIN messages ON (books.id = messages.book_id) WHERE messages.from_user_id=$1;", from_user_id)
    var books []Book
    for rows.Next(){
      b:=Book{}
      rows.Scan(&b.Id, &b.User_id, &b.Google_image, &b.Original_image, &b.Title, &b.State, &b.Price, &b.Note, &b.Liked, &b.Is_public, &b.Created_at, pq.Array(&b.Unread_messages))
      books = append(books, b)
    }
    db.Close()
    c.JSON(http.StatusOK, books)
  }
}

func getMyRates(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    to_user_id := c.Query("to_user_id")
    rows, _ := db.Query("SELECT id, rating, from_user_id, to_user_id FROM rates WHERE to_user_id=$1", to_user_id)
    var rates []Rate
    for rows.Next(){
      r:=Rate{}
      rows.Scan(&r.Id, &r.Rating, &r.From_user_id, &r.To_user_id)
      rates = append(rates, r)
    }
    db.Close()
    c.JSON(http.StatusOK, rates)
  }
}

func getRate(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    from_user_id := c.Query("from_user_id")
    to_user_id := c.Query("to_user_id")
    rate:=Rate{}
    db.QueryRow("SELECT id, rating, from_user_id, to_user_id FROM rates WHERE from_user_id=$1 AND to_user_id=$2 LIMIT 1", from_user_id, to_user_id).Scan(&rate.Id, &rate.Rating, &rate.From_user_id, &rate.To_user_id)
    db.Close()
    c.JSON(http.StatusOK, rate)
  }
}

func updateRate(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    from_user_id := c.PostForm("from_user_id")
    to_user_id := c.PostForm("to_user_id")
    rating := c.PostForm("rating")
    stmt, err := db.Prepare("UPDATE rates SET rating=$1 WHERE from_user_id=$2 AND to_user_id=$3")
    checkErr(err)
    stmt.Exec(rating, from_user_id, to_user_id)
    db.Close()
  }
}

func createRate(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    from_user_id := c.PostForm("from_user_id")
    to_user_id := c.PostForm("to_user_id")
    rating := c.PostForm("rating")
    stmt, err := db.Prepare("INSERT INTO rates(rating, from_user_id, to_user_id) VALUES($1, $2, $3)")
    checkErr(err)
    stmt.Exec(rating, from_user_id, to_user_id)
    db.Close()
  }
}

func userImage(c *gin.Context){
  profile_image := c.Param("profile_image")
  c.File("./profile_images/"+profile_image)
}

func updateUsers(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    user_id := c.PostForm("user_id")
    university := c.PostForm("university")
    username := c.PostForm("username")
    stmt, err := db.Prepare("UPDATE users SET username=$1, university=$2 WHERE id=$3")
    checkErr(err)
    stmt.Exec(username, university, user_id)
    db.Close()
  }
}

func updateUsersWithImg(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    user_id := c.PostForm("user_id")
    uid := c.PostForm("uid")
    university := c.PostForm("university")
    username := c.PostForm("username")
    file, _ := c.FormFile("profile_image")
    profile_image := uid+file.Filename
    c.SaveUploadedFile(file, "profile_images/"+profile_image)
    stmt, err := db.Prepare("UPDATE users SET username=$1, university=$2, profile_image=$3 , sns_image='' WHERE id=$4")
    checkErr(err)
    stmt.Exec(username, university, profile_image, user_id)
    db.Close()
  }
}

func readMessages(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    from_user_id := c.PostForm("from_user_id")
    to_user_id := c.PostForm("to_user_id")
    book_id := c.PostForm("book_id")
    stmt, err := db.Prepare("UPDATE messages SET is_read=true WHERE from_user_id=$1 AND to_user_id=$2 AND book_id=$3")
    checkErr(err)
    stmt.Exec(from_user_id, to_user_id, book_id)
    db.Close()
  }
}

func messagesUsers(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    book_id := c.Query("book_id")
    user_id := c.Query("user_id")
    rows, _ := db.Query("SELECT DISTINCT a.id, a.uid, a.username, a.university, a.profile_image, a.sns_image, ARRAY(SELECT b.id FROM messages b WHERE a.id=b.from_user_id AND b.is_read=false AND to_user_id=$1) AS unread_messages, (SELECT b.content from messages b WHERE a.id=b.from_user_id AND b.to_user_id=$2 ORDER BY b.created_at DESC LIMIT 1) AS new_message FROM users a INNER JOIN messages b ON b.from_user_id=a.id WHERE b.book_id=$3 AND b.to_user_id=$4", user_id, user_id, book_id, user_id)
    var users []User
    for rows.Next(){
      u:=User{}
      rows.Scan(&u.Id, &u.Uid, &u.Username, &u.University, &u.Profile_image, &u.Sns_image, pq.Array(&u.Unread_messages), &u.New_message)
      users = append(users, u)
    }
    db.Close()
    c.JSON(http.StatusOK, users)
  }
}

func messagesMyBooks(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    user_id := c.Query("user_id")
    rows, _ := db.Query("SELECT a.id, a.user_id, a.google_image, a.original_image, a.title, a.state, a.price, a.note, a.is_public, a.created_at, ARRAY(SELECT b.id FROM messages b WHERE b.book_id=a.id AND b.to_user_id=$1 AND b.is_read=false) AS messages_id FROM books a WHERE user_id=$2 ORDER BY a.created_at DESC", user_id, user_id)
    var mybooks []Book
    for rows.Next(){
      b:=Book{}
      rows.Scan(&b.Id, &b.User_id, &b.Google_image, &b.Original_image, &b.Title, &b.State, &b.Price, &b.Note, &b.Is_public, &b.Created_at, pq.Array(&b.Messages_id))
      mybooks = append(mybooks, b)
    }
    db.Close()
    c.JSON(http.StatusOK, mybooks)
  }
}

func send(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    stmt, err := db.Prepare("INSERT INTO messages(book_id, from_user_id, to_user_id, content, created_at) VALUES($1, $2, $3, $4, $5)")
    checkErr(err)
    book_id := c.PostForm("book_id")
    from_user_id := c.PostForm("from_user_id")
    to_user_id := c.PostForm("to_user_id")
    content := c.PostForm("content")
    created_at := c.PostForm("created_at")
    stmt.Exec(book_id, from_user_id, to_user_id, content, created_at)
    db.Close()
  }
}

func messages(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    user_id := c.Query("user_id")
    book_id := c.Query("book_id")
    rows, _ := db.Query("SELECT id, book_id, from_user_id, to_user_id, content, is_read, created_at FROM messages WHERE book_id=$1 AND (from_user_id=$2 OR to_user_id=$2) ORDER BY created_at ASC",book_id, user_id)
    var messages []Message
    for rows.Next(){
      m:=Message{}
      rows.Scan(&m.Id, &m.Book_id, &m.From_user_id, &m.To_user_id, &m.Content, &m.Is_read, &m.Created_at)
      messages = append(messages, m)
    }
    db.Close()
    c.JSON(http.StatusOK, messages)
  }
}

func likedBooks(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    user_id := c.Query("user_id")
    rows, _ := db.Query("SELECT books.id, books.user_id, google_image, original_image, title, state, price, note, liked, is_public, updated_at FROM books INNER JOIN likes ON (books.id = likes.book_id) WHERE likes.user_id=$1 ORDER BY likes.created_at DESC", user_id)
    var books []Book
    for rows.Next(){
      b:=Book{}
      rows.Scan(&b.Id, &b.User_id, &b.Google_image, &b.Original_image, &b.Title, &b.State, &b.Price, &b.Note, &b.Liked, &b.Is_public, &b.Updated_at)
      books = append(books, b)
    }
    db.Close()
    c.JSON(http.StatusOK, books)
  }
}

func likes(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    user_id := c.Query("user_id")
    rows, _ := db.Query("SELECT id, user_id, book_id FROM likes WHERE user_id=$1", user_id)
    var likes []Like
    for rows.Next(){
      like:=Like{}
      rows.Scan(&like.Id, &like.User_id, &like.Book_id)
      likes = append(likes, like)
    }
    db.Close()
    c.JSON(http.StatusOK, likes)
  }
}

func isLiked(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    user_id := c.Query("user_id")
    book_id := c.Query("book_id")
    like:=Like{}
    stmt, _ := db.Prepare("SELECT id, user_id, book_id FROM likes WHERE user_id=$1 AND book_id=$2 LIMIT 1")
    stmt.QueryRow(user_id, book_id).Scan(&like.Id, &like.User_id, &like.Book_id)
    db.Close()
    c.JSON(http.StatusOK, like)
  }
}

func deleteLike(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    stmt, err := db.Prepare("DELETE FROM likes  WHERE user_id=$1 AND book_id=$2")
    checkErr(err)
    user_id := c.PostForm("user_id")
    book_id := c.PostForm("book_id")
    stmt.Exec(user_id, book_id)

    stmt2, err := db.Prepare("UPDATE books SET liked=liked-1 WHERE id=$1")
    checkErr(err)
    stmt2.Exec(book_id)
    db.Close()
  }
}

func registerLike(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    stmt, err := db.Prepare("INSERT INTO likes(user_id, book_id) VALUES($1, $2)")
    checkErr(err)
    user_id := c.PostForm("user_id")
    book_id := c.PostForm("book_id")
    stmt.Exec(user_id, book_id)

    stmt2, err := db.Prepare("UPDATE books SET liked=liked+1 WHERE id=$1")
    checkErr(err)
    stmt2.Exec(book_id)
    db.Close()
  }
}

func bookDetail(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    id := c.Query("id")
    bookDetail := Book{}
    db.QueryRow("SELECT id, user_id, google_image, original_image, title, state, price, note, university, is_public, liked, created_at FROM books WHERE id=$1", id).Scan(&bookDetail.Id, &bookDetail.User_id, &bookDetail.Google_image, &bookDetail.Original_image, &bookDetail.Title, &bookDetail.State, &bookDetail.Price, &bookDetail.Note, &bookDetail.University, &bookDetail.Is_public, &bookDetail.Liked, &bookDetail.Created_at)
    db.Close()
    c.JSON(http.StatusOK, bookDetail)
  }
}

func books(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    university := c.Query("university")
    rows, _ := db.Query("SELECT id, user_id, google_image, original_image, title, state, price, note, liked, is_public, updated_at FROM books WHERE university=$1 ORDER BY created_at DESC", university)
    var books []Book
    for rows.Next(){
      b:=Book{}
      rows.Scan(&b.Id, &b.User_id, &b.Google_image, &b.Original_image, &b.Title, &b.State, &b.Price, &b.Note, &b.Liked, &b.Is_public, &b.Updated_at)
      books = append(books, b)
    }
    db.Close()
    c.JSON(http.StatusOK, books)
  }
}

func restartSelling(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    book_id := c.PostForm("book_id")
    db.Exec("UPDATE books SET is_public=true WHERE id=$1;", book_id)
    db.Close()
  }
}

func stopSelling(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    book_id := c.PostForm("book_id")
    db.Exec("UPDATE books SET is_public=false WHERE id=$1;", book_id)
    db.Close()
  }
}

func bookImages(c *gin.Context){
  original_image := c.Param("original_image")
  c.File("./book_images/"+original_image)
}

func mybooks(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    user_id := c.Query("user_id")
    rows, _ := db.Query("SELECT id, user_id, google_image, original_image, title, state, price, note, is_public, created_at FROM books WHERE user_id=$1 ORDER BY created_at DESC", user_id)
    var mybooks []Book
    for rows.Next(){
      b:=Book{}
      rows.Scan(&b.Id, &b.User_id, &b.Google_image, &b.Original_image, &b.Title, &b.State, &b.Price, &b.Note, &b.Is_public, &b.Created_at)
      mybooks = append(mybooks, b)
    }
    db.Close()
    c.JSON(http.StatusOK, mybooks)
  }
}

func sell(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    stmt, err := db.Prepare("INSERT INTO books(user_id, university, google_image, original_image, title, state, price, note, is_public) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id")
    checkErr(err)
    user_id := c.PostForm("user_id")
    university := c.PostForm("university")
    google_image := c.PostForm("google_image")
    original_image := ""
    title := c.PostForm("title")
    state := c.PostForm("state")
    price := c.PostForm("price")
    note := c.PostForm("note")
    is_public := true
    file, _ := c.FormFile("original_image")
    if(file!=nil){
      original_image = user_id+title+file.Filename
    c.SaveUploadedFile(file, "book_images/"+original_image)
    }

    stmt.Exec(user_id, university, google_image, original_image, title, state, price, note, is_public)
    db.Close()
  }

}

func getUser(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    id := c.Query("id")
    user:=User{}
    db.QueryRow("SELECT id, uid, username, university, profile_image, sns_image, is_signup_detail FROM users WHERE id=$1", id).Scan(&user.Id, &user.Uid, &user.Username, &user.University, &user.Profile_image, &user.Sns_image, &user.Is_signup_detail)
    db.Close()
    c.JSON(http.StatusOK, user)
  }
}


func user(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    uid := c.Query("uid")
    user:=User{}
    db.QueryRow("SELECT id, uid, username, university, profile_image, sns_image, is_signup_detail FROM users WHERE uid=$1", uid).Scan(&user.Id, &user.Uid, &user.Username, &user.University, &user.Profile_image, &user.Sns_image, &user.Is_signup_detail)
    db.Close()
    c.JSON(http.StatusOK, user)
  }
}

func signupWithImg(db *sql.DB) gin.HandlerFunc{
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

func signup(db *sql.DB) gin.HandlerFunc{
  return func (c *gin.Context){
    stmt, err := db.Prepare("INSERT INTO users(uid, username, university, profile_image, sns_image, is_signup_detail) VALUES($1, $2, $3, $4, $5, $6) RETURNING uid")
    checkErr(err)
    uid := c.PostForm("uid")
    username := c.PostForm("username")
    university := c.PostForm("university")
    profile_image := c.PostForm("profile_image")
    sns_image := c.PostForm("sns_image")
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

