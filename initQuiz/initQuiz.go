package initQuiz

import (
	"database/sql"
	"fmt"
	"os"

	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func InitQuiz(c *gin.Context){
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	err := godotenv.Load();

	if err != nil {
		log.Print("env load error")
		log.Print(err.Error())
		return
	}

	mysqlCred := os.Getenv("MYSQL_PWD")

	connectionString := fmt.Sprintf("root:%s@tcp(containers-us-west-166.railway.app:6421)/railway", mysqlCred)

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		log.Print("db connection err")
		log.Print(err.Error())
		return
	}
	defer db.Close()
	
	type ReqBody struct {
		Owner string `json:"owner"`
		Quiz_name string `json:"name"`
	}

	var initQuiz ReqBody

	if err := c.BindJSON(&initQuiz); err != nil {
		log.Print("json binding error")
		log.Print(err.Error())
		return
	}

	insert, dbErr := db.Prepare("INSERT INTO Users (username) VALUES (?)")

	if dbErr != nil {
		log.Print("db prepare err")
		log.Print(err.Error())
		return
	}

	res, err := insert.Exec(initQuiz.Owner)

	if err != nil {
		log.Print("db insert error")
		log.Print(err.Error())
		c.JSON(500, gin.H{
			"response": "db insert error",
		})
		return
	}

	lid, err := res.LastInsertId()

	if err != nil {
		log.Print("last ID error")
		log.Print(err.Error())
		return
	}
	
	quizInsert, dbErr := db.Prepare("INSERT INTO Quiz (quiz_title, owner_id) VALUES (?, ?)")

	if dbErr != nil {
		log.Print("db prepare error")
		log.Print(dbErr.Error())
		c.JSON(500, gin.H{
			"response": "db prepare error",
		})
		return
	}

	_, err = quizInsert.Exec(initQuiz.Quiz_name, lid)

	if err != nil {
		log.Print("db quiz insert error")
		log.Print(err.Error())
		c.JSON(500, gin.H{
			"response": "db quiz insert error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success!",
		"response": initQuiz,
		"newId": lid,
	})



}