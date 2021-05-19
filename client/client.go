package client

import (
	"fmt"
	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	. "github.com/rancher/types/client/management/v3"
)

type RancherClient struct {
	client    *Client
	ClusterID string
}

func New(url string, token string, clusterID string) (*RancherClient, error) {
	c, err := NewClient(&clientbase.ClientOpts{
		URL:      url,
		TokenKey: token,
		Insecure: true,
	})
	return &RancherClient{
		client:    c,
		ClusterID: clusterID,
	}, err
}

func (r *RancherClient) GetProjectList(projectName string) (*ProjectCollection, error) {
	return r.client.Project.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"clusterId": r.ClusterID,
			"name":      projectName,
		},
	})
}

func (r *RancherClient) canI(clusterRole *RoleTemplate, resource string, action string) (bool, error) {
	for _, rule := range clusterRole.Rules {
		if (contains(rule.Resources, resource) || contains(rule.Resources, "*")) && (contains(rule.Verbs, "*") || contains(rule.Verbs, action)) {
			return true, nil
		}
	}
	for _, crid := range clusterRole.RoleTemplateIDs {
		cr, err := r.client.RoleTemplate.ByID(crid)
		if err != nil {
			return false, err
		}
		result, err := r.canI(cr, resource, action)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

func contains(s []string, c string) bool {
	for _, v := range s {
		if v == c {
			return true
		}
	}
	return false
}

func (r *RancherClient) IsProjectMember(username string, projectID string) (bool, error, string) {
	prtb, err := r.client.ProjectRoleTemplateBinding.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"userId":    username,
			"projectId": projectID,
		},
	})
	if err != nil {
		return false, err, fmt.Sprintf("Rancher ProjectRoleTemplateBindings List API failed. username=%s, projectID=%s", username, projectID)
	}
	if len(prtb.Data) != 0 {
		return true, nil, fmt.Sprintf("Project already exists, and user is member of the project")
	}
	return false, nil, fmt.Sprintf("Project already exists, and user isn't member of the project. username=%s, projectID=%s", username, projectID)
}

func (r *RancherClient) HasCreateProjectPermission(username string, userGroups []string) (bool, error, string) {
	globalRoleBindings, err := r.client.GlobalRoleBinding.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"userId": username,
		},
	})
	if globalRoleBindings != nil && err == nil {
		for _, globalRoleBinding := range globalRoleBindings.Data {
			globalRole, _ := r.client.GlobalRole.ByID(globalRoleBinding.GlobalRoleID)
			for _, rule := range globalRole.Rules {
				if (contains(rule.Resources, "projects") || contains(rule.Resources, "*")) && (contains(rule.Verbs, "*") || contains(rule.Verbs, "create")) {
					return true, nil, ""
				}
			}
		}
	}
	for _, userGroup := range userGroups {
		globalRoleBindings, err = r.client.GlobalRoleBinding.List(&types.ListOpts{
			Filters: map[string]interface{}{
				"groupPrincipalId": userGroup,
			},
		})
		if globalRoleBindings != nil && err == nil {
			for _, globalRoleBinding := range globalRoleBindings.Data {
				globalRole, _ := r.client.GlobalRole.ByID(globalRoleBinding.GlobalRoleID)
				for _, rule := range globalRole.Rules {
					if (contains(rule.Resources, "projects") || contains(rule.Resources, "*")) && (contains(rule.Verbs, "*") || contains(rule.Verbs, "create")) {
						return true, nil, ""
					}
				}
			}
		}
	}

	clusterRoleBindings, err := r.client.ClusterRoleTemplateBinding.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"userId":    username,
			"clusterId": r.ClusterID,
		},
	})
	if err != nil {
		return false, err, fmt.Sprintf("Rancher ClusterRoleTemplateBindings List API failed. username=%s, clusterID=%s", username, r.ClusterID)
	}
	for _, clusterRoleBinding := range clusterRoleBindings.Data {
		clusterRole, err := r.client.RoleTemplate.ByID(clusterRoleBinding.RoleTemplateID)
		if err != nil {
			return false, err, fmt.Sprintf("Rancher ClusterRole List API failed. id=%s", clusterRoleBinding.ID)
		}
		result, err := r.canI(clusterRole, "projects", "create")
		if err != nil {
			return false, err, fmt.Sprintf("Rancher ClusterRole List API failed. id=%s", clusterRoleBinding.ID)
		}
		if result {
			return true, nil, ""
		}
	}
	for _, userGroup := range userGroups {
		clusterRoleBindings, err = r.client.ClusterRoleTemplateBinding.List(&types.ListOpts{
			Filters: map[string]interface{}{
				"groupPrincipalId": userGroup,
			},
		})
		if clusterRoleBindings != nil && err == nil {
			for _, clusterRoleBinding := range clusterRoleBindings.Data {
				clusterRole, _ := r.client.RoleTemplate.ByID(clusterRoleBinding.RoleTemplateID)
				if err != nil {
					return false, err, fmt.Sprintf("Rancher ClusterRole List API failed. id=%s", clusterRoleBinding.ID)
				}
				result, err := r.canI(clusterRole, "projects", "create")
				if err != nil {
					return false, err, fmt.Sprintf("Rancher ClusterRole List API failed. id=%s", clusterRoleBinding.ID)
				}
				if result {
					return true, nil, ""
				}
			}
		}
	}
	return false, nil, fmt.Sprintf("Specified user doesn't have the permission of creating project resource. username=%s", username)
}

func (r *RancherClient) CreateProject(projectName string) (*Project, error) {
	return r.client.Project.Create(&Project{ClusterID: r.ClusterID, Name: projectName})
}

func (r *RancherClient) AddProjectMember(username string, project *Project) error {
	roleList, err := r.client.RoleTemplate.ListAll(&types.ListOpts{
		Filters: map[string]interface{}{
			"context":               "project",
			"projectCreatorDefault": "true",
		},
	})
	if err != nil {
		return err
	}
	for _, role := range roleList.Data {
		_, err := r.client.ProjectRoleTemplateBinding.Create(
			&ProjectRoleTemplateBinding{
				ProjectID:      project.ID,
				UserID:         username,
				RoleTemplateID: role.ID},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
