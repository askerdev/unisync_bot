package mospolytech

import (
	"encoding/json"
	"net/http"
)

type api struct {
	scheduleURL string
	client      *http.Client
}

func NewAPI(
	scheduleURL string,
	client *http.Client,
) *api {
	return &api{
		scheduleURL: scheduleURL,
		client:      client,
	}
}

func (a *api) Schedule() (*SemesterSchedule, error) {
	httpresp, err := a.client.Get(a.scheduleURL)
	if err != nil {
		return nil, err
	}
	defer httpresp.Body.Close()

	var sch SemesterSchedule
	if err := json.NewDecoder(httpresp.Body).Decode(&sch); err != nil {
		return nil, err
	}

	return &sch, nil
}
