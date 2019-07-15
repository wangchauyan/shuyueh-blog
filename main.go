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

type Post struct {
	Title   string
	Content template.HTML
}

func initServer() *gin.Engine {
	server := gin.Default()
	server.Use(gin.Logger())
	server.Delims("{{", "}}")

	// load html template
	server.LoadHTMLGlob("./templates/*.tmpl.html")

	return server
}

func setupRootAPI(server *gin.Engine) {
	// handle the root access
	server.GET("/", func(c *gin.Context) {
		var posts []string

		files, err := ioutil.ReadDir("./markdown/")
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			fmt.Println(file.Name())
			posts = append(posts, file.Name())
		}

		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"posts": posts,
		})
	})
}

func setupPostAPI(server *gin.Engine) {
	server.GET("/:postName", func(c *gin.Context) {
		postName := c.Param("postName")

		mdfile, err := ioutil.ReadFile("./markdown/" + postName)

		if err != nil {
			fmt.Println(err)
			// handle the error page access
			c.HTML(http.StatusNotFound, "error.tmpl.html", nil)
			c.Abort()
			return
		}

		postHTML := template.HTML(blackfriday.Run([]byte(mdfile)))

		post := Post{Title: postName, Content: postHTML}

		c.HTML(http.StatusOK, "post.tmpl.html", gin.H{
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

	server.Run(":" + os.Args[1])
}