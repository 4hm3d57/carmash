package main

import (
	"car/db"
	"car/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	r := gin.Default()

	r.Static("/static", "./templates")
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/proxy", handlers.ProxyHandler)

	folderID := "10a3Ilc5o3YbBBdTVcqKUJJwVx9zJ7Re0" // Replace with your folder ID

	files, err := handlers.ListFilesInFolder(folderID)
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}

	for _, file := range files {
		if err := db.InsertImagesLinks(file); err != nil {
			log.Printf("Failed to insert file metadata for file %s: %v", file.Name, err)
		}
	}

	r.GET("/", handlers.ServePage)
	r.POST("/choose", handlers.SubmitChoice)

	log.Fatal(r.Run(":7070"))
}
