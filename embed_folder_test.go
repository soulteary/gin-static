package static_test

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	static "github.com/soulteary/gin-static"
	"github.com/stretchr/testify/assert"
)

//go:embed test/data/server
var testFS embed.FS

func TestEmbedFolderWithRedir(t *testing.T) {
	var tests = []struct {
		targetURL string // input
		httpCode  int    // expected http code
		httpBody  string // expected http body
		name      string // test name
	}{
		{"/404.html", 301, "<a href=\"/\">Moved Permanently</a>.\n\n", "Unknown file"},
		{"/", 200, "<h1>Hello Embed</h1>", "Root"},
		{"/index.html", 301, "", "Root by file name automatic redirect"},
		{"/static.html", 200, "<h1>Hello Gin Static</h1>", "Other file"},
	}

	router := gin.New()

	staticFiles, err := static.EmbedFolder(testFS, "test/data/server")
	if err != nil {
		log.Fatalln("initialization of embed folder failed:", err)
	} else {
		router.Use(static.Serve("/", staticFiles))
	}

	router.NoRoute(func(c *gin.Context) {
		fmt.Printf("%s doesn't exists, redirect on /\n", c.Request.URL.Path)
		c.Redirect(301, "/")
	})

	for _, tt := range tests {
		w := PerformRequest(router, "GET", tt.targetURL)
		assert.Equal(t, tt.httpCode, w.Code, tt.name)
		assert.Equal(t, tt.httpBody, w.Body.String(), tt.name)
	}
}

func TestEmbedFolderWithoutRedir(t *testing.T) {
	var tests = []struct {
		targetURL string // input
		httpCode  int    // expected http code
		httpBody  string // expected http body
		name      string // test name
	}{
		{"/404.html", 404, "404 page not found", "Unknown file"},
		{"/", 200, "<h1>Hello Embed</h1>", "Root"},
		{"/index.html", 301, "", "Root by file name automatic redirect"},
		{"/static.html", 200, "<h1>Hello Gin Static</h1>", "Other file"},
	}

	router := gin.New()
	staticFiles, err := static.EmbedFolder(testFS, "test/data/server")
	if err != nil {
		log.Fatalln("initialization of embed folder failed:", err)
	} else {
		router.Use(static.Serve("/", staticFiles))
	}

	for _, tt := range tests {
		w := PerformRequest(router, "GET", tt.targetURL)
		assert.Equal(t, tt.httpCode, w.Code, tt.name)
		assert.Equal(t, tt.httpBody, w.Body.String(), tt.name)
	}
}

func TestEmbedInitErrorPath(t *testing.T) {
	tests := []struct {
		name       string
		targetPath string
		haveErr    bool
		fs         embed.FS
	}{
		{
			name:       "ValidPath",
			targetPath: "test/data/server",
			haveErr:    false,
			fs:         testFS,
		},
		{
			name:       "InvalidPath",
			targetPath: "nonexistingdirectory/nonexistingdirectory",
			haveErr:    true,
			fs:         testFS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := static.EmbedFolder(tt.fs, tt.targetPath)
			assert.Equal(t, (err != nil), tt.haveErr, tt.name)
		})
	}
}

func TestCreateEmbed(t *testing.T) {
	_, err := static.EmbedFolder(testFS, "test/data/server")
	if err != nil {
		log.Fatalln("initialization of embed folder failed:", err)
	}
	_, err2 := static.EmbedFolder(testFS, "")
	if err2 != nil {
		log.Fatalln("initialization of embed folder failed:", err)
	}
}

func TestServEmbed(t *testing.T) {
	var tests = []struct {
		targetURL string // input
		httpCode  int    // expected http code
		httpBody  string // expected http body
		name      string // test name
	}{
		{"/404.html", 404, "404 page not found\n", "Unknown file"},
		{"/", 200, "<h1>Hello Embed</h1>", "Root"},
		{"/index.html", 301, "<a href=\"/\">Moved Permanently</a>.\n\n", "Root by file name automatic redirect"},
		{"/static.html", 200, "<h1>Hello Gin Static</h1>", "Other file"},
	}

	router := gin.New()

	router.Use(static.ServeEmbed("test/data/server", testFS))

	router.NoRoute(func(c *gin.Context) {
		fmt.Printf("%s doesn't exists, redirect on /\n", c.Request.URL.Path)
		c.Redirect(301, "/")
	})

	for _, tt := range tests {
		w := PerformRequest(router, "GET", tt.targetURL)
		assert.Equal(t, tt.httpCode, w.Code, tt.name)
		assert.Equal(t, tt.httpBody, w.Body.String(), tt.name)
	}

	router2 := gin.New()

	router2.Use(static.ServeEmbed("", testFS))

	router2.NoRoute(func(c *gin.Context) {
		fmt.Printf("%s doesn't exists, redirect on /\n", c.Request.URL.Path)
		c.Redirect(301, "/")
	})

	for _, tt := range tests {
		w := PerformRequest(router2, "GET", "/test/data/server"+tt.targetURL)
		assert.Equal(t, tt.httpCode, w.Code, tt.name)
		assert.Equal(t, tt.httpBody, w.Body.String(), tt.name)
	}
}

func TestServEmbedErr(t *testing.T) {
	tests := []struct {
		name   string
		URL    string
		Code   int
		Result string
	}{
		{
			name:   "Invalid Path",
			URL:    "/test/data/server/nonexistingdirectory/nonexistingdirectory",
			Code:   http.StatusInternalServerError,
			Result: "{\"error\":\"open 111111test/data/server: file does not exist\",\"message\":\"initialization of embed folder failed\"}",
		},
	}

	router := gin.New()
	router.Use(static.ServeEmbed("111111test/data/server", testFS))
	for _, tt := range tests {
		w := PerformRequest(router, "GET", tt.URL)
		assert.Equal(t, w.Code, tt.Code)
		assert.Equal(t, w.Body.String(), tt.Result)
	}
}
