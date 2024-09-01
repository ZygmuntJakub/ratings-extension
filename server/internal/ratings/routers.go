package ratings

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RatingsRoutes(router *gin.RouterGroup) {
	router.GET("/ratings", ratingsList)
}

func ratingsList(c *gin.Context) {
	result := []RatingModel{
		{
			Id:         1,
			MovieData:  map[Language]MovieDataModel{"pl-PL": MovieDataModel{Title: "Biuro"}},
			Streamings: map[StreamingName]StreamingModel{DEFAULT_STREAMING: StreamingModel{InternalId: "1"}},
			Ratings:    map[RatingVendor]Rating{DEFAULT_RATING_VENDOR: Rating{value: 7.62, count: 12_322}},
		},
		{
			Id:         2,
			MovieData:  map[Language]MovieDataModel{"pl-PL": MovieDataModel{Title: "Kogel Mogel"}},
			Streamings: map[StreamingName]StreamingModel{DEFAULT_STREAMING: StreamingModel{InternalId: "2"}},
			Ratings:    map[RatingVendor]Rating{DEFAULT_RATING_VENDOR: Rating{value: 3.64, count: 1_322}},
		},
	}

	serializer := RatingsSerializer{c, result}
	c.JSON(http.StatusOK, gin.H{"ratings": serializer.Response()})
}
