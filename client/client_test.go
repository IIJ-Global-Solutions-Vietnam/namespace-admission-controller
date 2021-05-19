package client

import (
	"github.com/rancher/norman/types"
	client "github.com/rancher/types/client/management/v3"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestContains(t *testing.T) {
	r := contains([]string{"delete", "create", "update"}, "create")
	if !r {
		t.Error("failed test")
	}
	r = contains([]string{"delete", "update"}, "create")
	if r {
		t.Error("failed test")
	}
}

type ClientMock struct {
	mock.Mock
	Project      ProjectMock
	RoleTemplate mock.Mock
}

type ProjectMock struct {
	mock.Mock
}

func (p *ProjectMock) List(opts types.ListOpts) (*client.ProjectCollection, error) {
	return &client.ProjectCollection{
		Data: []client.Project{},
	}, nil
}
