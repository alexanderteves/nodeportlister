package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"

	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getNodePorts() (*bytes.Buffer, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	writer := new(tabwriter.Writer)
	writer.Init(buf, 0, 4, 4, ' ', 0)
	fmt.Fprintln(writer, "NODEPORT(S)\tSERVICE\tNAMESPACE")
	for _, service := range services.Items {
		for _, port := range service.Spec.Ports {
			if port.NodePort > 0 {
				name := service.ObjectMeta.Name
				namespace := service.ObjectMeta.Namespace
				nodeport := strconv.Itoa(int(port.NodePort))
				fmt.Fprintln(writer, nodeport+"\t"+name+"\t"+namespace)
			}
		}
	}
	writer.Flush()
	return buf, nil
}

func main() {
	bindAddress := "127.0.0.1"
	if os.Getenv("BIND_ADDRESS") != "" {
		bindAddress = os.Getenv("BIND_ADDRESS")
	}

	listenPort := "8080"
	if os.Getenv("LISTEN_PORT") != "" {
		listenPort = os.Getenv("LISTEN_PORT")
	}

	httpRouter := http.NewServeMux()

	httpRouter.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	})

	httpRouter.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		buf, err := getNodePorts()
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", 500)
		} else {
			fmt.Fprint(w, buf.String())
		}
	})

	log.Println("Starting to serve on " + bindAddress + ":" + listenPort)
	http.ListenAndServe(bindAddress+":"+listenPort, httpRouter)
}
