package resources

import (
	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	clientsetnad "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	"golang.org/x/net/context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"sigs.k8s.io/yaml"
	"time"
)

var (

	// inside cluster client
	// creates the in-cluster config
	config, _ = rest.InClusterConfig()

	// creates the clientset
	clientset, _ = kubernetes.NewForConfig(config)
)

func QuotaLimits(nsName string) {

	//========================================================================================
	// apply resources quota, limit range and role binding
	time.Sleep(1 * time.Second)

	// ResourceQuota create
	a, err := os.ReadFile("/files/_resourcequota.yaml")
	if err != nil {
		panic(err)
	}
	// get yaml and convert it to  v1.ResourceQuota
	// to provide it import "k8s.io/apimachinery/pkg/util/yaml"
	quotaData := &v1.ResourceQuota{}
	err = yaml.Unmarshal(a, quotaData)
	if err != nil {
		panic(err)
	}

	// create quota for new namespaces
	_, err = clientset.CoreV1().ResourceQuotas(nsName).Create(context.TODO(), quotaData, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		log.Printf("Failed to create ResourceQuota for %s ", nsName)

	} else {
		log.Printf("Created ResourceQuota for %s", nsName)
	}

	// LimitRange create
	b, err := os.ReadFile("/files/_limitrange.yaml")
	if err != nil {
		panic(err)
	}
	// get yaml and convert it to  v1.ResourceQuota
	// to provide it import "k8s.io/apimachinery/pkg/util/yaml"
	limitRange := &v1.LimitRange{}
	err = yaml.Unmarshal(b, limitRange)
	if err != nil {
		panic(err)
	}

	// create quota for new namespaces
	_, err = clientset.CoreV1().LimitRanges(nsName).Create(context.TODO(), limitRange, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		log.Printf("Failed to create LimitRange for %s ", nsName)

	} else {
		log.Printf("Created LimitRange for %s", nsName)
	}

	// ---------------------------------------
	// create net-attach-def for new namespaces

	// nad client
	cfg, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		log.Printf("Error building kubeconfig: %v", err)
	}

	nadClient, err := clientsetnad.NewForConfig(cfg)
	if err != nil {
		log.Printf("Error building example clientset: %v", err)
	}

	nad := &nadv1.NetworkAttachmentDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "k8s.cni.cncf.io/v1",
			Kind:       "NetworkAttachmentDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "istio-cni",
		},
	}
	_, err = nadClient.K8sCniCncfIoV1().NetworkAttachmentDefinitions(nsName).Create(context.Background(), nad, metav1.CreateOptions{})
	if err != nil {
		log.Println("Cluster is k8s not OpenShift no need to create NetworkAttachmentDefinition ")
		log.Printf("Error %s", err)
	} else {
		log.Println("Created NetworkAttachmentDefinition ")
	}

}
