package main

import (
	"context"
)

func createRecord(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	res := map[string]interface{}{
		"data": data,
	}

	return res, nil
}
