package ratings

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RatingsRoutes(router *gin.RouterGroup) {
	router.GET("/ratings", ratingsList)
}

func ratingsList(c *gin.Context) {
	query := c.Query("query")

	// Temporary save requested data
	querySplitted := strings.Split(query, "|")

	ratingsToSave := make([]RatingModel, 0)
	existingRatings, err := ReadAll()
	counter := uint(len(existingRatings))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if len(querySplitted) > 1 {
	outer:
		for idx, value := range querySplitted {
			if idx%2 == 0 {
				for _, rating := range existingRatings {
					if rating.Streamings[DEFAULT_STREAMING].InternalId == value {
						continue outer
					}
				}
				counter++
				ratingsToSave = append(ratingsToSave, RatingModel{
					Id:         counter,
					Streamings: map[StreamingName]StreamingModel{DEFAULT_STREAMING: StreamingModel{InternalId: value}},
					MovieData:  map[Language]MovieDataModel{DEFAULT_LANGUAGE: MovieDataModel{Title: querySplitted[idx+1]}},
				})
			}
		}

		err = SaveAll(ratingsToSave)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	res, err := ReadAll()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	serializer := RatingsSerializer{c, res}
	c.JSON(http.StatusOK, gin.H{"ratings": serializer.Response()})
}
