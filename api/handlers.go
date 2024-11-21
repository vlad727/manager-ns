package api

import (
	"encoding/json"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/utils/strings/slices"
	"log"
	"manager-ns/clientset"
	"net/http"
	"os"
	"strings"
)

// Validate handlers accepts or rejects based on request contents
func Validate(w http.ResponseWriter, r *http.Request) {

	// var arReview with struct v1beta1.AdmissionReview{}
	arReview := v1beta1.AdmissionReview{}

	// decode arReview to json and check request
	if err := json.NewDecoder(r.Body).Decode(&arReview); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if arReview.Request == nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// https://kubernetes.io/docs/reference/config-api/apiserver-admission.v1/
	// about admission request and response
	// get namespace name and put it to var
	nsName := arReview.Request.Namespace
	// get requester name and put it to var
	userInfo := arReview.Request.UserInfo.Username
	log.Printf("Requested namespace is %s", nsName)
	log.Printf("Requester for namespace is %s", userInfo)

	// object struct AdmissionResponse
	arReview.Response = &v1beta1.AdmissionResponse{
		UID:     arReview.Request.UID,
		Allowed: true,
	}
	//log.Println("The end of func validate.bac")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&arReview)

	//========================================================================================
	// is it allowed to create quota and limit range?
	//log.Printf("start to check namespace %s", nsName)
	// get list namespaces from configmap
	exceptionNs, err := os.ReadFile("/files/_namespacelist")
	if err != nil {
		log.Println("config not found...")
		log.Println(err)
	}
	// convert it to string
	nsToStr := string(exceptionNs)

	// convert to slice
	stringToSlice := strings.Split(nsToStr, "\n")
	log.Printf("Checking namespace, is it allowed to create resources for %s?", nsName)
	// check namespace name is forbidden to create resources?
	if slices.Contains(stringToSlice, nsName) {
		log.Println("Allowed False")
		log.Printf("Create resources for namespace %s is forbidden", nsName)
	} else {
		log.Println("Allowed True")
		//log.Println("Start to create resources")
		go clientset.CreateObjects(nsName, userInfo)
	}

}
