package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const markdownPath =  "./markdowns/"
const templatePath =  "./templates/"

type Post struct {
	Title   string
	Content template.HTML
}

func initServer() *gin.Engine {
	server := gin.Default()
	server.Use(gin.Logger())
	server.Delims("{{", "}}")

	// load html template files
	server.LoadHTMLGlob(templatePath + "*.html")

	return server
}

func setupRootAPI(server *gin.Engine) {
	// handle the root access
	server.GET("/", func(c *gin.Context) {
		var posts []string

		files, err := ioutil.ReadDir(markdownPath)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			fmt.Println(file.Name())
			posts = append(posts, file.Name())
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"posts": posts,
		})
	})
}

func setupPostAPI(server *gin.Engine) {
	server.GET("/:postName", func(c *gin.Context) {
		postName := c.Param("postName")

		mdfile, err := ioutil.ReadFile(markdownPath + postName)

		if err != nil {
			fmt.Println(err)
			// handle the error page access
			c.HTML(http.StatusNotFound, "error.html", nil)
			c.Abort()
			return
		}

		postHTML := template.HTML(blackfriday.Run([]byte(mdfile)))

		post := Post{Title: postName, Content: postHTML}

		c.HTML(http.StatusOK, "post.html", gin.H{
			"Title":   post.Title,
			"Content": post.Content,
		})
	})
}


func main()  {


	server := initServer()

	// setup home API
	setupRootAPI(server)

	// handle the post access
	setupPostAPI(server)

	var portNumber string
	if len(os.Args) == 1 {
		portNumber = "8000"
	} else {
		portNumber = os.Args[1]
	}
 	server.Run(":" + portNumber)
}