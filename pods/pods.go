package pods

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

func GetPods(token string, serverURL string, namespace string) (pods *v1.PodList, err error) {

	if len(serverURL) == 0 || serverURL == "" || len(namespace) == 0 || namespace == "" || len(token) == 0 || token == "" {
		return nil, errors.New("token, serverURL, and namespace are required")
	}
	config := &rest.Config{
		Host:        serverURL,
		BearerToken: token,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	podClient := corev1client.NewForConfigOrDie(config)
	pods, err = podClient.Pods(namespace).List(context.Background(), metav1.ListOptions{})
	return pods, err
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("PODS")
	// for _, pod := range pods.Items {
	// 	fmt.Println(pod.Name)
	// }
}
