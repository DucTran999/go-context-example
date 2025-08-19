package main

import (
	"context"
	"errors"
	"log"
	"net/http"
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

func UpdateArticleHandler(c *gin.Context) {
	// ⏱ Set a timeout of 2 seconds for the request
	ctx, cancel := context.WithDeadline(c.Request.Context(), time.Now().Add(time.Second*2))
	defer cancel()

	id := c.Param("id") // get request param named id from gin.Context

	reqBody := new(UpdateArticleRequest)
	c.Bind(&reqBody) // get request body form gin.Context

	// pass context.Context to business layer to update the article
	perm, err := UpdateArticleBiz(ctx, id, reqBody.Content)
	if err != nil {
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

func UpdateArticleBiz(ctx context.Context, id, content string) (string, error) {
	// Simulate a slow operation
	select {
	case <-time.After(3 * time.Second):
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

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Apply middleware globally
	router.Use(AuthJWTMiddleware())
	router.PUT("/article/:id", UpdateArticleHandler)

	router.Run()
}
