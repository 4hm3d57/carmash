package handlers

import (
	"math"
)

func base_rating(ratingA, ratingB float64) float64 {

	return 1 / (1 + math.Pow(10, (ratingB-ratingA)/400))

}

func CalculateEloRating(wRating, lRating, kFactor float64) (float64, float64) {

	expected_w := base_rating(wRating, lRating)
	expected_l := base_rating(lRating, wRating)

	newWinner := expected_w + kFactor*(1-expected_w)
	newLoser := expected_l + kFactor*(0-expected_l)

	return newWinner, newLoser

}
