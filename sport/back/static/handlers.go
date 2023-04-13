package static

import (
	"path"

	"github.com/gin-gonic/gin"
)

/*
NOTE THAT CURRENTLY STATIC SERVER IS NOT SECURE
AND IS MEANT TO BE USED __ONLY__ IN DEV ENV
*/

func (c *Ctx) GetFile() gin.HandlerFunc {
	return func(g *gin.Context) {
		if err := findPathTraversal(g.Params); err != nil {
			g.AbortWithError(400, err)
			return
		}
		path := path.Join(c.Config.Basepath, g.Param("dir"), g.Param("file"))

		g.Header("Content-Description", "File Transfer")
		g.Header("Content-Transfer-Encoding", "binary")
		g.Header("Content-Disposition", "attachment; filename="+g.Param("file"))
		g.File(path)
	}
}
