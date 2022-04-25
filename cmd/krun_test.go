package cmd

import (
	"context"
	"io/ioutil"
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	id "github.com/iter8-tools/iter8/driver"

	"github.com/iter8-tools/iter8/base"
)

func TestKRun(t *testing.T) {
	tests := []cmdTestCase{
		// k report
		{
			name:   "k run",
			cmd:    "k run -g default --revision 1 --namespace default",
			golden: base.CompletePath("../testdata", "output/krun.txt"),
		},
	}

	// mock the environment
	base.SetupWithMock(t)
	// fake kube cluster
	*kd = *id.NewFakeKubeDriver(settings)
	kd.Revision = 1
	byteArray, _ := ioutil.ReadFile(base.CompletePath("../testdata", "experiment.yaml"))
	kd.Clientset.CoreV1().Secrets("default").Create(context.TODO(), &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-1-spec",
			Namespace: "default",
		},
		StringData: map[string]string{"experiment.yaml": string(byteArray)},
	}, metav1.CreateOptions{})

	kd.Clientset.BatchV1().Jobs("default").Create(context.TODO(), &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-1-job",
			Namespace: "default",
		},
	}, metav1.CreateOptions{})

	runTestActionCmd(t, tests)
}