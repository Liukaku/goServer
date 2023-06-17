package indexPage

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)


type Sections struct {
	Id int `json:"section_id"`
	Section_title string `json:"name"`
	Section_bg string `json:"background"`
}

type Questions struct {
	Id int `json:"question_id"`
	Question_title string `json:"questionTitle"`
	Question_bg string `json:"question_bg"`
	Question_type string `json:"type"`
	Section_id int `json:"from_section_id"`
	Order int `json:"order"`
	Answer []Answers `json:"answer"`
}

type Answers struct {
	Id int `json:"id"`
	AnswerType string `json:"answerType"`
	Correct bool `json:"correct"`
	Order int `json:"order"`
	Title string `json:"title"`
}

type RetObj struct {
	Owner Owner `json:"Owner"`
	Sections map[string]Sections `json:"Sections"`
	Questions map[string][]Questions `json:"Questions"`
}

type Owner struct {
	OwnerName int `json:"ownerName"`
	QuizName string `json:"quizName"`
	Id int `json:"id"`
}

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
	log.Print(db)
	selectPrep, err := db.Prepare("SELECT * FROM Quiz INNER JOIN Sections ON Quiz.id = (?) INNER JOIN Questions ON Questions.section_id = Sections.id;")

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
		Order int `json:"order"`
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
			&ret.Order,
		)

		if err != nil {
			log.Print("db select error")
			log.Print(err.Error())
			return
		}

		retArr = append(retArr, ret)
	}


	p := RetObj{}
	q := map[string]Sections{}
	r := map[string][]Questions{}
	fakeAns := Answers{0, "SINGLE_CHOICE", true, 0, "answer one"}
	ans := []Answers{fakeAns}
	fmt.Println(retArr)
	for _, retRow := range retArr {
		strSecQid := strconv.Itoa(retRow.From_section_id)
		if _, ok := q[strSecQid]; !ok{
			q[strSecQid] = Sections{retRow.Section_id, retRow.Section_title, retRow.Section_background}
		}

		qBg := retRow.Question_background 

		if qBg == nil {
			var alo = "NULL"
			qBg = &alo
		}

		r[strSecQid] = append(r[strSecQid], Questions{retRow.Question_id, retRow.Question_title, *qBg , retRow.Question_type, retRow.From_section_id, retRow.Order, ans})
	}

	owner := Owner{retArr[0].Owner_id, retArr[0].Quiz_title, retArr[0].Id}

	p = RetObj{owner, q, r}

	c.JSON(200, gin.H{
		"message": p,
	})
}

func arrayContains(sections []Sections, secId int)bool{

	for _, section := range sections {
		if section.Id == secId {
			return true
		}
	}

	return false
}