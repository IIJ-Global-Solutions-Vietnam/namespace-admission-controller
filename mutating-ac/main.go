package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gitlab-vn.iij-globalps.jp/iij-k8s-team/namespace-admission-controller/client"
	"gitlab-vn.iij-globalps.jp/iij-k8s-team/namespace-admission-controller/config"
	"gitlab-vn.iij-globalps.jp/iij-k8s-team/namespace-admission-controller/consts"

	"github.com/sirupsen/logrus"
	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwhlog "github.com/slok/kubewebhook/v2/pkg/log"
	kwhlogrus "github.com/slok/kubewebhook/v2/pkg/log/logrus"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
)

type Mutator struct {
	client *client.RancherClient
}

const webhookID = "projectFieldAnnotate"

func NewMutator (url string, token string, clusterID string) (*Mutator, error) {
	c, err := client.New(url, token, clusterID)
	return &Mutator{
		client: c,
	}, err
}

func (m *Mutator) Mutate(_ context.Context, ar *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhmutating.MutatorResult, error) {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		return &kwhmutating.MutatorResult{}, nil
	}

	if ar.UserInfo.Username == consts.IgnoreUser {
		return &kwhmutating.MutatorResult{}, nil
	}
	if ns.Annotations == nil {
		ns.Annotations = make(map[string]string)
	}

	if ns.Annotations[consts.ProjectField] == "" {
		l, err := m.client.GetProjectList(ns.Name)
		if err != nil {
			return &kwhmutating.MutatorResult{}, nil
		}
		var projectId string
		if len(l.Data) != 0 {
			isMember, err, message := m.client.IsProjectMember(ar.UserInfo.Username, l.Data[0].ID)
			if err != nil {
				kwhlogs.Warningf(message)
				return &kwhmutating.MutatorResult{
					Warnings: []string{message},
				}, err
			}
			if !isMember {
				kwhlogs.Warningf(message)
				return &kwhmutating.MutatorResult{
					Warnings: []string{message},
				}, nil
			}
			projectId = l.Data[0].ID
		} else {
			hasPermission, err, message := m.client.HasCreateProjectPermission(ar.UserInfo.Username, ar.UserInfo.Groups)
			if err != nil {
				kwhlogs.Warningf(message)
				return &kwhmutating.MutatorResult{
					Warnings: []string{message},
				}, err
			}
			if !hasPermission{
				kwhlogs.Warningf(message)
				return &kwhmutating.MutatorResult{
					Warnings: []string{message},
				}, err
			}
			p, err := m.client.CreateProject(ns.Name)
			if err != nil {
				return &kwhmutating.MutatorResult{}, err
			}
			if err = m.client.AddProjectMember(ar.UserInfo.Username, p); err != nil {
				return &kwhmutating.MutatorResult{}, err
			}
			projectId = p.ID
		}
		ns.Annotations[consts.ProjectField] = projectId
	}

	return &kwhmutating.MutatorResult{
		MutatedObject: ns,
	}, nil
}


var kwhlogs kwhlog.Logger

func run() error {
	logrusLogEntry := logrus.NewEntry(logrus.New())
	if config.Config.Debug {
		logrusLogEntry.Logger.SetLevel(logrus.DebugLevel)
	}else {
		logrusLogEntry.Logger.SetLevel(logrus.InfoLevel)
	}
	logger := kwhlogrus.NewLogrus(logrusLogEntry)
	kwhlogs = logger

	m, err := NewMutator(config.Config.RancherURL, config.Config.RancherAPIToken, config.Config.ClusterID)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return err
	}

	mt := kwhmutating.MutatorFunc(m.Mutate)

	mcfg := kwhmutating.WebhookConfig{
		ID:      webhookID,
		Obj:     &corev1.Namespace{},
		Mutator: mt,
		Logger:  logger,
	}
	wh, err := kwhmutating.NewWebhook(mcfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating webhook: %s", err)
		return err
	}

	webhookHandler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{Webhook: wh, Logger: logger})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating webhook handler: %s", err)
		return err
	}
	logger.Infof("Listening on :8080")
	err = http.ListenAndServeTLS(":8080", config.Config.CertFilePath, config.Config.KeyFilePath, webhookHandler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error serving webhook: %s", err)
		return err
	}
	return nil
}

func main() {
	if err := run();err != nil {
		fmt.Fprintf(os.Stderr, "error runnig app: %s", err)
		os.Exit(1)
	}
}
