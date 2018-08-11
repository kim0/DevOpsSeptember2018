/*
Copyright 2018 Ahmed Kamal
*/

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var zero int64
	podnames := make([]string, 0)
	podnames = append(podnames, "azure", "loves", "devops")
	rand.Seed(time.Now().Unix()) // initialize the global PRNG
	var kubeconfig *string
	var clientset *kubernetes.Clientset

	_, err := net.LookupHost("kubernetes") // If can resolve, then inside cluster
	if err != nil {
		// Outside cluster
		fmt.Println("Detected outside cluster.")
		if home := os.Getenv("HOME"); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		// create the clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	} else {
		// Inside cluster
		fmt.Println("Detected inside cluster.")
		// creates the in-cluster config
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}

	time.Sleep(2 * time.Second) // Don't fork too fast!
	pods, err := clientset.CoreV1().Pods("default").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Hello world! Found %d pods in the default namespace\n", len(pods.Items))
	if len(pods.Items) < 32 { // Fork new instance
		fmt.Println("Too few virii found .. Let's procreate")

		forkTrials := 0
		for forked := 0; forked < 2 && forkTrials < 10; { // Fork 2 virii, spread the vlove, but don't try too hard
			forkTrials++
			parts := []string{
				podnames[rand.Intn(len(podnames))], podnames[rand.Intn(len(podnames))], podnames[rand.Intn(len(podnames))],
			}
			newpodname := strings.Join(parts, "-")
			fmt.Println("Attempting to create pod with name: ", newpodname)
			_, err = clientset.CoreV1().Pods("default").Create(&corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: newpodname,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  newpodname,
							Image: "kim0/k8svirus:1.0",
						},
					},
					TerminationGracePeriodSeconds: &zero,
					RestartPolicy:                 "Never",
				},
			})

			if err != nil {
				fmt.Printf("Error procreating! Got: %s", err.Error())
				if forkTrials < 10 {
					fmt.Printf(" ... Retrying\n")
				} else {
					fmt.Printf(" ... Giving up\n")
				}
				continue
			}
			forked++
		}
	}

	time.Sleep(60 * time.Second)
	for hostname, err := os.Hostname(); hostname == "azure-loves-devops"; {
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Azure", "❤️", " DevOps!") // Extra space needed!! At least on iterm2 on OSX
		// Chosen host lives forever!
		time.Sleep(10 * time.Second)
	}
	// Everyone else, good bye!
}
