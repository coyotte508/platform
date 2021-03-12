package test

import (
	"context"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/tidepool-org/platform/permission"
)

type GetUserPermissionsInput struct {
	Context       context.Context
	Request       rest.Request
	RequestUserID string
	TargetUserID  string
}

type GetUserPermissionsOutput struct {
	Permissions permission.Permissions
	Error       error
}

type Client struct {
	GetUserPermissionsInvocations int
	GetUserPermissionsInputs      []GetUserPermissionsInput
	GetUserPermissionsStub        func(ctx context.Context, req *rest.Request, requestUserID string, targetUserID string) (permission.Permissions, error)
	GetUserPermissionsOutputs     []GetUserPermissionsOutput
	GetUserPermissionsOutput      *GetUserPermissionsOutput
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetUserPermissions(ctx context.Context, req *rest.Request, requestUserID string, targetUserID string) (permission.Permissions, error) {
	c.GetUserPermissionsInvocations++
	c.GetUserPermissionsInputs = append(
		c.GetUserPermissionsInputs,
		GetUserPermissionsInput{
			Context:       ctx,
			Request:       *req,
			RequestUserID: requestUserID,
			TargetUserID:  targetUserID,
		},
	)
	if c.GetUserPermissionsStub != nil {
		return c.GetUserPermissionsStub(ctx, req, requestUserID, targetUserID)
	}
	if len(c.GetUserPermissionsOutputs) > 0 {
		output := c.GetUserPermissionsOutputs[0]
		c.GetUserPermissionsOutputs = c.GetUserPermissionsOutputs[1:]
		return output.Permissions, output.Error
	}
	if c.GetUserPermissionsOutput != nil {
		return c.GetUserPermissionsOutput.Permissions, c.GetUserPermissionsOutput.Error
	}
	panic("GetUserPermissions has no output")
}

func (c *Client) AssertOutputsEmpty() {
	if len(c.GetUserPermissionsOutputs) > 0 {
		panic("GetUserPermissionsOutputs is not empty")
	}
}
