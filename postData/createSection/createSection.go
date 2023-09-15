package createSection

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CreateSection(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	err := godotenv.Load()

	handleErr(err, c)

	mysqlUn := os.Getenv("MYSQL_NAME")
	mysqlCred := os.Getenv("MYSQL_PWD")
	mysqlUrl := os.Getenv("DB_URL")

	connectionString := fmt.Sprintf("%s:%s@%s", mysqlUn, mysqlCred, mysqlUrl)

	db, err := sql.Open("mysql", connectionString)

	handleErr(err, c)

	defer db.Close()

	reqBody := map[string]string{}

	handleErr(err, c)

	fmt.Println(reqBody)

	rowExistsPrep, err := db.Prepare("SELECT * FROM Sections WHERE ID = (?)")

	handleErr(err, c)

	rows, err := rowExistsPrep.Query(reqBody["secId"])

	handleErr(err, c)

	i := 0
	for rows.Next() {
		i++
	}

	if i == 0 {
		createNewSection(reqBody, c)
	} else {
		editSection(reqBody, c)
	}

}

func handleErr (err error, c *gin.Context){
	if err != nil {
		fmt.Println("oops", err.Error())
		c.JSON(500, gin.H{
			"res": err.Error(),
		})
		return
	}
}

func createNewSection(reqBody map[string]string, c *gin.Context){
	
}

func editSection(reqBody map[string]string, c *gin.Context){

}