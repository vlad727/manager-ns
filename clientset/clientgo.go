package clientset

import (
	"context"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/utils/strings/slices"
	"log"
	"manager-ns/annotations"
	"manager-ns/resources"
	"strings"
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

	// ---------------------------------------
	// requester may be kubernetes-admin or system:serviceaccount:vlku4:vlku4 or just a login from ldap someusername
	// parse requester user or service account
	//log.Printf("Requester %s", userInfo)
	// send requester name and namespace name for annotation namespace
	go annotations.SetAnnotation(userInfo, nsName)
	// empty slice
	sliceUser := []string{}
	// append to slice
	sliceUser = append(sliceUser, userInfo)
	parsedUser := strings.Split(userInfo, ":")
	if len(parsedUser) > 1 { // create rolebinding for sa

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
		resources.QuotaLimits(nsName) // create quota and limit for namespace

	} else if slices.Contains(parsedUser, "kubernetes-admin") {
		log.Println("No need to create Rolebinding for kubernetes admin")
		log.Println("No need to create ResourceQuota for kubernetes admin")
		log.Println("No need to create LimitRange for kubernetes admin")

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

		_, err := clientset.RbacV1().RoleBindings(nsName).Create(context.Background(), &roleBinding, metav1.CreateOptions{})
		if err != nil {
			log.Println(err)
			log.Printf("Failed to create rolebinding for %s ", userInfo)

		} else {
			log.Printf("Created rolebinding for %s", userInfo)
		}
		resources.QuotaLimits(nsName)
	}

}
