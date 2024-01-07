/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package check

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	kubeconfig string
	logger     bool
	namespaces []string
	data       [][]string
)

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check your Kubernetes cluster to see it's current status",
	Run: func(cmd *cobra.Command, args []string) {
		// Log file to store the application output
		if logger {
			logFile, err := os.Create("kubecollect.log")
			if err != nil {
				fmt.Println("Cannot create log file: ", err)
			}
			defer logFile.Close()
			// Write to both logfile and standard output
			multi := io.MultiWriter(logFile, os.Stdout)
			log.SetOutput(multi)
		} else { // Else write to standard output only
			log.SetOutput(io.Writer(os.Stdout))
		}

		namespaces, err := cmd.Flags().GetStringSlice("namespaces")
		if err != nil {
			log.Printf("Error accessing namespaces flag: %v", err)
		}

		// If kubeconfig is not set, then use the default kubeconfig from the home directory
		if kubeconfig == "" {
			if home := homedir.HomeDir(); home != "" {
				kubeconfig = filepath.Join(home, ".kube", "config")
				if _, err := os.Stat(kubeconfig); err == nil {
					log.Print("Using the Kubeconfig in the home directory")
				} else {
					log.Fatal("Kubeconfig file not found. Program will exit.")
				}
			}
		}

		if err := readData(kubeconfig, namespaces); err != nil {
			log.Fatalf("Failed to get the status for your cluster due to: %v", err)
		}
	},
}

// TODO: Fun feat: Display status as Haiku. Doesnt have to include everything
func readData(kubeconfig string, namespaces []string) error {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// If namespaces are not set, have a slice of all the namespaces in the cluster
	if len(namespaces) == 0 {
		ns, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Unable to get namespaces due to %v", err)
		}
		for _, namespace := range ns.Items {
			namespaces = append(namespaces, namespace.ObjectMeta.Name)
		}
	}

	// Loop through the namespaces and print out the relevant details
	for _, namespace := range namespaces {
		// Get podlist
		podList, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}

		for _, pod := range podList.Items {
			// Calculate the age of the pod
			podCreationTime := pod.GetCreationTimestamp()
			age := time.Since(podCreationTime.Time).Round(time.Second)

			// Get the status of each of the pods
			podStatus := pod.Status

			var containerRestarts int32
			var containerReady int
			var totalContainers int
			var containerReasonNotReady string

			// If a pod has multiple containers, get the status from all
			for container := range pod.Spec.Containers {
				// Grab actual reason instead of showing just the podStatus
				if !podStatus.ContainerStatuses[container].Ready {
					if ok := podStatus.ContainerStatuses[container].State.Waiting; ok != nil {
						containerReasonNotReady += podStatus.ContainerStatuses[container].State.Waiting.Reason
					}
					if ok := podStatus.ContainerStatuses[container].State.Terminated; ok != nil {
						containerReasonNotReady += podStatus.ContainerStatuses[container].State.Terminated.Reason
					}
				}

				containerRestarts += podStatus.ContainerStatuses[container].RestartCount
				if podStatus.ContainerStatuses[container].Ready {
					containerReady++
				}
				totalContainers++
			}

			// Get the values from the pod status
			name := pod.GetName()
			ready := fmt.Sprintf("%v/%v", containerReady, totalContainers)
			var actualStatus string
			if len(containerReasonNotReady) > 0 {
				actualStatus = containerReasonNotReady
			} else {
				actualStatus = fmt.Sprintf("%v", podStatus.Phase)
			}

			restarts := fmt.Sprintf("%v", containerRestarts)
			ageS := age.String()

			// Append this to data to be printed in a table
			data = append(data, []string{name, ready, actualStatus, restarts, ageS})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Ready", "Status", "Restarts", "Age"})

	for _, row := range data {
		table.Append(row)
	}

	table.Render() // Send output

	return nil
}

func init() {
	CheckCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "Kubeconfig file location")
	CheckCmd.Flags().BoolVarP(&logger, "logger", "l", false, "Enable logging")
	CheckCmd.Flags().StringSliceP("namespaces", "n", []string{}, "Provide a namespace to get status from")
}
