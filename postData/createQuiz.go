package postData

import (
	"database/sql"
	"fmt"
	"os"

	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func CreateQuiz(c *gin.Context){
	
	err := godotenv.Load()

	if err != nil {
		log.Print("env load error")
		log.Print(err.Error())
		return
	}

	type Quiz struct {
		Quiz_title string `json:"quiz_title"`
		Owner_id int `json:"owner_id"`
	}

	var newQuiz Quiz
	
	mysqlUn := os.Getenv("MYSQL_NAME")
	mysqlCred := os.Getenv("MYSQL_PWD")
	mysqlUrl := os.Getenv("DB_URL")

	connectionString := fmt.Sprintf("%s:%s@%s", mysqlUn, mysqlCred, mysqlUrl)
	
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Print("connection error")
		log.Print(err.Error())
		return
	}

	defer db.Close()
	
	if err := c.BindJSON(&newQuiz); err != nil{
		log.Print("json binding error")
		log.Print(err.Error())
		return
	}

	insert, dbErr := db.Prepare("INSERT INTO Quiz (quiz_title, owner_id) VALUES (?, ?)")
	
	if dbErr != nil {
		log.Print("db prepare error")
		log.Print(dbErr.Error())
		c.JSON(500, gin.H{
			"response": "db prepare error",
		})
		return
	}

	_, err = insert.Exec(newQuiz.Quiz_title, newQuiz.Owner_id)

	if err != nil {
		log.Print("db insert error")
		log.Print(err.Error())
		c.JSON(500, gin.H{
			"response": "db insert error",
		})
		return
	}

	log.Print(db)
	c.JSON(200, gin.H{
		"message": "success!",
		"response": newQuiz,
	})

}