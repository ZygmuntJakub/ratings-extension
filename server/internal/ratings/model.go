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
		rm.Ratings[DEFAULT_RATING_VENDOR].SubType,
		rm.Ratings[DEFAULT_RATING_VENDOR].OriginalTitle,
		rm.Ratings[DEFAULT_RATING_VENDOR].PosterPath,
		rm.Ratings[DEFAULT_RATING_VENDOR].Value,
		rm.Ratings[DEFAULT_RATING_VENDOR].Count,
		rm.Ratings[DEFAULT_RATING_VENDOR].WantToSee,
	}, nil
}

func getRating(row []string) Rating {
	var InternalId, Title, Year, Type, SubType, OriginalTitle, PosterPath, Value, Count, WantToSee string
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
		SubType = row[7]
	}
	if len(row) >= 9 {
		OriginalTitle = row[8]
	}
	if len(row) >= 10 {
		PosterPath = row[9]
	}
	if len(row) >= 11 {
		Value = row[10]
	}
	if len(row) >= 12 {
		Count = row[11]
	}
	if len(row) >= 13 {
		WantToSee = row[12]
	}

	return Rating{
		InternalId:    InternalId,
		Title:         Title,
		Year:          Year,
		Type:          Type,
		SubType:       SubType,
		OriginalTitle: OriginalTitle,
		PosterPath:    PosterPath,
		Value:         Value,
		Count:         Count,
		WantToSee:     WantToSee,
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

func (rm RatingModel) GetId() uint {
	return rm.Id
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
	InternalId    string // Rating company internal id
	Title         string
	Year          string
	Type          string
	SubType       string
	OriginalTitle string
	PosterPath    string
	Value         string
	Count         string
	WantToSee     string
}

func ReadAll() ([]RatingModel, error) {
	return db.ReadAll()
}

func ReadOnlyWithRating() ([]RatingModel, error) {
	return db.Read(func(entity RatingModel) bool {
		return entity.Ratings[DEFAULT_RATING_VENDOR].Count != ""
	})
}

func SaveAll(rm []RatingModel) error {
	return db.SaveAll(rm)
}

func ReadFirstEmpty(bq map[uint]bool) (RatingModel, error) {
	return db.ReadFirst(func(entity RatingModel) bool {
		if bq[entity.Id] { // Skip broken movies
			return false
		}
		if entity.Ratings[DEFAULT_RATING_VENDOR].Value == "" {
			return true
		}
		return false
	})
}

func Update(rm RatingModel) error {
	return db.Update(rm)
}
