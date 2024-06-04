package clientset

import (
	"context"
	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	clientsetnad "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"strings"
	"time"
)

var (

	// outside cluster client
	/*
		config, _       = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
		clientset, _    = kubernetes.NewForConfig(config)
	*/

	// inside cluster client
	// creates the in-cluster config
	config, _ = rest.InClusterConfig()

	// creates the clientset
	clientset, _ = kubernetes.NewForConfig(config)
)

func CreateObjects(nsName, userInfo string) {

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

	// ---------------------------------------
	// parse requester user or service account
	log.Printf("Requester %s", userInfo)
	// empty slice
	sliceUser := []string{}
	// append to slice
	sliceUser = append(sliceUser, userInfo)
	parsedUser := strings.Split(userInfo, ":")
	if len(parsedUser) > 1 {

		log.Printf("Parsed user %s", parsedUser[3])
		saUser := parsedUser[3]
		saNs := parsedUser[2]
		//system:serviceaccount:vlku4:vlku4

		roleBinding := rbacv1.RoleBinding{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "rbac.authorization.k8s.io/v1",
				Kind:       "RoleBinding",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: saUser + "-admin-" + nsName,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "admin",
			},
			Subjects: []rbacv1.Subject{{Kind: "ServiceAccount", Name: saUser, Namespace: saNs}},
		}
		// create role binding for service account
		_, err := clientset.RbacV1().RoleBindings(nsName).Create(context.Background(), &roleBinding, metav1.CreateOptions{})
		if err != nil {
			log.Println(err)
			log.Printf("failed to create rolebinding for %s ", saUser)

		} else {
			log.Printf("Created RoleBinding for %s", saUser)
		}

	} else {
		// create role binding for user
		//log.Printf("creating rolebinding for %s", userInfo)
		roleBinding := rbacv1.RoleBinding{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "rbac.authorization.k8s.io/v1",
				Kind:       "RoleBinding",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: userInfo + "-admin-" + nsName,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "admin",
			},
			Subjects: []rbacv1.Subject{{Kind: "User", Name: userInfo}},
		}

		_, err = clientset.RbacV1().RoleBindings(nsName).Create(context.Background(), &roleBinding, metav1.CreateOptions{})
		if err != nil {
			log.Println(err)
			log.Printf("failed to create rolebinding for %s ", userInfo)

		} else {
			log.Printf("created rolebinding for %s", userInfo)
		}

	}

}
