package static

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	return err == nil
}

func EmbedFolder(fsEmbed embed.FS, reqPath string) (ServeFileSystem, error) {
	targetPath := strings.TrimSpace(reqPath)
	if targetPath == "" {
		return embedFileSystem{
			FileSystem: http.FS(fsEmbed),
		}, nil
	}

	fsys, _ := fs.Sub(fsEmbed, targetPath)
	_, err := fsEmbed.Open(targetPath)
	if err != nil {
		return nil, err
	}

	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}, nil
}

func ServeEmbed(reqPath string, fsEmbed embed.FS) gin.HandlerFunc {
	embedFS, err := EmbedFolder(fsEmbed, reqPath)
	if err != nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "initialization of embed folder failed",
				"error":   err.Error(),
			})
		}
	}
	return gin.WrapH(http.FileServer(embedFS))
}
