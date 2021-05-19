package main

import (
	"context"
	"fmt"
	"gitlab-vn.iij-globalps.jp/iij-k8s-team/namespace-admission-controller/client"
	"gitlab-vn.iij-globalps.jp/iij-k8s-team/namespace-admission-controller/config"
	"gitlab-vn.iij-globalps.jp/iij-k8s-team/namespace-admission-controller/consts"
	"net/http"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sirupsen/logrus"
	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwhlogrus "github.com/slok/kubewebhook/v2/pkg/log/logrus"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhvalidating "github.com/slok/kubewebhook/v2/pkg/webhook/validating"
)

type Validator struct {
	client *client.RancherClient
}

const webhookID = "projectFieldValidate"

func NewValidator(url string, token string, clusterID string) (*Validator, error) {
	c, err := client.New(url, token, clusterID)
	return &Validator{
		client: c,
	}, err
}

func (v *Validator) Validate(c context.Context, ar *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhvalidating.ValidatorResult, error) {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		return &kwhvalidating.ValidatorResult{
			Valid:   false,
			Message: "object isn't namespace resource.",
		}, nil
	}
	if ar.UserInfo.Username == consts.IgnoreUser {
		return &kwhvalidating.ValidatorResult{
			Valid:   true,
			Message: "ok",
		}, nil
	}

	if ns.Annotations == nil {
		ns.Annotations = make(map[string]string)
	}
	if ns.Annotations[consts.ProjectField] == "" {
		return &kwhvalidating.ValidatorResult{
			Valid:   false,
			Message: fmt.Sprintf("projectID doesn't exist. namespace=%s", ns.Name),
		}, nil
	}
	return &kwhvalidating.ValidatorResult{
		Valid:   true,
		Message: "ok",
	}, nil
}

func run() error {
	logrusLogEntry := logrus.NewEntry(logrus.New())
	if config.Config.Debug {
		logrusLogEntry.Logger.SetLevel(logrus.DebugLevel)
	} else {
		logrusLogEntry.Logger.SetLevel(logrus.InfoLevel)
	}
	logger := kwhlogrus.NewLogrus(logrusLogEntry)

	v, err := NewValidator(config.Config.RancherURL, config.Config.RancherAPIToken, config.Config.ClusterID)

	vt := kwhvalidating.ValidatorFunc(v.Validate)

	vcfg := kwhvalidating.WebhookConfig{
		ID:        webhookID,
		Obj:       &corev1.Namespace{},
		Validator: vt,
		Logger:    logger,
	}
	wh, err := kwhvalidating.NewWebhook(vcfg)
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
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error runnig app: %s", err)
		os.Exit(1)
	}
}
