package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"

	// "io/fs"
	"net/http"
	"publisher/openapi"

	"github.com/gin-gonic/gin"
	"github.com/mushoffa/go-library/redis"
	server "github.com/mushoffa/go-library/server/http"
)

type PublisherService interface {
	Publish(*gin.Context)
}

type PublisherHandler struct {
	r redis.RedisService
}

type Publisher struct {
	Topic string      `json:"topic" binding:"required"`
	Data  interface{} `json:"data" binding:"required"`
}

//go:embed static
var content embed.FS

func main() {

	// redisHost := "redis-pubsub:6379"
	redisHost := ":6380"
	redisPassword := ""
	redisDB := 0

	redis, err := redis.NewRedisClient(redisHost, redisPassword, redisDB)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	fmt.Println("Redis PING: ", redis.Ping())

	handler := PublisherHandler{redis}

	router := gin.Default()
	router.POST("/publish", handler.Publish)

	// fsys, _ := fs.Sub(content, "static")
	// swagger := http.StripPrefix("/static/", http.FileServer(http.FS(fsys)))
	// router.GET("/swagger/*any", gin.WrapH(swagger))
	// router.StaticFS("/swagger", http.FileServer(http.FS(fsys)))
	router.Static("/swagger", "static/swagger-ui/")
	// router.GET("/swagger/*any", http.StripPrefix("/static/", http.FileServer(http.FS(fsys))))

	openapi.RegisterOpenAPI(router)

	server := server.NewHttpServer(9001, router)
	server.Run()
}

func (c *PublisherHandler) Publish(ctx *gin.Context) {
	parent := context.Background()
	defer parent.Done()

	request := Publisher{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var b bytes.Buffer

	if err := json.NewEncoder(&b).Encode(request.Data); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	res := c.r.Publish(request.Topic, b.Bytes())
	if err := res.Err(); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
	return
}
