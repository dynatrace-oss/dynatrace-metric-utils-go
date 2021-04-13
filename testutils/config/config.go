package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	defaultMetricWaitRetries          = 40
	defaultMetricRetryIntervalSeconds = 10
)

// Config holds all configurations from the dtconfig.json file
type Config struct {
	Endpoint                   string `json:"endpoint"`
	APIToken                   string `json:"api_token"`
	MetricWaitRetries          int    `json:"metric_wait_retries"`
	MetricRetryIntervalSeconds int    `json:"metric_retry_interval_seconds"`
}

// LoadConfig loads the configuration from the dtconfig.json file in the
// current working directory
func LoadConfig() *Config {
	// random nonce ensures the metric is unique so tests do not interfere with each other
	rand.Seed(time.Now().Unix())

	cfg := Config{
		Endpoint:                   os.Getenv("DT_METRICS_QUERY_ENDPOINT"),
		APIToken:                   os.Getenv("DT_API_TOKEN"),
		MetricWaitRetries:          defaultMetricWaitRetries,
		MetricRetryIntervalSeconds: defaultMetricRetryIntervalSeconds,
	}

	if jsonFile, err := os.Open("dtconfig.json"); err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("failed to open dtconfig.json: %s", err)
		}
	} else {
		defer jsonFile.Close()
		if byteValue, err := ioutil.ReadAll(jsonFile); err != nil {
			log.Fatalf("failed to read dtconfig.json: %s", err)
		} else {
			if err = json.Unmarshal(byteValue, &cfg); err != nil {
				log.Fatalf("failed to unmarshal dtconfig.json: %s", err)
			}
		}
	}

	if cfg.Endpoint == "" {
		log.Fatalf("Missing required configuration: endpoint")
	}

	if cfg.APIToken == "" {
		log.Fatalf("Missing required configuration: api_token")
	}

	return &cfg
}
