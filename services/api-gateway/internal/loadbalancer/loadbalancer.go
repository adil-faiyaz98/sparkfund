package loadbalancer

import (
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

// Strategy defines the load balancing strategy
type Strategy string

const (
	// RoundRobin distributes requests in a circular order
	RoundRobin Strategy = "round_robin"
	// LeastConnections routes to the server with the least active connections
	LeastConnections Strategy = "least_connections"
	// IPHash consistently routes requests from the same client IP to the same server
	IPHash Strategy = "ip_hash"
	// WeightedRoundRobin distributes requests based on server weights
	WeightedRoundRobin Strategy = "weighted_round_robin"
	// Random randomly selects a server
	Random Strategy = "random"
	// LeastResponseTime routes to the server with the lowest response time
	LeastResponseTime Strategy = "least_response_time"
)

// Config holds load balancer configuration
type Config struct {
	Strategy            Strategy      `mapstructure:"strategy"`
	HealthCheckEnabled  bool          `mapstructure:"health_check_enabled"`
	HealthCheckPath     string        `mapstructure:"health_check_path"`
	HealthCheckInterval time.Duration `mapstructure:"health_check_interval"`
	HealthCheckTimeout  time.Duration `mapstructure:"health_check_timeout"`
	RetryCount          int           `mapstructure:"retry_count"`
	RetryWaitTime       time.Duration `mapstructure:"retry_wait_time"`
	MaxIdleConns        int           `mapstructure:"max_idle_conns"`
	MaxConnsPerHost     int           `mapstructure:"max_conns_per_host"`
	IdleConnTimeout     time.Duration `mapstructure:"idle_conn_timeout"`
}

// Server represents a backend server
type Server struct {
	URL             *url.URL
	Weight          int
	Active          bool
	CurrentLoad     int64
	ResponseTime    int64 // in milliseconds
	LastHealthCheck time.Time
	FailCount       int
}

// LoadBalancer manages load balancing between servers
type LoadBalancer struct {
	servers  []*Server
	strategy Strategy
	current  uint64
	mutex    sync.RWMutex
	client   *http.Client
	logger   *logrus.Logger
	config   Config
	rand     *rand.Rand
	stopCh   chan struct{}
}

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer(config Config, logger *logrus.Logger) *LoadBalancer {
	// Create HTTP client with custom transport
	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConns,
		MaxConnsPerHost:     config.MaxConnsPerHost,
		IdleConnTimeout:     config.IdleConnTimeout,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	// Initialize random number generator with seed
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	lb := &LoadBalancer{
		servers:  make([]*Server, 0),
		strategy: config.Strategy,
		client:   client,
		logger:   logger,
		config:   config,
		rand:     rng,
		stopCh:   make(chan struct{}),
	}

	// Start health check if enabled
	if config.HealthCheckEnabled {
		go lb.healthCheck()
	}

	return lb
}

// AddServer adds a server to the load balancer
func (lb *LoadBalancer) AddServer(serverURL string, weight int) error {
	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		return err
	}

	server := &Server{
		URL:             parsedURL,
		Weight:          weight,
		Active:          true,
		CurrentLoad:     0,
		ResponseTime:    0,
		LastHealthCheck: time.Now(),
		FailCount:       0,
	}

	lb.mutex.Lock()
	lb.servers = append(lb.servers, server)
	lb.mutex.Unlock()

	// Perform initial health check
	if lb.config.HealthCheckEnabled {
		go lb.checkServerHealth(server)
	}

	return nil
}

// RemoveServer removes a server from the load balancer
func (lb *LoadBalancer) RemoveServer(serverURL string) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	for i, server := range lb.servers {
		if server.URL.String() == serverURL {
			// Remove server from slice
			lb.servers = append(lb.servers[:i], lb.servers[i+1:]...)
			break
		}
	}
}

// NextServer returns the next server based on the load balancing strategy
func (lb *LoadBalancer) NextServer(clientIP string) *Server {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	// Check if there are any active servers
	activeServers := lb.getActiveServers()
	if len(activeServers) == 0 {
		return nil
	}

	// Select server based on strategy
	switch lb.strategy {
	case RoundRobin:
		return lb.roundRobin(activeServers)
	case LeastConnections:
		return lb.leastConnections(activeServers)
	case IPHash:
		return lb.ipHash(activeServers, clientIP)
	case WeightedRoundRobin:
		return lb.weightedRoundRobin(activeServers)
	case Random:
		return lb.random(activeServers)
	case LeastResponseTime:
		return lb.leastResponseTime(activeServers)
	default:
		// Default to round robin
		return lb.roundRobin(activeServers)
	}
}

// getActiveServers returns all active servers
func (lb *LoadBalancer) getActiveServers() []*Server {
	activeServers := make([]*Server, 0)
	for _, server := range lb.servers {
		if server.Active {
			activeServers = append(activeServers, server)
		}
	}
	return activeServers
}

// roundRobin implements the round robin strategy
func (lb *LoadBalancer) roundRobin(servers []*Server) *Server {
	if len(servers) == 0 {
		return nil
	}

	// Get next server index
	current := atomic.AddUint64(&lb.current, 1) - 1
	index := int(current % uint64(len(servers)))
	return servers[index]
}

// leastConnections implements the least connections strategy
func (lb *LoadBalancer) leastConnections(servers []*Server) *Server {
	if len(servers) == 0 {
		return nil
	}

	// Find server with least connections
	var minServer *Server
	minLoad := int64(^uint64(0) >> 1) // Max int64 value

	for _, server := range servers {
		load := atomic.LoadInt64(&server.CurrentLoad)
		if load < minLoad {
			minLoad = load
			minServer = server
		}
	}

	return minServer
}

// ipHash implements the IP hash strategy
func (lb *LoadBalancer) ipHash(servers []*Server, clientIP string) *Server {
	if len(servers) == 0 {
		return nil
	}

	// Hash client IP to determine server
	hash := 0
	for i := 0; i < len(clientIP); i++ {
		hash = 31*hash + int(clientIP[i])
	}
	index := hash % len(servers)
	if index < 0 {
		index = -index
	}
	return servers[index]
}

// weightedRoundRobin implements the weighted round robin strategy
func (lb *LoadBalancer) weightedRoundRobin(servers []*Server) *Server {
	if len(servers) == 0 {
		return nil
	}

	// Calculate total weight
	totalWeight := 0
	for _, server := range servers {
		totalWeight += server.Weight
	}

	if totalWeight == 0 {
		// If all weights are 0, use simple round robin
		return lb.roundRobin(servers)
	}

	// Get next server based on weight
	current := atomic.AddUint64(&lb.current, 1) - 1
	pos := int(current % uint64(totalWeight))

	runningSum := 0
	for _, server := range servers {
		runningSum += server.Weight
		if pos < runningSum {
			return server
		}
	}

	// Fallback to first server
	return servers[0]
}

// random implements the random strategy
func (lb *LoadBalancer) random(servers []*Server) *Server {
	if len(servers) == 0 {
		return nil
	}

	// Select random server
	lb.mutex.Lock()
	index := lb.rand.Intn(len(servers))
	lb.mutex.Unlock()
	return servers[index]
}

// leastResponseTime implements the least response time strategy
func (lb *LoadBalancer) leastResponseTime(servers []*Server) *Server {
	if len(servers) == 0 {
		return nil
	}

	// Find server with lowest response time
	var minServer *Server
	minTime := int64(^uint64(0) >> 1) // Max int64 value

	for _, server := range servers {
		responseTime := atomic.LoadInt64(&server.ResponseTime)
		if responseTime < minTime {
			minTime = responseTime
			minServer = server
		}
	}

	return minServer
}

// IncrementLoad increments the load counter for a server
func (lb *LoadBalancer) IncrementLoad(server *Server) {
	if server != nil {
		atomic.AddInt64(&server.CurrentLoad, 1)
	}
}

// DecrementLoad decrements the load counter for a server
func (lb *LoadBalancer) DecrementLoad(server *Server) {
	if server != nil {
		atomic.AddInt64(&server.CurrentLoad, -1)
	}
}

// UpdateResponseTime updates the response time for a server
func (lb *LoadBalancer) UpdateResponseTime(server *Server, responseTime time.Duration) {
	if server != nil {
		// Convert to milliseconds
		ms := responseTime.Milliseconds()
		// Use exponential moving average
		current := atomic.LoadInt64(&server.ResponseTime)
		if current == 0 {
			atomic.StoreInt64(&server.ResponseTime, ms)
		} else {
			// Weight: 0.8 current, 0.2 new
			newTime := (current * 8 / 10) + (ms * 2 / 10)
			atomic.StoreInt64(&server.ResponseTime, newTime)
		}
	}
}

// healthCheck periodically checks the health of all servers
func (lb *LoadBalancer) healthCheck() {
	ticker := time.NewTicker(lb.config.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-lb.stopCh:
			return
		case <-ticker.C:
			lb.mutex.RLock()
			servers := make([]*Server, len(lb.servers))
			copy(servers, lb.servers)
			lb.mutex.RUnlock()

			for _, server := range servers {
				go lb.checkServerHealth(server)
			}
		}
	}
}

// checkServerHealth checks the health of a server
func (lb *LoadBalancer) checkServerHealth(server *Server) {
	healthURL := server.URL.String()
	if lb.config.HealthCheckPath != "" {
		healthURL = server.URL.String() + lb.config.HealthCheckPath
	}

	// Create request
	req, err := http.NewRequest("GET", healthURL, nil)
	if err != nil {
		lb.logger.Errorf("Failed to create health check request for %s: %v", server.URL.String(), err)
		return
	}

	// Set timeout
	ctx, cancel := context.WithTimeout(context.Background(), lb.config.HealthCheckTimeout)
	defer cancel()
	req = req.WithContext(ctx)

	// Send request
	start := time.Now()
	resp, err := lb.client.Do(req)
	responseTime := time.Since(start)

	// Update response time
	lb.UpdateResponseTime(server, responseTime)

	// Update server status
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	server.LastHealthCheck = time.Now()

	if err != nil {
		server.FailCount++
		lb.logger.Warnf("Health check failed for %s: %v", server.URL.String(), err)
		if server.FailCount >= 3 {
			if server.Active {
				lb.logger.Warnf("Marking server %s as inactive after %d consecutive failures", server.URL.String(), server.FailCount)
				server.Active = false
			}
		}
		return
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Health check successful
		if !server.Active {
			lb.logger.Infof("Marking server %s as active after successful health check", server.URL.String())
		}
		server.Active = true
		server.FailCount = 0
	} else {
		// Health check failed
		server.FailCount++
		lb.logger.Warnf("Health check failed for %s with status code %d", server.URL.String(), resp.StatusCode)
		if server.FailCount >= 3 {
			if server.Active {
				lb.logger.Warnf("Marking server %s as inactive after %d consecutive failures", server.URL.String(), server.FailCount)
				server.Active = false
			}
		}
	}
}

// Stop stops the load balancer
func (lb *LoadBalancer) Stop() {
	close(lb.stopCh)
}

// GetServers returns all servers
func (lb *LoadBalancer) GetServers() []*Server {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	servers := make([]*Server, len(lb.servers))
	copy(servers, lb.servers)
	return servers
}

// GetStrategy returns the current load balancing strategy
func (lb *LoadBalancer) GetStrategy() Strategy {
	return lb.strategy
}

// SetStrategy sets the load balancing strategy
func (lb *LoadBalancer) SetStrategy(strategy Strategy) {
	lb.strategy = strategy
}
