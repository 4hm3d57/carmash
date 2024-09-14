package handlers

import (
	"car/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChoiceRequest struct {
	WinnerCarID string `form:"winnerCarID"`
	LoserCarID  string `form:"loserCarID"`
}

func SubmitChoice(c *gin.Context) {
	var choice ChoiceRequest
	if err := c.ShouldBind(&choice); err != nil {
		c.String(http.StatusBadRequest, "Invalid form submission")
		return
	}

	winnerCar, err := db.GetCarID(choice.WinnerCarID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to get winner car")
		return
	}
	loserCar, err := db.GetCarID(choice.LoserCarID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to get loser car")
		return
	}

	newWinnerRating, newLoserRating := CalculateEloRating(winnerCar.Elorating, loserCar.Elorating, 32)

	err = db.UpdateEloRating(winnerCar.ID.Hex(), loserCar.ID.Hex(), newWinnerRating, newLoserRating)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to update Elo ratings")
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}
