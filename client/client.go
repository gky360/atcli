package client

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gky360/atsrv/handlers"
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
	endpoint := "login"
	return c.client.R().
		SetBody(models.User{ID: userID, Password: password}).
		SetResult(&user).
		Post(endpoint)
}

func (c *AtcliClient) Logout(user *models.User) (*resty.Response, error) {
	endpoint := "logout"
	return c.client.R().
		SetResult(&user).
		Post(endpoint)
}

func (c *AtcliClient) Me(user *models.User) (*resty.Response, error) {
	endpoint := "me"
	return c.client.R().
		SetResult(&user).
		Get(endpoint)
}

func (c *AtcliClient) GetContest(contestID string, contest *models.Contest) (*resty.Response, error) {
	endpoint := path.Join("contests", contestID)
	return c.client.R().
		SetResult(&contest).
		Get(endpoint)
}

func (c *AtcliClient) GetTasks(contestID string, tasks []models.Task) (*resty.Response, error) {
	endpoint := path.Join("contests", contestID, "tasks")
	rspGetTasks := new(handlers.RspGetTasks)
	resp, err := c.client.R().
		SetResult(&rspGetTasks).
		Get(endpoint)
	tasks = rspGetTasks.Tasks
	return resp, err
}

func (c *AtcliClient) GetTask(contestID string, taskID string, task *models.Task) (*resty.Response, error) {
	endpoint := path.Join("contests", contestID, "tasks", taskID)
	return c.client.R().
		SetResult(&task).
		Get(endpoint)
}

func (c *AtcliClient) GetSubmissions(contestID string, taskID string, sbms []models.Submission) (*resty.Response, error) {
	endpoint := path.Join("contests", contestID, "tasks", taskID, "submissions")
	rspGetSubmissions := new(handlers.RspGetSubmissions)
	resp, err := c.client.R().
		SetResult(&rspGetSubmissions).
		Get(endpoint)
	sbms = rspGetSubmissions.Submissions
	return resp, err
}

func (c *AtcliClient) GetSubmission(contestID string, taskID string, sbmId int, sbm *models.Submission) (*resty.Response, error) {
	endpoint := path.Join("contests", contestID, "tasks", taskID, "submissions", string(sbmId))
	return c.client.R().
		SetResult(&sbm).
		Get(endpoint)
}

func (c *AtcliClient) PostSubmission(contestID string, taskID string, sbmId int, sbmSource string, sbm *models.Submission) (*resty.Response, error) {
	endpoint := path.Join("contests", contestID, "tasks", taskID, "submissions", string(sbmId))
	return c.client.R().
		SetResult(&sbm).
		Post(endpoint)
}
