package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/testutils/config"
)

// GetMostRecentValueForMetric returns the most recent value as a json.Number for the metric identified by the selector or an error if
func GetMostRecentValueForMetric(s selector, cfg config.Config) (*json.Number, error) {
	log.Printf("Getting metrics for %s", s.String())
	for cnt := 1; cnt <= cfg.MetricWaitRetries; cnt++ {
		metrics, err := GetMetrics(s, cfg.Endpoint, cfg.APIToken)
		if err != nil {
			return nil, err
		}
		if len(metrics) > 0 {
			log.Printf("Got result: %+v", metrics)
			return &metrics[len(metrics)-1], nil
		}
		log.Printf("Attempt %d / %d no results", cnt, cfg.MetricWaitRetries)
		time.Sleep(time.Duration(cfg.MetricRetryIntervalSeconds) * time.Second)
	}

	return nil, fmt.Errorf("failed to get metrics for %s", s.String())
}

func GetMetrics(s selector, endpoint, apiToken string) ([]json.Number, error) {
	req, err := http.NewRequest("GET", endpoint, bytes.NewBufferString(""))

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("accept", "application/json; charset=utf-8")
	req.Header.Add("Authorization", fmt.Sprintf("Api-Token %s", apiToken))

	q := req.URL.Query()
	q.Add("metricSelector", s.String())
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 400 {
		// metric is simply not yet available
		return []json.Number{}, nil
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("could not get metric: %s", resp.Status)
	}

	responseBody := metricQueryResponse{}
	if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %s", err.Error())
	}

	if count, err := responseBody.TotalCount.Int64(); err != nil {
		return nil, err
	} else if count == 0 {
		return []json.Number{}, nil
	}

	return responseBody.Result[0].Data[0].Values, nil
}

type metricQueryResponse struct {
	TotalCount json.Number          `json:"totalCount"`
	Result     []metricsQueryResult `json:"result"`
}

type metricsQueryResult struct {
	ID   string             `json:"metricId"`
	Data []metricResultData `json:"data"`
}

type metricResultData struct {
	Values []json.Number `json:"values"`
}
