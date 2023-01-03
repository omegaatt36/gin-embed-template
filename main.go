package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omegeett36/gin-embed-template/templates"
)

func generateGeneralTemplate(c *gin.Context) {
	var req struct {
		NickName string `form:"nick_name" binding:"required"`
		FullName string `form:"full_name" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	mailText, err := templates.GenerateHTML(templates.TemplateNameGeneralTmpl, req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	log.Println(mailText)

	c.JSON(http.StatusOK, mailText)
}

func main() {
	router := gin.Default()
	templates.SetHTMLTemplate(router)

	router.GET("hello", generateGeneralTemplate)

	srv := &http.Server{
		Addr:    "localhost:8787",
		Handler: router,
	}

	log.Println("starts serving...")
	if err := srv.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}
