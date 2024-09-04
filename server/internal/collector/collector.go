package collector

import (
	"fmt"
	"io"
	"net/http"
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
		fmt.Println("Waiting for next collect...")
		time.Sleep(5 * time.Second)

		emptyEntity, err := getEmptyEntity()
		if err != nil {
			fmt.Println(err)
		}
		getSearchResult(emptyEntity)

	}
}

func getEmptyEntity() (ratings.RatingModel, error) {
	return ratings.ReadFirstEmpty()
}

func getSearchResult(rm ratings.RatingModel) (any, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://www.filmweb.pl/api/v1/search?query=%s",
			url.QueryEscape(rm.Ratings[ratings.DEFAULT_RATING_VENDOR].Title),
		),
		nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Locale", "pl_PL")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	return nil, nil
}
