package ratings

import (
	"strconv"

	"github.com/ZygmuntJakub/mkino-extension/internal/database"
)

var db = database.DB[RatingModel]{
	File:           "/internal/ratings/db.csv",
	RowSerialize:   RowSerialize,
	RowDeserialize: RowDeserialize,
}

func RowSerialize(rm RatingModel) ([]string, error) {
	return []string{
		strconv.FormatUint(uint64(rm.Id), 10),
		rm.MovieData[DEFAULT_LANGUAGE].Title,
		rm.Streamings[DEFAULT_STREAMING].InternalId,
	}, nil
}

func RowDeserialize(row []string) (RatingModel, error) {
	Id, err := strconv.ParseUint(row[0], 10, 32)
	if err != nil {
		return RatingModel{}, err
	}

	return RatingModel{
		Id:         uint(Id),
		MovieData:  map[Language]MovieDataModel{"pl-PL": MovieDataModel{Title: row[1]}},
		Streamings: map[StreamingName]StreamingModel{DEFAULT_STREAMING: StreamingModel{InternalId: row[2]}},
	}, nil
}

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
	InternalId string // Rating company internal id
	Value      string
	Count      string
	Url        string
}

func ReadAll() ([]RatingModel, error) {
	return db.ReadAll()
}

func SaveAll(rm []RatingModel) error {
	return db.SaveAll(rm)
}
