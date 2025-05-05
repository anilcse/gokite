package web

import (
	"net/http"

	"github.com/anilcse/gokite/internal/engine"

	"github.com/gin-gonic/gin"
)

// Router builds a minimal REST API to list rules and toggle enable/disable.
func Router(eng *engine.Engine) http.Handler {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	r.GET("/rules", func(c *gin.Context) {
		// naive exposure of in‑memory rule slice
		rules := eng.AllRules()
		c.JSON(200, rules)
	})

	// TODO: add POST /rules, PATCH /rules/:id etc.

	return r
}
