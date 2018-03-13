package constants

import (
	"fmt"

	"gopkg.in/resty.v1"
)

const (
	APIHost = "http://localhost:1323"
)

var (
	Client = resty.
		SetHostURL(APIHost).
		OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			fmt.Printf("Status: %v\n", resp.Status())
			fmt.Printf("Body:   %v\n", resp)
			return nil
		})
)
