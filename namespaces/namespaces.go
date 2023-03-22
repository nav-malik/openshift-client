package namespaces

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	v1 "github.com/openshift/api/project/v1"
	projectv1client "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func GetNamespaces(token string, serverURL string) (namespaces *v1.ProjectList, err error) {

	if len(serverURL) == 0 || serverURL == "" || len(token) == 0 || token == "" {
		return nil, errors.New("token, and serverURL are required")
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

	clientset, err := projectv1client.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	namespaceClient := clientset.Projects()
	namespaces, err = namespaceClient.List(context.Background(), metav1.ListOptions{})
	return namespaces, err
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("namespaces")
	// for _, namespace := range namespaces.Items {
	// 	fmt.Println(namespace.Name)
	// }

}
