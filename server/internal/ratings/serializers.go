package ratings

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
)

type RatingSerializer struct {
	C      *gin.Context
	Rating RatingModel
}

type RatingResponse struct {
	StreamingVendorId string
	Value             string
	Link              string
}

func (s *RatingSerializer) Response() RatingResponse {
	streamingVendor := s.Rating.Streamings[DEFAULT_STREAMING]
	rating := s.Rating.Ratings[DEFAULT_RATING_VENDOR]
	Link := fmt.Sprintf(
		"https://www.filmweb.pl/%s/%s-%s-%s",
		rating.Type,
		url.QueryEscape(rating.Title),
		rating.Year,
		rating.InternalId,
	)

	return RatingResponse{
		StreamingVendorId: streamingVendor.InternalId,
		Value:             rating.Value,
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
