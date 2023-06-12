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
	selectPrep, err := db.Prepare("SELECT * FROM Quiz INNER JOIN Sections ON Quiz.id = (?) INNER JOIN Questions ON Sections.quiz_id = Questions.section_id;")

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

	type ReturnObj struct {
		Id int `json:"id"`
		Quiz_title string `json:"quiz_title"`
		Owner_id int `json:"owner_id"`
		Section_id int `json:"section_id"`
		Section_title string `json:"section_title"`
		Section_background string `json:"section_background"`
		Quiz_id int `json:"quiz_id"`
		Question_id int `json:"question_id"`
		Question_title string `json:"question_title"`
		Question_background *string `json:"question_background"`
		Question_type string `json:"question_type"`
		From_section_id int `json:"from_section_id"`
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	defer rows.Close()

	retArr := []ReturnObj{}
	for rows.Next() {
		var ret ReturnObj
		fmt.Println(ret, " bbb")
		err = rows.Scan(
			&ret.Id,
			&ret.Quiz_title,
			&ret.Owner_id,
			&ret.Section_id,
			&ret.Section_title,
			&ret.Section_background,
			&ret.Quiz_id,
			&ret.Question_id,
			&ret.Question_title,
			&ret.Question_background,
			&ret.Question_type,
			&ret.From_section_id,
		)

		if err != nil {
			log.Print("db select error")
			log.Print(err.Error())
			return
		}

		retArr = append(retArr, ret)
	}

	c.JSON(200, gin.H{
		"message": retArr,
	})
}