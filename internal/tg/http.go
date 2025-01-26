package tg

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type response[T any] struct {
	Ok     bool `json:"ok"`
	Result T    `json:"result"`
}

func Request[T any](
	url string,
	in any,
	client *http.Client,
) (T, error) {
	var data T
	body, err := json.Marshal(in)
	if err != nil {
		return data, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return data, err
	}
	req.Header.Set("Content-Type", "application/json")

	httpresp, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer httpresp.Body.Close()

	var response response[T]

	if err := json.NewDecoder(httpresp.Body).Decode(&response); err != nil {
		return data, err
	}
	if !response.Ok {
		return data, errors.New("tg: recieved !ok from telegram")
	}

	return response.Result, nil
}
