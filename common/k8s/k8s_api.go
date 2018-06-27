// +build !debug

package k8s

import (
	"context"
	"os"
	"strconv"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
)

func GetEndpoints(namespace, service string) ([]*Endpoint, error) {
	client, err := k8s.NewInClusterClient()
	if err != nil {
		return nil, err
	}

	var ips []*Endpoint
	var endpoints corev1.Endpoints
	err = client.Get(context.Background(), namespace, service, &endpoints)
	if err != nil {
		return nil, err
	}

	for _, endpoint := range endpoints.Subsets {
		for _, address := range endpoint.Addresses {
			index := getIndex(*address.Hostname)
			item := NewEndpoint()
			item.IP = *address.Ip
			item.Index = index
			for _, port := range endpoint.Ports {
				item.Ports[*port.Name] = int(*port.Port) + index
			}
			ips = append(ips, item)
		}
	}

	return ips
}

func GetVaildPort(namespace, service string) (map[string]int, error) {
	client, err := k8s.NewInClusterClient()
	if err != nil {
		return nil, err
	}
	var endpoints corev1.Endpoints
	err = client.Get(context.Background(), namespace, service, &endpoints)
	if err != nil {
		return nil, err
	}

	ports := make(map[string]int)

	for _, endpoint := range endpoints.Subsets {
		for _, port := range endpoint.Ports {
			ports[*port.Name] = int(*port.Port) + getIndex(os.Getenv("POD_NAME"))
		}
		break
	}

	return ports
}

func getIndex(service string, name string) int {
	if len(service) >= len(name) {
		return 0
	}
	id, err := strconv.Atoi(name[len(service)+1:])
	if err != nil {
		return 0
	}
	return id
}
