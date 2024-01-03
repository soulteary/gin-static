package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed public
var EmbedFS embed.FS

func main() {
	r := gin.Default()

	// method 1: use as Gin Router
	// trim embedfs path `public/page`, and use it as url path `/`
	// r.GET("/", static.ServeEmbed("public/page", EmbedFS))

	// method 2: use as middleware
	// trim embedfs path `public/page`, the embedfs path start with `/`
	// r.Use(static.ServeEmbed("public/page", EmbedFS))

	// method 2.1: use as middleware
	// trim embedfs path `public/page`, the embedfs path start with `/public/page`
	// r.Use(static.ServeEmbed("", EmbedFS))

	// method 3: use as manual
	// trim embedfs path `public/page`, the embedfs path start with `/public/page`
	// staticFiles, err := static.EmbedFolder(EmbedFS, "public/page")
	// if err != nil {
	// 	log.Fatalln("initialization of embed folder failed:", err)
	// } else {
	// 	r.Use(static.Serve("/", staticFiles))
	// }

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "test")
	})

	r.NoRoute(func(c *gin.Context) {
		fmt.Printf("%s doesn't exists, redirect on /\n", c.Request.URL.Path)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
