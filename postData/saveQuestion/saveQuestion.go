package postData

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// this route is for updating existing questions
func SaveQuestion(c *gin.Context){
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	err := godotenv.Load();
	handleErr(c, err, "env load error")

	
	sqlConnStr := loadConnDetails()

	db, err := sql.Open("mysql", sqlConnStr)
	handleErr(c, err, "error connecting to db")

	defer db.Close()
	
	
}

func loadConnDetails() (string) {
	mysqlUn := os.Getenv("MYSQL_NAME")
	mysqlCred := os.Getenv("MYSQL_PWD")
	mysqlUrl := os.Getenv("DB_URL")

	connectionString := fmt.Sprintf("%s:%s@%s", mysqlUn, mysqlCred, mysqlUrl)

	return connectionString
}

func handleErr(c *gin.Context, err error, errStr string) {
	if err != nil {
		log.Print("env load error")
		log.Print(err.Error())
		c.JSON(500, gin.H{
			"response": "db prepare error",
		})
		panic(err.Error())
	}
}