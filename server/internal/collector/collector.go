package collector

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ZygmuntJakub/mkino-extension/internal/ratings"
)

type QueueEntity struct {
	id    uint
	title string
}

func RunCollector() {
	for {
		fmt.Println("Start collecting...")

		emptyEntity, err := getEmptyEntity()
		if err != nil {
			fmt.Println(err)
		}

		searchResultResponse, err := getSearchResult(emptyEntity)
		if err != nil {
			fmt.Println(err)
		}

		infoResultResponse, err := getInfoResult(searchResultResponse.SearchHits[0].Id)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(infoResultResponse)

		fmt.Println("Wait for next request...")
		time.Sleep(5 * time.Second)
	}
}

func getEmptyEntity() (ratings.RatingModel, error) {
	return ratings.ReadFirstEmpty()
}

type SearchHit struct {
	Id uint `json:"id"`
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

type InfoResultResponse struct {
	Title         string `json:"title"`
	OriginalTitle string `json:"originalTitle"`
	Year          uint   `json:"year"`
	Type          string `json:"type"`
	SubType       string `json:"subType"`
	PosterPath    string `json:"posterPath"`
}

func getInfoResult(id uint) (InfoResultResponse, error) {
	fmt.Printf("Getting info for: %d\n", id)
	url := fmt.Sprintf(
		"https://www.filmweb.pl/api/v1/title/%d/info",
		id,
	)
	fmt.Printf("Search URL: %s\n", url)

	return MakeGetRequest[InfoResultResponse](url)
}

func saveInfoResult(irs InfoResultResponse) error {
	return nil
}
