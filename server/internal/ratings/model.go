package ratings

type RatingModel struct {
	Id         uint
	MovieData  map[Language]MovieDataModel
	Streamings map[StreamingName]StreamingModel
	Ratings    map[RatingVendor]Rating
}

type Language string

type MovieDataModel struct {
	Title string
}

type StreamingName string

type StreamingModel struct {
	InternalId string // Streaming company internal id
}

type RatingVendor string

type Rating struct {
	value float64
	count uint
}
