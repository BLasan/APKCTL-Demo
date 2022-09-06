package utils

import (
	"github.com/ghodss/yaml"
)

func YamlToJson(yamlData []byte) ([]byte, error) {
	return yaml.YAMLToJSON(yamlData)
}
