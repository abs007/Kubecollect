# Kubecollect

Go tool to scan your Kubernetes cluster and provide you with a quick summary

## Installation

You can run `go build -o kcl` to build the binary locally or download from the **Releases** tab

## Features

1. Quick summary of your Kubernetes cluster status and the option to point to specific namespaces
2. Customizable logging

## Requirements

1. Go version 1.21.5 or higher
2. Access to a Kubernetes cluster

## Issues

Open issues if you're having trouble with the tool or you feel that the documentation is lacking

## Flags

Currently there exists a single sub-command: `check` with 3 flags:

1. **--kubeconfig**: This flag stores the path to the kubeconfig file for your cluster. If not specified, it will try to access the file from the default location in `~/.kube/config`

2. **--logger**: This flag enables logging for the application run and creates a `kubecollect.log` based on the path where your run the app from

3. **--namespaces**: This flag stores a string slice of the namespaces that you want to get cluster status from. If no namespaces are provided then, all the namespaces from your cluster are used to get the cluster status update 

## Usage

![kcl command screenshot 1](/images/image1.png)

![kcl command screenshot 2](/images/image2.png)