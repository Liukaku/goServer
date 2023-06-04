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

func GetById(c *gin.Context){
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	err := godotenv.Load()

	if err != nil {
		log.Print("env load error")
		log.Print(err.Error())
		return
	}

	routeId, pathBool := c.Params.Get("id")

	if pathBool == false {
		log.Print("path bool is false")
		c.JSON(500, gin.H{
			"response": "path bool is false",
		})
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
	selectPrep, err := db.Prepare("SELECT * FROM Quiz WHERE id = (?)")

	if err != nil {
		log.Print("db prepare error")
		log.Print(err.Error())
		c.JSON(500, gin.H{
			"response": "db prepare error",
		})
		return
	}

	rows, err := selectPrep.Query(routeId)

	if err != nil {
		log.Print("db select error")
		log.Print(err.Error())
		c.JSON(500, gin.H{
			"response": "db select error",
		})
		return
	}

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