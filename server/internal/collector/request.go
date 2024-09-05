package collector

import (
	"encoding/json"
	"io"
	"net/http"
)

func MakeGetRequest[T any](url string) (T, error) {
	var response T
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return response, err
	}

	req.Header.Set("X-Locale", "pl_PL")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
