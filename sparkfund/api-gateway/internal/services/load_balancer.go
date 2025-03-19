package services

import (
	"github.com/adil-faiyaz98/sparkfund/api-gateway/internal/errors"
	"github.com/adil-faiyaz98/sparkfund/api-gateway/internal/models"
)

type ServiceInstance interface {
	ForwardRequest(req *models.Request) (*models.Response, error)
	HealthCheck() bool
}

type LoadBalancer interface {
	GetInstance(serviceName string) (ServiceInstance, error)
	AddInstance(serviceName string, instance ServiceInstance)
	RemoveInstance(serviceName string, instance ServiceInstance)
	GetInstances(serviceName string) []ServiceInstance
}

type loadBalancer struct {
	instances map[string][]ServiceInstance
}

func NewLoadBalancer() LoadBalancer {
	return &loadBalancer{
		instances: make(map[string][]ServiceInstance),
	}
}

func (lb *loadBalancer) GetInstance(serviceName string) (ServiceInstance, error) {
	instances := lb.instances[serviceName]
	if len(instances) == 0 {
		return nil, errors.ErrNoInstancesAvailable
	}
	// TODO: Implement load balancing strategy (e.g., round-robin, least connections)
	return instances[0], nil
}

func (lb *loadBalancer) AddInstance(serviceName string, instance ServiceInstance) {
	lb.instances[serviceName] = append(lb.instances[serviceName], instance)
}

func (lb *loadBalancer) RemoveInstance(serviceName string, instance ServiceInstance) {
	instances := lb.instances[serviceName]
	for i, inst := range instances {
		if inst == instance {
			lb.instances[serviceName] = append(instances[:i], instances[i+1:]...)
			break
		}
	}
}

func (lb *loadBalancer) GetInstances(serviceName string) []ServiceInstance {
	return lb.instances[serviceName]
}
