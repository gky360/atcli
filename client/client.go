package client

import (
	"fmt"
	"net/http"
	"path/filepath"

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
		OnBeforeRequest(onBeforeRequest).
		OnAfterResponse(onAfterResponse)
	return c
}

func onBeforeRequest(c *resty.Client, req *resty.Request) error {
	fmt.Printf("API: %s %s\n", req.Method, req.URL)
	return nil
}

func onAfterResponse(c *resty.Client, resp *resty.Response) error {
	fmt.Printf("Status: %v\n", resp.Status())
	fmt.Printf("Body:   %v\n", resp)
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("atsrv returned an error: %s\n%v", resp.Status(), resp)
	}
	return nil
}

func (c *AtcliClient) SetBasicAuth(username, token string) {
	c.client.SetBasicAuth(username, token)
}

func (c *AtcliClient) Me(user *models.User) (*resty.Response, error) {
	endpoint := "/me"
	return c.client.R().
		SetResult(&user).
		Get(endpoint)
}

func (c *AtcliClient) GetContest(contestID string, contest *models.Contest) (*resty.Response, error) {
	if contestID == "" {
		return nil, fmt.Errorf("Contest id is required.")
	}

	endpoint := filepath.Join("/contests", contestID)
	return c.client.R().
		SetResult(&contest).
		Get(endpoint)
}

func (c *AtcliClient) Join(contestID string, contest *models.Contest) (*resty.Response, error) {
	if contestID == "" {
		return nil, fmt.Errorf("Contest id is required.")
	}

	endpoint := filepath.Join("/contests", contestID, "join")
	return c.client.R().
		SetResult(&contest).
		Post(endpoint)
}

func (c *AtcliClient) GetTasks(contestID string, isFull bool) (*resty.Response, []*models.Task, error) {
	if contestID == "" {
		return nil, nil, fmt.Errorf("Contest id is required.")
	}

	endpoint := filepath.Join("/contests", contestID, "tasks")
	rspGetTasks := new(handlers.RspGetTasks)
	req := c.client.R()
	if isFull {
		req.SetQueryParam("full", "true")
	}
	resp, err := req.
		SetResult(&rspGetTasks).
		Get(endpoint)
	return resp, rspGetTasks.Tasks, err
}

func (c *AtcliClient) GetTask(contestID string, taskName string, task *models.Task) (*resty.Response, error) {
	if contestID == "" {
		return nil, fmt.Errorf("Contest id is required.")
	}
	if taskName == "" {
		return nil, fmt.Errorf("Task Name is required.")
	}

	endpoint := filepath.Join("/contests", contestID, "tasks", taskName)
	return c.client.R().
		SetResult(&task).
		Get(endpoint)
}

func (c *AtcliClient) GetSubmissions(contestID string, taskName string) (*resty.Response, []*models.Submission, error) {
	if contestID == "" {
		return nil, nil, fmt.Errorf("Contest id is required.")
	}

	endpoint := filepath.Join("/contests", contestID, "submissions")
	rspGetSubmissions := new(handlers.RspGetSubmissions)
	req := c.client.R()
	if taskName != "" {
		req.SetQueryParam("task_name", taskName)
	}
	resp, err := req.
		SetResult(&rspGetSubmissions).
		Get(endpoint)
	return resp, rspGetSubmissions.Submissions, err
}

func (c *AtcliClient) GetSubmission(contestID string, sbmID int, sbm *models.Submission) (*resty.Response, error) {
	if contestID == "" {
		return nil, fmt.Errorf("Contest id is required.")
	}
	if sbmID == 0 {
		return nil, fmt.Errorf("Submission id is required.")
	}

	endpoint := filepath.Join("/contests", contestID, "submissions", fmt.Sprintf("%d", sbmID))
	return c.client.R().
		SetResult(&sbm).
		Get(endpoint)
}

func (c *AtcliClient) PostSubmission(contestID string, taskName string, sbmSource string, sbm *models.Submission) (*resty.Response, error) {
	if contestID == "" {
		return nil, fmt.Errorf("Contest id is required.")
	}
	if taskName == "" {
		return nil, fmt.Errorf("Task name is required.")
	}
	if sbmSource == "" {
		return nil, fmt.Errorf("Submission source is required.")
	}

	endpoint := filepath.Join("/contests", contestID, "submissions")
	return c.client.R().
		SetQueryParam("task_name", taskName).
		SetBody(&models.Submission{Source: sbmSource}).
		SetResult(&sbm).
		Post(endpoint)
}
