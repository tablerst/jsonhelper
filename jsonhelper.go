package jsonhelper

import "github.com/tablerst/jsonhelper/pkg/jsonutil"

func Parse(jsonStr string) (interface{}, error) {
	result, err := jsonutil.Parse(jsonStr)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Encode(parseData interface{}, pretty bool) (string, error) {
	output, err := jsonutil.Encode(parseData, pretty)
	if err != nil {
		return "", err
	}
	return output, nil
}
