package gansible

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type LoadOptions struct {
	UseRawInput bool
}

func LoadOptionsFromEnv() LoadOptions {
	useRaw, _ := strconv.ParseBool(os.Getenv("GANSIBLE_USE_RAW_INPUT"))

	return LoadOptions{
		UseRawInput: useRaw,
	}
}

func LoadOptionsDefault() LoadOptions {
	return LoadOptions{UseRawInput: false}
}

func stripAnsibleFields(input []byte) ([]byte, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(input, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input: %w", err)
	}

	for key := range raw {
		if strings.HasPrefix(key, "__ansible") || strings.HasPrefix(key, "_ansible") {
			delete(raw, key)
		}
	}

	return json.Marshal(raw)
}

func InitModule[T any](opts LoadOptions) (T, Result) {
	args, result := LoadArgs[T](opts)
	filterArgsByOS(&args, &result)
	validateRequiredFields(&args, &result)
	return args, result
}
