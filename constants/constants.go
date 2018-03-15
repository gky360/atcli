package constants

import (
	"fmt"
	"net/http"

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
			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("atsrv returned an error: %s\n%v", resp.Status(), resp)
			}
			return nil
		})
)
