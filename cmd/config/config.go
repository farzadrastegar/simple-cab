package config

import (
	"github.com/smallfish/simpleyaml"
	"io/ioutil"
	"log"
	"regexp"
)

type Handlers struct {
	logger *log.Logger
	yaml *simpleyaml.Yaml
}

func (h *Handlers) ReadYaml(filename string) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		h.logger.Fatalf("ERROR: reading config file failed => %s", filename)
	}
	yaml, err := simpleyaml.NewYaml(source)
	if err != nil {
		h.logger.Fatalf("ERROR: reading config file failed +> %s", filename)
	}

	h.yaml = yaml
}

func (h *Handlers) getYamlValue (args []interface{}) *simpleyaml.Yaml{
	vYaml := h.yaml

	for i := 0; i < len(args); i++ {
		switch args[i].(type) {
		case int:
			idx, ok := args[i].(int)
			if !ok {
				h.logger.Fatalf("ERROR: missing parameter in yaml => %#v", args)
			}
			vYaml = vYaml.GetIndex(idx)
		case string:
			attr, ok := args[i].(string)
			if !ok {
				h.logger.Fatalf("ERROR: missing parameter in yaml => %#v", args)
			}
			vYaml = vYaml.Get(attr)
		default:
			h.logger.Fatalf("ERROR: missing parameter in yaml => %#v", args)
		}
	}

	return vYaml
}

func (h *Handlers) GetPathInMuxFormat(args ...interface{}) string {
	pathYaml := h.getYamlValue(args)
	pathStr, err := pathYaml.String()
	if err != nil {
		h.logger.Fatalf("ERROR: path not found in yaml => %#v", args)
	}

	re := regexp.MustCompile("/\\:([a-zA-Z][a-zA-Z0-9]*)")

	return re.ReplaceAllString(pathStr, "/{$1}")
}

func (h *Handlers) GetYamlValueStr(args ...interface{}) string {
	value := h.getYamlValue(args)
	valueStr, err := value.String()
	if err != nil {
		h.logger.Fatalf("ERROR: missing parameter in yaml => %#v", args)
	}
	return valueStr
}

func (h *Handlers) GetYamlValue(args ...interface{}) *simpleyaml.Yaml {
	return h.getYamlValue(args)
}

func NewConfig(logger *log.Logger) *Handlers {
	return &Handlers{
		logger:logger,
	}
}
