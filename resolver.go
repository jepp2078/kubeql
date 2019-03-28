package main

import (
	"context"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Resolver struct {
	kubeQLClient *KubeQLClient
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Pods(ctx context.Context) ([]Pod, error) {
	outPod := make([]Pod, 0)
	pods, err := r.kubeQLClient.client.CoreV1().Pods("default").List(metav1.ListOptions{})
	if err != nil {
		glog.V(1).Info(err)
		return nil, err
	}

	for _, pod := range pods.Items {
		objMeta := &ObjectMeta{Name: pod.ObjectMeta.Name}
		newPod := &Pod{ObjectMeta: *objMeta}
		outPod = append(outPod, *newPod)
	}

	return outPod, nil
}
func (r *queryResolver) Pod(ctx context.Context, name string) (*Pod, error) {
	panic("not implemented")
}
