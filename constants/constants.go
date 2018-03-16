package constants

import (
	"github.com/gky360/atcli/client"
)

const (
	APIHost = "http://localhost:1323"
)

var (
	Client = client.NewClient(APIHost)
)
