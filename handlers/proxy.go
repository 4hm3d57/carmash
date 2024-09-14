package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func ProxyHandler(c *gin.Context) {
	imageURL := c.Query("url")
	if imageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing URL"})
		return
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Follow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}

	resp, err := client.Get(imageURL)
	if err != nil {
		fmt.Printf("Error fetching URL: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch image"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Received status code: %d\n", resp.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch image, status code: " + resp.Status})
		return
	}

	// Set headers
	for name, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(name, value)
		}
	}
	c.Writer.WriteHeader(resp.StatusCode)

	// Copy the response body to the gin response writer
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		fmt.Printf("Error copying response body: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy response body"})
	}
}
