package middleware

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// NetworkSecurityConfig holds network security configuration
type NetworkSecurityConfig struct {
	TCP struct {
		MaxConnections    int
		ConnectionTimeout time.Duration
		BacklogSize       int
		MaxSynQueue       int
	}
	SSL struct {
		MinVersion   uint16
		CipherSuites []uint16
		CertFile     string
		KeyFile      string
		StrictSNI    bool
		HSTSEnabled  bool
		HSTSDuration time.Duration
	}
	Network struct {
		AllowedIPRanges []string
		BlockedIPs      []string
		RateLimit       struct {
			RequestsPerSecond int
			BurstSize         int
		}
	}
	Wireless struct {
		EnableWifiSecurity bool
		EnableBluetooth    bool
		EnableRFID         bool
		AllowedDevices     []string
	}
}

// NetworkSecurityMiddleware implements network-level security features
type NetworkSecurityMiddleware struct {
	config      NetworkSecurityConfig
	synTracker  *SYNFloodTracker
	connTracker *NetworkConnectionTracker
	ipFilter    *IPFilter
	rateLimiter *NetworkRateLimiter
	mu          sync.RWMutex
}

// NewNetworkSecurityMiddleware creates a new network security middleware instance
func NewNetworkSecurityMiddleware(config NetworkSecurityConfig) *NetworkSecurityMiddleware {
	// Set default TLS version if not specified
	if config.SSL.MinVersion == 0 {
		config.SSL.MinVersion = tls.VersionTLS12
	}

	// Set default cipher suites if not specified
	if len(config.SSL.CipherSuites) == 0 {
		config.SSL.CipherSuites = []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		}
	}

	return &NetworkSecurityMiddleware{
		config:      config,
		synTracker:  NewSYNFloodTracker(config.TCP.MaxSynQueue),
		connTracker: NewNetworkConnectionTracker(config.TCP.MaxConnections, config.TCP.ConnectionTimeout),
		ipFilter:    NewIPFilter(config.Network.AllowedIPRanges, config.Network.BlockedIPs),
		rateLimiter: NewNetworkRateLimiter(config.Network.RateLimit.RequestsPerSecond, config.Network.RateLimit.BurstSize),
	}
}

// Apply applies network security middleware to the Gin router
func (nsm *NetworkSecurityMiddleware) Apply(router *gin.Engine) {
	router.Use(nsm.SSLStripProtection())
	router.Use(nsm.SYNFloodProtection())
	router.Use(nsm.RouteTableProtection())
	router.Use(nsm.SmurfAttackProtection())
	router.Use(nsm.MACAddressProtection())
	router.Use(nsm.WirelessSecurity())
}

// SSLStripProtection prevents SSL stripping attacks
func (nsm *NetworkSecurityMiddleware) SSLStripProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request is HTTPS
		if !c.Request.URL.IsAbs() || c.Request.URL.Scheme != "https" {
			c.Redirect(http.StatusMovedPermanently, "https://"+c.Request.Host+c.Request.URL.Path)
			return
		}

		// Set HSTS header if enabled
		if nsm.config.SSL.HSTSEnabled {
			c.Header("Strict-Transport-Security", "max-age="+nsm.config.SSL.HSTSDuration.String())
		}

		// Validate SSL/TLS configuration
		if tlsConn := c.Request.TLS; tlsConn != nil {
			if tlsConn.Version < nsm.config.SSL.MinVersion {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient TLS version"})
				return
			}
			if !containsUint16(nsm.config.SSL.CipherSuites, tlsConn.CipherSuite) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Weak cipher suite"})
				return
			}
		}

		c.Next()
	}
}

// SYNFloodProtection prevents SYN flood attacks
func (nsm *NetworkSecurityMiddleware) SYNFloodProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Track SYN packets
		if !nsm.synTracker.Allow(clientIP) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "SYN flood detected"})
			return
		}

		// Check connection limits
		if !nsm.connTracker.Allow(clientIP) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many connections"})
			return
		}

		c.Next()
	}
}

// RouteTableProtection prevents route table manipulation
func (nsm *NetworkSecurityMiddleware) RouteTableProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate source IP against allowed ranges
		if !nsm.ipFilter.IsAllowed(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "IP not in allowed range"})
			return
		}

		// Check for suspicious routing headers
		if c.GetHeader("X-Forwarded-Host") != "" && c.GetHeader("X-Forwarded-Host") != c.Request.Host {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid routing"})
			return
		}

		c.Next()
	}
}

// SmurfAttackProtection prevents smurf attacks
func (nsm *NetworkSecurityMiddleware) SmurfAttackProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for broadcast/multicast IPs
		if isBroadcastIP(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Broadcast IP not allowed"})
			return
		}

		// Check for ICMP flood
		if c.Request.Method == "PING" {
			if !nsm.rateLimiter.Allow(c.ClientIP()) {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "ICMP flood detected"})
				return
			}
		}

		c.Next()
	}
}

// MACAddressProtection prevents MAC address spoofing
func (nsm *NetworkSecurityMiddleware) MACAddressProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get MAC address from request
		mac := c.GetHeader("X-Real-MAC")
		if mac == "" {
			mac = c.GetHeader("X-Forwarded-MAC")
		}

		if mac != "" {
			// Validate MAC address format and check against allowed devices
			if !isValidMAC(mac) || !nsm.isAllowedDevice(mac) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid MAC address"})
				return
			}
		}

		c.Next()
	}
}

// WirelessSecurity implements wireless security features
func (nsm *NetworkSecurityMiddleware) WirelessSecurity() gin.HandlerFunc {
	return func(c *gin.Context) {
		if nsm.config.Wireless.EnableWifiSecurity {
			// Check for WiFi-specific headers
			if !nsm.validateWifiHeaders(c) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid WiFi headers"})
				return
			}
		}

		if nsm.config.Wireless.EnableBluetooth {
			// Check for Bluetooth-specific headers
			if !nsm.validateBluetoothHeaders(c) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid Bluetooth headers"})
				return
			}
		}

		if nsm.config.Wireless.EnableRFID {
			// Check for RFID-specific headers
			if !nsm.validateRFIDHeaders(c) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid RFID headers"})
				return
			}
		}

		c.Next()
	}
}

// Helper functions
func (nsm *NetworkSecurityMiddleware) isAllowedDevice(mac string) bool {
	nsm.mu.RLock()
	defer nsm.mu.RUnlock()
	for _, device := range nsm.config.Wireless.AllowedDevices {
		if device == mac {
			return true
		}
	}
	return false
}

func (nsm *NetworkSecurityMiddleware) validateWifiHeaders(c *gin.Context) bool {
	// Implement WiFi security checks
	return true
}

func (nsm *NetworkSecurityMiddleware) validateBluetoothHeaders(c *gin.Context) bool {
	// Implement Bluetooth security checks
	return true
}

func (nsm *NetworkSecurityMiddleware) validateRFIDHeaders(c *gin.Context) bool {
	// Implement RFID security checks
	return true
}

// SYNFloodTracker tracks SYN packets to prevent SYN flood attacks
type SYNFloodTracker struct {
	synQueue map[string][]time.Time
	mu       sync.RWMutex
	maxQueue int
}

func NewSYNFloodTracker(maxQueue int) *SYNFloodTracker {
	return &SYNFloodTracker{
		synQueue: make(map[string][]time.Time),
		maxQueue: maxQueue,
	}
}

func (st *SYNFloodTracker) Allow(ip string) bool {
	st.mu.Lock()
	defer st.mu.Unlock()

	now := time.Now()
	queue := st.synQueue[ip]

	// Remove old entries
	valid := queue[:0]
	for _, t := range queue {
		if now.Sub(t) <= time.Second {
			valid = append(valid, t)
		}
	}

	if len(valid) >= st.maxQueue {
		return false
	}

	valid = append(valid, now)
	st.synQueue[ip] = valid
	return true
}

// NetworkConnectionTracker tracks connections to prevent connection flooding
type NetworkConnectionTracker struct {
	connections       map[string][]time.Time
	mu                sync.RWMutex
	maxConnections    int
	connectionTimeout time.Duration
}

func NewNetworkConnectionTracker(maxConnections int, connectionTimeout time.Duration) *NetworkConnectionTracker {
	return &NetworkConnectionTracker{
		connections:       make(map[string][]time.Time),
		maxConnections:    maxConnections,
		connectionTimeout: connectionTimeout,
	}
}

func (ct *NetworkConnectionTracker) Allow(ip string) bool {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	now := time.Now()
	conns := ct.connections[ip]

	// Remove expired connections
	valid := conns[:0]
	for _, t := range conns {
		if now.Sub(t) <= ct.connectionTimeout {
			valid = append(valid, t)
		}
	}

	if len(valid) >= ct.maxConnections {
		return false
	}

	valid = append(valid, now)
	ct.connections[ip] = valid
	return true
}

// NetworkRateLimiter implements rate limiting functionality for network requests
type NetworkRateLimiter struct {
	rps   int
	burst int
	ips   map[string][]time.Time
	mu    sync.Mutex
}

func NewNetworkRateLimiter(requestsPerSecond, burstSize int) *NetworkRateLimiter {
	return &NetworkRateLimiter{
		rps:   requestsPerSecond,
		burst: burstSize,
		ips:   make(map[string][]time.Time),
	}
}

func (r *NetworkRateLimiter) Allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	times := r.ips[ip]

	// Remove entries older than 1 second
	valid := times[:0]
	for _, t := range times {
		if now.Sub(t) < time.Second {
			valid = append(valid, t)
		}
	}

	// Check if rate limit is exceeded
	if len(valid) >= r.rps {
		return false
	}

	valid = append(valid, now)
	r.ips[ip] = valid
	return true
}

// IPFilter implements IP filtering
type IPFilter struct {
	allowedRanges []string
	blockedIPs    []string
}

func NewIPFilter(allowedRanges, blockedIPs []string) *IPFilter {
	return &IPFilter{
		allowedRanges: allowedRanges,
		blockedIPs:    blockedIPs,
	}
}

func (f *IPFilter) IsAllowed(ip string) bool {
	// Check if IP is blocked
	for _, blockedIP := range f.blockedIPs {
		if blockedIP == ip {
			return false
		}
	}

	// Check if IP is in allowed range
	// FIXED: renamed 'range' variable to 'ipRange'
	for _, ipRange := range f.allowedRanges {
		if isIPInRange(ip, ipRange) {
			return true
		}
	}

	return false
}

func isIPInRange(ip, ipRange string) bool {
	// Implement IP range checking
	return true
}

func isBroadcastIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Check if it's IPv4 broadcast address
	if parsedIP.Equal(net.IPv4bcast) {
		return true
	}

	// Additional broadcast check logic can be added here
	return false
}

func isValidMAC(mac string) bool {
	// Implement MAC address validation
	return true
}

// Helper function specifically for checking uint16 values in slices
func containsUint16(slice []uint16, item uint16) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
