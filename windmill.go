package windmill

import (
	"os"
)

func GetClient() (ClientWithWorkspace, error) {
	base_url := os.Getenv("BASE_INTERNAL_URL")
	workspace := os.Getenv("WM_WORKSPACE")
	token := os.Getenv("WM_TOKEN")

	return NewClient(base_url, token, workspace)
}
func newBool(b bool) *bool {
	return &b
}

func GetVariable(path string) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", err
	}
	return client.GetVariable(path)
}

func GetResource(path string) (interface{}, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	return client.GetResource(path)
}

func SetResource(path string, value interface{}, resourceTypeOpt ...string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}
	return client.SetResource(path, value, resourceTypeOpt...)
}

func SetVariable(path string, value string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}
	return client.SetVariable(path, value)
}

func GetStatePath() string {
	value := os.Getenv("WM_STATE_PATH_NEW")
	if len(value) == 0 {
		return os.Getenv("WM_STATE_PATH")
	}
	return value
}

func GetState() (interface{}, error) {
	return GetResource(GetStatePath())
}

func SetState(state interface{}) error {
	err := SetResource(GetStatePath(), state)
	if err != nil {
		return err
	}
	return nil
}

type ResumeUrls struct {
	ApprovalPage string `json:"approvalPage"`
	Cancel       string `json:"cancel"`
	Resume       string `json:"resume"`
}

func GetResumeUrls(approver string) (ResumeUrls, error) {
	client, err := GetClient()
	if err != nil {
		return ResumeUrls{}, err
	}
	return client.GetResumeUrls(approver)
}
