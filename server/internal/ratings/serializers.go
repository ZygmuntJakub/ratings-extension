package ratings

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type RatingSerializer struct {
	C      *gin.Context
	Rating RatingModel
}

type RatingResponse struct {
	StreamingVendorId string
	Value             string
	Count             string
	Link              string
}

func (s *RatingSerializer) Response() RatingResponse {
	streamingVendor := s.Rating.Streamings[DEFAULT_STREAMING]
	rating := s.Rating.Ratings[DEFAULT_RATING_VENDOR]
	Link := fmt.Sprintf(
		"https://www.filmweb.pl/%s/%s-%s-%s",
		rating.Type,
		rating.Title,
		rating.Year,
		rating.InternalId,
	)

	return RatingResponse{
		StreamingVendorId: streamingVendor.InternalId,
		Value:             rating.Value,
		Count:             rating.Count,
		Link:              Link,
	}
}

type RatingsSerializer struct {
	C       *gin.Context
	Ratings []RatingModel
}

func (s *RatingsSerializer) Response() []RatingResponse {
	response := make([]RatingResponse, len(s.Ratings))

	for idx, r := range s.Ratings {
		serializer := RatingSerializer{s.C, r}
		response[idx] = serializer.Response()
	}

	return response
}
