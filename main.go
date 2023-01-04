package main

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/omegeett36/gin-embed-template/random"
	"github.com/omegeett36/gin-embed-template/templates"
)

type service struct {
	sync.Mutex

	m map[string]string
}

func newService() *service {
	return &service{
		m: make(map[string]string),
	}
}

func (s *service) generateToken(c *gin.Context) {
	var req struct {
		Name string `form:"name" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	/*
		mail to xxx.
		mailText, err := templates.GenerateHTML(templates.TemplateNameGeneralTmpl, req)
		if err != nil {
		 	c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err := mailService.MailTo(ooxx, mailText)
		if err != nil {
		 	c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, nil)
	*/

	var token string
	s.Lock()
	for {
		token = random.RandStringBytesMaskImprSrc(16)
		_, ok := s.m[token]
		if ok {
			continue
		}

		s.m[token] = req.Name
		break
	}

	s.Unlock()

	c.String(http.StatusOK, "http://localhost:8787/verify?token=%s", token)
}

func (s *service) verify(c *gin.Context) {
	var req struct {
		Token string `form:"token" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if _, ok := s.m[req.Token]; !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("bad token"))
		return
	}

	c.HTML(
		http.StatusOK,
		templates.TemplateNameVerifyTmpl.String(),
		map[string]any{"Name": s.m[req.Token]},
	)
}

func main() {
	router := gin.Default()
	templates.SetHTMLTemplate(router)

	s := newService()

	router.GET("token", s.generateToken)
	router.GET("verify", s.verify)

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
