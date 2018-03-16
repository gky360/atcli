package client

import (
	"fmt"
	"net/http"

	"github.com/gky360/atsrv/models"
	"gopkg.in/resty.v1"
)

type AtcliClient struct {
	client *resty.Client
}

func NewClient(host string) *AtcliClient {
	c := new(AtcliClient)
	c.client = resty.
		SetHostURL(host).
		OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			fmt.Printf("Status: %v\n", resp.Status())
			fmt.Printf("Body:   %v\n", resp)
			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("atsrv returned an error: %s\n%v", resp.Status(), resp)
			}
			return nil
		})
	return c
}

func (c *AtcliClient) SetAuthToken(token string) {
	c.client.SetAuthToken(token)
}

func (c *AtcliClient) Login(userID string, password string, user *models.User) (*resty.Response, error) {
	return c.client.R().
		SetBody(models.User{ID: userID, Password: password}).
		SetResult(&user).
		Post("/login")
}

func (c *AtcliClient) Logout(user *models.User) (*resty.Response, error) {
	return c.client.R().
		SetResult(&user).
		Post("/logout")
}
