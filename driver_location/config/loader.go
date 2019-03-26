package config

import (
        "net/http"
        "fmt"
        "io/ioutil"
        "encoding/json"
        "github.com/spf13/viper"
)

// LoadConfigurationFromBranch loads configurations into viper.
func LoadConfigurationFromBranch(configServerUrl string, appName string, profile string, branch string) {
        url := fmt.Sprintf("%s/%s/%s/%s", configServerUrl, appName, profile, branch)
        fmt.Printf("Loading config from %s\n", url)
        body, err := fetchConfiguration(url)
        if err != nil {
                panic("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
        }
        parseConfiguration(body)
}

// fetchConfiguration gets config from given url.
func fetchConfiguration(url string) ([]byte, error) {
        resp, err := http.Get(url)
        if err != nil {
                panic("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
        }
        body, err := ioutil.ReadAll(resp.Body)
        return body, err
}

// parseConfiguration reads input, which is received from config server, into viper as key, value pairs.
func parseConfiguration(body []byte) {
        var cloudConfig springCloudConfig
        err := json.Unmarshal(body, &cloudConfig)
        if err != nil {
                panic("Cannot parse configuration, message: " + err.Error())
        }

        for key, value := range cloudConfig.PropertySources[0].Source {
                viper.Set(key, value)
                //fmt.Printf("Loading config property %v => %v\n", key, value)
        }
        if viper.IsSet("urls.driverLocations.nsq.topic") {
                fmt.Printf("Successfully loaded configuration for topic %s\n", viper.GetString("urls.driverLocations.nsq.topic"))
        }
}

type springCloudConfig struct {
        Name            string           `json:"name"`
        Profiles        []string         `json:"profiles"`
        Label           string           `json:"label"`
        Version         string           `json:"version"`
        PropertySources []propertySource `json:"propertySources"`
}

type propertySource struct {
        Name   string                 `json:"name"`
        Source map[string]interface{} `json:"source"`
}
