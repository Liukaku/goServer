package indexPage

import (
	"database/sql"
	"fmt"
	"os"

	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func Index(c *gin.Context){
	
	err := godotenv.Load()

	if err != nil {
		log.Print("env load error")
		log.Print(err.Error())
		return
	}

	mysqlCred := os.Getenv("MYSQL_PWD")

	connectionString := fmt.Sprintf("root:%s@tcp(containers-us-west-166.railway.app:6421)/railway", mysqlCred)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Print("connection error")
		log.Print(err.Error())
		return
	}
	defer db.Close()
	log.Print(db)
	rows, dbErr := db.Query("SELECT * FROM Quiz")


	if dbErr != nil {
		log.Print("db select error")
		log.Print(dbErr.Error())
		c.JSON(500, gin.H{
			"response": "db select error",
		})
		return
	}
	log.Print(rows)

	columns, err:= rows.Columns()
	if err != nil {
		log.Print("db select error")
		log.Print(err.Error())
		c.JSON(500, gin.H{
			"response": "db select error",
		})
		return
	}


	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	// log.Print(columns)
	// log.Print(values)
	defer rows.Close()
		

	type Quiz struct {
		Quiz_id string `json:"quiz_id"`
		Quiz_title string `json:"quiz_title"`
		Owner_id string `json:"owner_id"`
	}
	quizArr := make([]any, 0)

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Print("db select error")
			log.Print(err.Error())
			return
		}
		returnObj := make(map[string]string)
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			returnObj[columns[i]] = value
			}

			// append to the end of the array
			quizArr = append(quizArr, returnObj)
}

	c.JSON(200, gin.H{
		"message": quizArr,
	})
}