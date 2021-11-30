package main

import (
	"SpaxFiz/LaydownLight/core/domain"
	"SpaxFiz/LaydownLight/core/storage"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	r := gin.Default()
	route := r.Group("/api")
	{
		route.GET("/em_account_data", Handler(&domain.EmAccount{}))
		route.GET("/lg_pe_data", Handler(&domain.PETrend{}))
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
		log.Fatal("Server Shutdown:", err)
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
