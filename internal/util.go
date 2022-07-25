package internal

import (
	"io/ioutil"
	"os"
	"strings"
)

func ReadFile(path string) ([]string, error) {
	fileHanle, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	defer fileHanle.Close()

	readBytes, err := ioutil.ReadAll(fileHanle)
	if err != nil {
		return nil, err
	}

	results := strings.Split(string(readBytes), "\n")
	res := make([]string, 0)
	for _, item := range results {
		if item != "" {
			res = append(res, item)
		}
	}
	return res, nil
}
