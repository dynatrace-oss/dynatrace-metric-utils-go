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

func GetMostRecentValueForMetric(s selector, cfg config.Config) json.Number {
	log.Printf("Getting metrics for %s", s.String())
	for cnt := 1; cnt <= cfg.MetricWaitRetries; cnt++ {
		metrics := GetMetrics(s, cfg.Endpoint, cfg.APIToken)
		if len(metrics) > 0 {
			log.Printf("Got result: %+v", metrics)
			return metrics[len(metrics)-1]
		}
		log.Printf("Attempt %d / %d no results", cnt, cfg.MetricWaitRetries)
		time.Sleep(time.Duration(cfg.MetricRetryIntervalSeconds) * time.Second)
	}

	log.Fatalf("Failed to get metrics for %s", s.String())
	return json.Number("")
}

func GetMetrics(s selector, endpoint, apiToken string) []json.Number {
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
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 400 {
		// metric is simply not yet available
		return nil
	}

	if resp.StatusCode != 200 {
		log.Printf("Could not get metric: %s", resp.Status)
		log.Printf("%s", string(bodyBytes[:]))
		return nil
	}

	responseBody := metricQueryResponse{}
	if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
		log.Fatalf("failed to unmarshal response: %s", err.Error())
	}

	if count, err := responseBody.TotalCount.Int64(); err != nil {
		log.Fatalf(err.Error())
	} else if count == 0 {
		return nil
	}

	return responseBody.Result[0].Data[0].Values
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
