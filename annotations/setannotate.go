package annotations

import (
	"encoding/json"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

var (

	// outside cluster client
	/*
		config, _    = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
		clientset, _ = kubernetes.NewForConfig(config)

	*/
	// inside cluster client
	// creates the in-cluster config
	config, _ = rest.InClusterConfig()

	// creates the clientset
	clientset, _ = kubernetes.NewForConfig(config)
)

// struct for json
type Y struct {
	Metadata Annotations `json:"metadata"`
}

type Annotations struct {
	Annotations Requester `json:"annotations"`
}
type Requester struct {
	Requester string `json:"requester"`
}

func SetAnnotaion(reqname, nsName string) {

	setAnnotation := Y{
		Metadata: Annotations{
			Requester{reqname},
		},
	}

	// marshal var setAnnotation to json
	bytes, _ := json.Marshal(setAnnotation)

	// set annotation to namespace
	//Note: that type used MergePatchType (allow add new piece of json to namespace)
	_, err := clientset.CoreV1().Namespaces().Patch(context.TODO(), nsName, types.MergePatchType, bytes, metav1.PatchOptions{})
	if err != nil {

		log.Println(err)
	}
	log.Println("Namespace has been annotated ", string(bytes))
}

//Result:
/*
	apiVersion: v1
	kind: Namespace
	metadata:
	  annotations:
	    requester: admin
*/
//https://stackoverflow.com/questions/69125257/golang-kubernetes-client-patching-an-existing-resource-with-a-label <<< diff merge and json
