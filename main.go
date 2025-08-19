package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateArticleRequest struct {
	Content string `json:"content"`
}

// AuthJWTMiddleware simulates JWT auth and pass role value via context
func AuthJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add context value "user_permissions" with value "admin"
		ctx := context.WithValue(c.Request.Context(), "user_permissions", "admin")

		//  Update the request with the new context
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func UpdateArticleBiz(ctx context.Context, id, content string) (string, error) {
	// Simulate a slow operation
	select {
	case <-time.After(5 * time.Second):
		// Retrieve the user permission from context
		value := ctx.Value("user_permissions")
		if valueStr, ok := value.(string); ok {
			return valueStr, nil
		}
		return "", errors.New("user_permission not found")
	case <-ctx.Done():
		log.Println("update article got error:", ctx.Err())
		return "", ctx.Err()
	}
}

func UpdateArticleHandler(c *gin.Context) {
	id := c.Param("id") // get request param named id from gin.Context

	reqBody := new(UpdateArticleRequest)
	c.Bind(&reqBody) // get request body form gin.Context

	perm, err := UpdateArticleBiz(c.Request.Context(), id, reqBody.Content)
	if err != nil {
		log.Println("[ERROR] update article failed: ", err.Error())
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":              id,
		"message":         "update success",
		"content":         reqBody.Content,
		"user_permission": perm,
	})
}

func main() {
	// Create new context listen from syscall interrupt or terminate
	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Apply middleware globally
	router.Use(AuthJWTMiddleware())
	router.PUT("/article/:id", UpdateArticleHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	// start HTTP on different go routine
	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// If got signal interrupt or terminate, gracefully shutdown the server
	<-appCtx.Done()
	log.Printf("shutdown signal received\n")

	// Give the server up to 10 seconds to finish processing ongoing requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
