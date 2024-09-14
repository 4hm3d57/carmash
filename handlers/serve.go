package handlers

import (
	"car/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
    "log"
)

type CarPair struct {
	CarID1   string
	CarID2   string
	CarLink1 string
	CarLink2 string
}

func ServePage(c *gin.Context) {
	// Get the car pair
	cars, err := db.GetPair()
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to get image pairs")
		return
	}

	if len(cars) < 2 {
		c.String(http.StatusInternalServerError, "not enough cars")
		return
	}

	carID1 := cars[0].ID.Hex()
	carID2 := cars[1].ID.Hex()

	// Prepare the car pair data
	carPairData := CarPair{
		CarID1:   carID1,
		CarID2:   carID2,
		CarLink1: fmt.Sprintf("https://drive.google.com/uc?id=%s", cars[0].FileID),
		CarLink2: fmt.Sprintf("https://drive.google.com/uc?id=%s", cars[1].FileID),
	}

    log.Printf("first image link: %s", carPairData.CarLink1)
    log.Printf("second image link: %s", carPairData.CarLink2)


	// Get all car data for points
	allCars, err := db.GetAllPics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving points"})
		return
	}

	// Render the template with both car pair and all car data
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to parse templates")
		return
	}

	data := gin.H{
		"carPair": carPairData,
		"cars":    allCars,
	}

	t.Execute(c.Writer, data)
}
