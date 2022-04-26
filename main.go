package main

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"time"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		klog.Fatal("$HOME/.kube/config not  found")
	}

	clientSet, err := kubernetes.NewForConfig(config)

	informerFactory := informers.NewSharedInformerFactory(clientSet, 30*time.Second)
	deployInformer := informerFactory.Apps().V1().Deployments()

	informer := deployInformer.Informer()
	deployLister := deployInformer.Lister()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAddDeploy,
		UpdateFunc: onUpdateDeploy,
		DeleteFunc: onDeleteDeploy,
	})

	stopCh := make(chan struct{})
	defer close(stopCh)

	// srart informer Listen && watch deployments
	informerFactory.Start(stopCh)

	// wait for all informer handlers to be initialized
	informerFactory.WaitForCacheSync(stopCh)

	deloyList, err := deployLister.Deployments("default").List(labels.Everything())

	for index, deloyment := range deloyList {

		fmt.Printf("%d ----> %s\n", index, deloyment.Name)

	}

	<-stopCh

}

func onAddDeploy(obj interface{}) {
	deployment := obj.(*v1.Deployment)
	fmt.Printf("Deployment %s added\n", deployment.Name)
}

func onUpdateDeploy(old, new interface{}) {
	oldDeployment := old.(*v1.Deployment)
	newDeployment := new.(*v1.Deployment)
	klog.Infoln("update deployment ", oldDeployment.Status.Replicas, newDeployment.Status.Replicas)
}

func onDeleteDeploy(obj interface{}) {
	deployment := obj.(*v1.Deployment)
	fmt.Printf("Deployment %s delete\n", deployment.Name)
}
