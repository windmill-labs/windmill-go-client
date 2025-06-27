package windmill

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"

	"github.com/google/uuid"
	api "github.com/windmill-labs/windmill-go-client/api"
)

type ClientWithWorkspace struct {
	Client    *api.ClientWithResponses
	Workspace string
}

func NewClient(baseUrl, token, workspace string) (ClientWithWorkspace, error) {
	client, err := api.NewClientWithResponses(baseUrl, func(c *api.Client) error {
		c.RequestEditors = append(c.RequestEditors, func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
			return nil
		})
		return nil
	})
	return ClientWithWorkspace{
		Client:    client,
		Workspace: workspace,
	}, err
}

func (c ClientWithWorkspace) GetVariable(path string) (string, error) {
	res, err := c.Client.GetVariableValueWithResponse(context.Background(), c.Workspace, path)
	if err != nil {
		return "", err
	}
	if res.StatusCode()/100 != 2 {
		return "", errors.New(string(res.Body))
	}
	return *res.JSON200, nil
}

func (c ClientWithWorkspace) GetResource(path string) (interface{}, error) {
	params := api.GetResourceValueInterpolatedParams{}
	res, err := c.Client.GetResourceValueInterpolatedWithResponse(context.Background(), c.Workspace, path, &params)
	if err != nil {
		return nil, err
	}
	if res.StatusCode()/100 != 2 {
		return nil, errors.New(string(res.Body))
	}
	return *res.JSON200, nil
}

func (c ClientWithWorkspace) SetResource(path string, value interface{}, resourceTypeOpt ...string) error {
	params := api.GetResourceValueInterpolatedParams{}
	getRes, getErr := c.Client.GetResourceValueInterpolatedWithResponse(context.Background(), c.Workspace, path, &params)
	if getErr != nil {
		return getErr
	}
	if getRes.StatusCode() == 404 {
		resourceType := "any"
		if len(resourceTypeOpt) > 0 {
			resourceType = resourceTypeOpt[0]
		}
		res, err := c.Client.CreateResourceWithResponse(context.Background(), c.Workspace, &api.CreateResourceParams{
			UpdateIfExists: newBool(true),
		}, api.CreateResource{Value: &value, Path: path, ResourceType: resourceType})
		if err != nil {
			return err
		}
		if res.StatusCode()/100 != 2 {
			return errors.New(string(res.Body))
		}
	} else {
		res, err := c.Client.UpdateResourceValueWithResponse(context.Background(), c.Workspace, path, api.UpdateResourceValueJSONRequestBody{
			Value: &value,
		})
		if err != nil {
			return err
		}
		if res.StatusCode()/100 != 2 {
			return errors.New(string(res.Body))
		}
	}
	return nil
}

func (c ClientWithWorkspace) SetVariable(path string, value string) error {
	f := false
	res, err := c.Client.UpdateVariableWithResponse(context.Background(), c.Workspace, path, &api.UpdateVariableParams{AlreadyEncrypted: &f}, api.EditVariable{Value: &value})
	if err != nil {
		f = true
	}
	if res.StatusCode()/100 != 2 {
		f = true
	}
	if f {
		res, err := c.Client.CreateVariableWithResponse(context.Background(), c.Workspace, &api.CreateVariableParams{},
			api.CreateVariableJSONRequestBody{
				Path:  path,
				Value: value,
			})

		if err != nil {
			return err
		}
		if res.StatusCode()/100 != 2 {
			return errors.New(string(res.Body))
		}
	}
	return nil
}

func (c ClientWithWorkspace) GetResumeUrls(approver string) (ResumeUrls, error) {
	var urls ResumeUrls
	jobId, err := uuid.Parse(os.Getenv("WM_JOB_ID"))
	if err != nil {
		return urls, err
	}
	params := api.GetResumeUrlsParams{Approver: &approver}
	nonce := rand.Intn(int(math.MaxUint32))
	res, err := c.Client.GetResumeUrlsWithResponse(context.Background(),
		c.Workspace,
		jobId,
		nonce,
		&params,
	)
	if err != nil {
		return urls, err
	}
	if res.StatusCode()/100 != 2 {
		return urls, errors.New(string(res.Body))
	}
	urls = *res.JSON200
	return urls, nil
}
