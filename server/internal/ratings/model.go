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
		rm.Ratings[DEFAULT_RATING_VENDOR].InternalId,
		rm.Ratings[DEFAULT_RATING_VENDOR].Title,
		rm.Ratings[DEFAULT_RATING_VENDOR].Year,
		rm.Ratings[DEFAULT_RATING_VENDOR].Type,
		rm.Ratings[DEFAULT_RATING_VENDOR].Value,
		rm.Ratings[DEFAULT_RATING_VENDOR].Count,
	}, nil
}

func getRating(row []string) Rating {
	var InternalId, Title, Year, Type, Value, Count string
	if len(row) >= 4 {
		InternalId = row[3]
	}
	if len(row) >= 5 {
		Title = row[4]
	}
	if len(row) >= 6 {
		Year = row[5]
	}
	if len(row) >= 7 {
		Type = row[6]
	}
	if len(row) >= 8 {
		Value = row[7]
	}
	if len(row) >= 9 {
		Count = row[8]
	}

	return Rating{
		InternalId: InternalId,
		Title:      Title,
		Year:       Year,
		Type:       Type,
		Value:      Value,
		Count:      Count,
	}
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
		Ratings:    map[RatingVendor]Rating{DEFAULT_RATING_VENDOR: getRating(row)},
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
	Title      string
	Year       string
	Type       string
	Value      string
	Count      string
}

func ReadAll() ([]RatingModel, error) {
	return db.ReadAll()
}

func SaveAll(rm []RatingModel) error {
	return db.SaveAll(rm)
}

func ReadFirstEmpty() (RatingModel, error) {
	return db.ReadFirst(func(entity RatingModel) bool {
		if entity.Ratings[DEFAULT_RATING_VENDOR].InternalId == "" {
			return true
		}
		return false
	})
}
