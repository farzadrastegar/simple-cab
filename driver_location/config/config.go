package config

import (
	"io/ioutil"
	"regexp"

	logger "github.com/sirupsen/logrus"
	"github.com/smallfish/simpleyaml"
)

type Handlers struct {
	yaml *simpleyaml.Yaml
}

func (h *Handlers) ReadYaml(filename string) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Fatalf("ERROR: reading config file failed => %s", filename)
	}
	yaml, err := simpleyaml.NewYaml(source)
	if err != nil {
		logger.Fatalf("ERROR: reading config file failed +> %s", filename)
	}

	h.yaml = yaml
}

func (h *Handlers) getYamlValue(args []interface{}) *simpleyaml.Yaml {
	vYaml := h.yaml

	for i := 0; i < len(args); i++ {
		switch args[i].(type) {
		case int:
			idx, ok := args[i].(int)
			if !ok {
				logger.Fatalf("ERROR: missing parameter in yaml => %#v", args)
			}
			vYaml = vYaml.GetIndex(idx)
		case string:
			attr, ok := args[i].(string)
			if !ok {
				logger.Fatalf("ERROR: missing parameter in yaml => %#v", args)
			}
			vYaml = vYaml.Get(attr)
		default:
			logger.Fatalf("ERROR: missing parameter in yaml => %#v", args)
		}
	}

	return vYaml
}

func (h *Handlers) GetPathInMuxFormat(args ...interface{}) string {
	pathYaml := h.getYamlValue(args)
	pathStr, err := pathYaml.String()
	if err != nil {
		logger.Fatalf("ERROR: path not found in yaml => %#v", args)
	}

	re := regexp.MustCompile("/\\:([a-zA-Z][a-zA-Z0-9]*)")

	return re.ReplaceAllString(pathStr, "/{$1}")
}

func (h *Handlers) GetYamlValueStr(args ...interface{}) string {
	value := h.getYamlValue(args)
	valueStr, err := value.String()
	if err != nil {
		logger.Fatalf("ERROR: missing parameter in yaml => %#v", args)
	}
	return valueStr
}

func (h *Handlers) GetYamlValue(args ...interface{}) *simpleyaml.Yaml {
	return h.getYamlValue(args)
}

func NewConfig() *Handlers {
	return &Handlers{}
}

//Code Example
//func (b *Bus) readParams0() error {
//	//create a config handler
//	configHandler := config.NewConfig(b.Logger)
//
//	//read yaml config
//	configHandler.ReadYaml(driver_location.GetConfigFilename())
//	yaml := configHandler.GetYamlValue("urls", "driverLocations", "nsq")
//
//	//set parameters from yaml
//	msg := "ERROR: reading NSQ parameters failed"
//	nsqdLookupAddress, err := yaml.Get("nsqLookupdAddress").String()
//	if err != nil {
//		return driver_location.Error(msg)
//	}
//	nsqdAddress, err := yaml.Get("nsqdAddress").String()
//	if err != nil {
//		return driver_location.Error(msg)
//	}
//	topic, err := yaml.Get("topic").String()
//	if err != nil {
//		return driver_location.Error(msg)
//	}
//	channel, err := yaml.Get("channel").String()
//	if err != nil {
//		return driver_location.Error(msg)
//	}
//	maxInFlight, err := yaml.Get("maxInFlight").Int()
//	if err != nil {
//		return driver_location.Error(msg)
//	}
//	handlerConcurrency, err := yaml.Get("handlerConcurrency").Int()
//	if err != nil {
//		return driver_location.Error(msg)
//	}
//
//	b.ParamsSet = true
//
//	b.Params = &Params{
//		NSQLookupdAddress:  nsqdLookupAddress,
//		NSQDAddress:        nsqdAddress,
//		Topic:              topic,
//		MaxInFlight:        maxInFlight,
//		Channel:            channel,
//		HandlerConcurrency: handlerConcurrency,
//	}
//	return nil
//}

//logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
//configHandler := config.NewConfig(logger)
//configHandler.ReadYaml(driver_location.GetConfigFilename())
//addr := configHandler.GetYamlValueStr("server", "address")
//port := configHandler.GetYamlValueStr("server", "port")

////how to make a fake github commit event using rabbitmq's /monitor
//curl -H "X-Github-Event: push" -H "Content-Type: application/json" -X POST -d '{"commits": [{"modified": ["driver_location-test.yaml"]}],"name":"what is this?"}' -ki http://localhost:8888/monitor
