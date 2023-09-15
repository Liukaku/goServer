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

// TO DO:
// separate out the stucts to a different file
// break out a lot into helper functions so it's not one mega function
// 	^^--that will make it so much easier to ready this bastard file


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
	selectPrep, err := db.Prepare("SELECT * FROM Quiz LEFT OUTER JOIN Sections ON Sections.quiz_id = Quiz.id LEFT OUTER JOIN Questions ON Questions.section_id = Sections.id WHERE Quiz.id = (?);")

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
		log.Print("db select error 1")
		log.Print(err.Error())
		c.JSON(500, gin.H{
			"response": "db select error",
		})
		return
	}

	columns, err:= rows.Columns()
	if err != nil {
		log.Print("db select error 2")
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
		Section_id sql.NullInt16 `json:"section_id"`
		Section_title sql.NullString `json:"section_title"`
		Section_background sql.NullString `json:"section_background"`
		Quiz_id sql.NullInt64 `json:"quiz_id"`
		Question_id sql.NullInt64 `json:"question_id"`
		Question_title sql.NullString `json:"question_title"`
		Question_background sql.NullString `json:"question_background"`
		Question_type sql.NullString `json:"question_type"`
		From_section_id sql.NullInt64 `json:"from_section_id"`
		Order sql.NullInt64 `json:"order"`
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
			log.Print("db select error 3")
			log.Print(err.Error())
			return
		}

		retArr = append(retArr, ret)
	}

	p := RetObj{}
	q := map[string]Sections{}
	r := map[string][]Questions{}
	fakeAns := Answers{0, "SINGLE_CHOICE", true, 0, "answer one12"}
	ans := []Answers{fakeAns}
	fmt.Println(retArr)
	for _, retRow := range retArr {
		secId := retRow.From_section_id
		var setSecId int
		
		// if section id is valid that means there are sections and questions, so go ahead and populate the structs
		// if there aren't then don't
		// any questions aren't valid as they should be deleted with the sections so don't worry about that 
		if secId.Valid {
			setSecId = int(secId.Int64)
			strSecQid := strconv.Itoa(int(setSecId))
			if _, ok := q[strSecQid]; !ok{
				q[strSecQid] = Sections{setSecId, retRow.Section_title.String, retRow.Section_background.String}
			}
	
			qBg := retRow.Question_background 
	
			r[strSecQid] = append(r[strSecQid], Questions{int(retRow.Question_id.Int64), retRow.Question_title.String, qBg.String , retRow.Question_type.String, setSecId, int(retRow.Order.Int64), ans})
		}
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