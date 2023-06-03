package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/liukaku/go-server-ws/indexPage"
	"github.com/liukaku/go-server-ws/initQuiz"
	"github.com/liukaku/go-server-ws/postData"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	  },
}

func ws(c *gin.Context){
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer ws.Close()
	for { mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Print(mt)
		log.Printf("recv: %s", message)
		replyMsg:= []byte("hello, client!")
		ws.WriteMessage(mt, replyMsg)
	}
}

func main() {
	err := godotenv.Load();

	if err != nil {
		log.Print(".env error")
	}

	fmt.Println("websocket server start!")

	PORT := os.Getenv("PORT")

	bindAddress := fmt.Sprintf(":%s", PORT)

	r:= gin.Default()
	r.GET("/ws", ws)
	r.GET("/", indexPage.Index)
	r.POST("/createQuiz", postData.CreateQuiz)
	r.POST("/initQuiz", initQuiz.InitQuiz)
	r.Run(bindAddress)

	// r := gin.Default()
	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "Hello, World!",
	// 	})
	// })
	// r.Run()
}