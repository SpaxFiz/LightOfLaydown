package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	JSON "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/spaxfiz/unjuanable/core/domain"
	"github.com/spaxfiz/unjuanable/core/storage"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

func main() {
	logrus.SetOutput(os.Stdout)

	r := gin.Default()
	route := r.Group("/api")
	{
		route.GET("/background", background)
		route.POST("/cipher", cipher)
		route.GET("/em_account_data", Handler(&domain.EmAccount{}))
		route.GET("/lg_pe_data", Handler(&domain.PETrend{}))
		route.GET("/industry_pe_data", Handler(&domain.IndustryPETrend{}))
		route.GET("/single_fund", domain.SingleFundHandler)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	//r.GET("/", func(context *gin.Context) {
	//	context.JSON(200, gin.H{"message": "hello world"})
	//})

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	<-interrupt
	if err := storage.GetCache().PersistCache(); err != nil {
		logrus.Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
}

func Handler(iface domain.CrawlData) func(*gin.Context) {
	return func(c *gin.Context) {
		if err := iface.Fetch(); err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}
		if data, e := iface.Render(); e != nil {
			c.JSON(500, gin.H{"message": e.Error()})
			return
		} else {
			c.JSON(200, data)
		}
	}
}

func cipher(c *gin.Context) {
	raw, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusForbidden)
		_, _ = c.Writer.WriteString("something went wrong")
		return
	}
	data := struct {
		Cipher string
	}{}
	if err := JSON.Unmarshal(raw, &data); err != nil {
		logrus.Error(err)
		c.Writer.WriteHeader(http.StatusForbidden)
		_, _ = c.Writer.WriteString("something went wrong")
	}
	if time.Now().Format("20060102") == data.Cipher {
		c.Writer.WriteHeader(http.StatusOK)
		return
	}
	c.Writer.WriteHeader(http.StatusForbidden)
	_, _ = c.Writer.WriteString("permission deny")
}

func background(c *gin.Context) {
	gopath := os.Getenv("GOPATH")
	picPath := path.Join(gopath, "src", "github.com", "spaxfiz", "unjuanable", "lib", "assets", "BG.jpg")
	fmt.Println(picPath)
	pic, _ := ioutil.ReadFile(picPath)
	fmt.Println(len(pic))

	c.Header("Content-Type", "image/jpg")
	c.Header("Content-Disposition", "attachment;filename=\"background.jpg\"")
	c.Data(http.StatusOK, "image/jpg", pic)
}
