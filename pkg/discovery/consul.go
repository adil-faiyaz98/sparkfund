package discovery

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

// ServiceDiscovery defines the interface for service registration and discovery
type ServiceDiscovery interface {
	Register(name, host string, port int, tags []string) error
	Deregister() error
	DiscoverService(service string) (string, error)
}

// ConsulClient implements ServiceDiscovery using Consul
type ConsulClient struct {
	client       *consulapi.Client
	serviceID    string
	registration *consulapi.AgentServiceRegistration
}

// NewConsulClient creates a new ConsulClient
func NewConsulClient(addr string) (*ConsulClient, error) {
	config := consulapi.DefaultConfig()
	config.Address = addr
	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConsulClient{client: client}, nil
}

// Register registers a service with Consul
func (c *ConsulClient) Register(name, host string, port int, tags []string) error {
	c.serviceID = fmt.Sprintf("%s-%s-%d", name, host, port)

	reg := &consulapi.AgentServiceRegistration{
		ID:      c.serviceID,
		Name:    name,
		Tags:    tags,
		Port:    port,
		Address: host,
		Check: &consulapi.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", host, port),
			Interval:                       "10s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	c.registration = reg
	return c.client.Agent().ServiceRegister(reg)
}

// Deregister removes a service from Consul
func (c *ConsulClient) Deregister() error {
	return c.client.Agent().ServiceDeregister(c.serviceID)
}

// DiscoverService finds a service by name
func (c *ConsulClient) DiscoverService(service string) (string, error) {
	services, _, err := c.client.Health().Service(service, "", true, nil)
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("service '%s' not found", service)
	}

	// Basic round-robin load balancing - could be enhanced
	service = services[0]
	return fmt.Sprintf("http://%s:%d", service.Service.Address, service.Service.Port), nil
}
