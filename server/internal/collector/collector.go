package collector

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/ZygmuntJakub/mkino-extension/internal/ratings"
)

type QueueEntity struct {
	id    uint
	title string
}

var (
	brokenQueue = make(map[uint]bool)
	TIME_OFFSET = 50
)

func beforeContinue(id uint, message string, err error) {
	fmt.Println(message)
	fmt.Println(err)
	brokenQueue[id] = true
	time.Sleep(time.Duration(rand.Intn(5)+TIME_OFFSET) * time.Second)
}

func RunCollector() {
	for {
		fmt.Println("Start collecting...")

		rm, err := getEmptyRating(brokenQueue)
		if err != nil {
			beforeContinue(0, "Cannot get next empty rating. Ending...", err)
			break
		}

		ratingId := rm.Ratings[ratings.DEFAULT_RATING_VENDOR].InternalId
		if ratingId == "" {
			srr, err := getSearchResult(rm)
			if err != nil {
				beforeContinue(rm.Id, "Cannot get search result.", err)
				continue
			}

			ratingId, err = getSearchResultId(srr.SearchHits)
			if err != nil {
				beforeContinue(rm.Id, "Cannot get search result id.", err)
				continue
			}

			irs, err := getInfo(ratingId)
			if err != nil {
				beforeContinue(rm.Id, "Cannot get info result.", err)
				continue
			}

			rm.Ratings[ratings.DEFAULT_RATING_VENDOR] = ratings.Rating{
				InternalId:    ratingId,
				Title:         irs.Title,
				Year:          strconv.FormatUint(uint64(irs.Year), 10),
				Type:          irs.Type,
				SubType:       irs.SubType,
				OriginalTitle: irs.OriginalTitle,
				PosterPath:    irs.PosterPath,
			}
		}

		if rm.Ratings[ratings.DEFAULT_RATING_VENDOR].Value == "" {
			fmt.Printf("Getting rating for: %s\n", rm.Ratings[ratings.DEFAULT_RATING_VENDOR].Title)
			r, err := getRating(ratingId)
			if err != nil {
				beforeContinue(rm.Id, "Cannot get rating.", err)
				continue
			}
			fmt.Println(r)
			res := rm.Ratings[ratings.DEFAULT_RATING_VENDOR]
			res.Value = strconv.FormatFloat(r.Rate, 'f', -1, 64)
			res.Count = strconv.FormatUint(uint64(r.Count), 10)
			res.WantToSee = strconv.FormatUint(uint64(r.WantToSee), 10)
			rm.Ratings[ratings.DEFAULT_RATING_VENDOR] = res
		}

		err = saveInfoResult(&rm)
		if err != nil {
			beforeContinue(rm.Id, "Cannot save info result.", err)
			continue
		}

		fmt.Println("Done. Wait for next request...")
		time.Sleep(time.Duration(rand.Intn(5)+TIME_OFFSET) * time.Second)
	}
}

func getEmptyRating(bq map[uint]bool) (ratings.RatingModel, error) {
	return ratings.ReadFirstEmpty(bq)
}

const (
	CHARACTER_TYPE = "character"
)

type SearchHit struct {
	Id   uint   `json:"id"`
	Type string `json:"type"`
}

type SearchResultResponse struct {
	SearchHits []SearchHit `json:"searchHits"`
}

func getSearchResult(rm ratings.RatingModel) (SearchResultResponse, error) {
	fmt.Printf("Looking for: %s\n", rm.MovieData[ratings.DEFAULT_LANGUAGE].Title)
	params := url.Values{}
	params.Add("query", rm.MovieData[ratings.DEFAULT_LANGUAGE].Title)
	url := fmt.Sprintf(
		"https://www.filmweb.pl/api/v1/search?%s",
		params.Encode(),
	)
	fmt.Printf("Search URL: %s\n", url)

	return MakeGetRequest[SearchResultResponse](url)
}

func getSearchResultId(sh []SearchHit) (string, error) {
	for _, value := range sh {
		if value.Type == CHARACTER_TYPE {
			continue
		}

		return strconv.FormatUint(uint64(value.Id), 10), nil
	}
	return "", errors.New("Cannot find movie/series type in search result")
}

type InfoResponse struct {
	Title         string `json:"title"`
	OriginalTitle string `json:"originalTitle"`
	Year          uint   `json:"year"`
	Type          string `json:"type"`
	SubType       string `json:"subType"`
	PosterPath    string `json:"posterPath"`
}

func getInfo(ratingId string) (InfoResponse, error) {
	fmt.Printf("Getting info for: %s\n", ratingId)
	url := fmt.Sprintf(
		"https://www.filmweb.pl/api/v1/title/%s/info",
		ratingId,
	)
	fmt.Printf("Search URL: %s\n", url)

	return MakeGetRequest[InfoResponse](url)
}

type RatingResponse struct {
	Count     uint    `json:"count"`
	Rate      float64 `json:"rate"`
	WantToSee uint    `json:"countWantToSee"`
}

func getRating(ratingId string) (RatingResponse, error) {
	fmt.Printf("Get rating: %s", ratingId)
	url := fmt.Sprintf(
		"https://www.filmweb.pl/api/v1/film/%s/rating",
		ratingId,
	)
	fmt.Printf("Get rating URL: %s\n", url)

	return MakeGetRequest[RatingResponse](url)
}

func saveInfoResult(rm *ratings.RatingModel) error {
	fmt.Printf("Saving info: %s\n", rm.Ratings[ratings.DEFAULT_RATING_VENDOR].Title)

	return ratings.Update(*rm)

}
