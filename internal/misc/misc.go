package misc

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"math"
	"os"
	"strings"
)

func CalculatePercentage(base float64, percentage float64) float64 {
	return Round(base * percentage / 100)
}

func Round(f float64) float64 {
	return math.Round(f*100) / 100
}

func LoadEnvs(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error when reading %s file: %v", path, err)
	}
	result := make(map[string]string)
	err = yaml.Unmarshal(b, &result)
	if err != nil {
		return fmt.Errorf("error when Unmarshal %s file: %v", path, err)
	}
	var errors []string
	for k, v := range result {
		err = os.Setenv(k, v)
		if err != nil {
			errors = append(errors, fmt.Errorf("error when setting %s env var: %v", k, err).Error())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("errors when setting env vars: %s", strings.Join(errors, ", "))
	}
	return nil
}
