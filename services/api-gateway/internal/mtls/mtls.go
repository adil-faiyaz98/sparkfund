package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

// Config holds mTLS configuration
type Config struct {
	Enabled            bool   `mapstructure:"enabled"`
	CertFile           string `mapstructure:"cert_file"`
	KeyFile            string `mapstructure:"key_file"`
	CAFile             string `mapstructure:"ca_file"`
	VerifyClient       bool   `mapstructure:"verify_client"`
	CertRotationPeriod string `mapstructure:"cert_rotation_period"`
}

// Manager manages TLS certificates
type Manager struct {
	config     Config
	certMu     sync.RWMutex
	cert       *tls.Certificate
	certPool   *x509.CertPool
	logger     *logrus.Logger
	watcher    *fsnotify.Watcher
	stopCh     chan struct{}
	rotationCh chan struct{}
}

// NewManager creates a new TLS manager
func NewManager(config Config, logger *logrus.Logger) (*Manager, error) {
	if !config.Enabled {
		return nil, nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	m := &Manager{
		config:     config,
		logger:     logger,
		watcher:    watcher,
		stopCh:     make(chan struct{}),
		rotationCh: make(chan struct{}, 1),
	}

	// Load initial certificates
	if err := m.loadCertificates(); err != nil {
		return nil, err
	}

	// Watch certificate files for changes
	for _, file := range []string{config.CertFile, config.KeyFile, config.CAFile} {
		if file != "" {
			if err := watcher.Add(filepath.Dir(file)); err != nil {
				return nil, fmt.Errorf("failed to watch certificate file %s: %w", file, err)
			}
		}
	}

	// Start certificate rotation
	go m.watchCertificates()
	go m.rotateCertificates()

	return m, nil
}

// GetTLSConfig returns the TLS configuration
func (m *Manager) GetTLSConfig() *tls.Config {
	if !m.config.Enabled {
		return nil
	}

	m.certMu.RLock()
	defer m.certMu.RUnlock()

	clientAuth := tls.NoClientCert
	if m.config.VerifyClient {
		clientAuth = tls.RequireAndVerifyClientCert
	}

	return &tls.Config{
		Certificates: []tls.Certificate{*m.cert},
		ClientAuth:   clientAuth,
		ClientCAs:    m.certPool,
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		PreferServerCipherSuites: true,
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			m.certMu.RLock()
			defer m.certMu.RUnlock()
			return m.cert, nil
		},
		GetClientCertificate: func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			m.certMu.RLock()
			defer m.certMu.RUnlock()
			return m.cert, nil
		},
	}
}

// GetHTTPTransport returns an HTTP transport with mTLS configured
func (m *Manager) GetHTTPTransport() *http.Transport {
	if !m.config.Enabled {
		return &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		}
	}

	return &http.Transport{
		TLSClientConfig:       m.GetTLSConfig(),
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// Close closes the manager
func (m *Manager) Close() error {
	if m.watcher != nil {
		close(m.stopCh)
		return m.watcher.Close()
	}
	return nil
}

// loadCertificates loads certificates from files
func (m *Manager) loadCertificates() error {
	m.certMu.Lock()
	defer m.certMu.Unlock()

	// Load certificate and key
	cert, err := tls.LoadX509KeyPair(m.config.CertFile, m.config.KeyFile)
	if err != nil {
		return fmt.Errorf("failed to load certificate and key: %w", err)
	}
	m.cert = &cert

	// Load CA certificate
	if m.config.CAFile != "" {
		caData, err := ioutil.ReadFile(m.config.CAFile)
		if err != nil {
			return fmt.Errorf("failed to read CA certificate: %w", err)
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caData) {
			return fmt.Errorf("failed to parse CA certificate")
		}
		m.certPool = certPool
	}

	m.logger.Info("Loaded TLS certificates successfully")
	return nil
}

// watchCertificates watches certificate files for changes
func (m *Manager) watchCertificates() {
	for {
		select {
		case <-m.stopCh:
			return
		case event, ok := <-m.watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				// Check if the modified file is one of our certificate files
				filename := filepath.Base(event.Name)
				certFilename := filepath.Base(m.config.CertFile)
				keyFilename := filepath.Base(m.config.KeyFile)
				caFilename := filepath.Base(m.config.CAFile)

				if filename == certFilename || filename == keyFilename || filename == caFilename {
					m.logger.Infof("Certificate file changed: %s", event.Name)
					// Wait a bit for the file to be fully written
					time.Sleep(100 * time.Millisecond)
					if err := m.loadCertificates(); err != nil {
						m.logger.Errorf("Failed to reload certificates: %v", err)
					} else {
						m.logger.Info("Certificates reloaded successfully")
					}
				}
			}
		case err, ok := <-m.watcher.Errors:
			if !ok {
				return
			}
			m.logger.Errorf("Certificate watcher error: %v", err)
		}
	}
}

// rotateCertificates periodically triggers certificate rotation
func (m *Manager) rotateCertificates() {
	if m.config.CertRotationPeriod == "" {
		return
	}

	rotationPeriod, err := time.ParseDuration(m.config.CertRotationPeriod)
	if err != nil {
		m.logger.Errorf("Invalid certificate rotation period: %v", err)
		return
	}

	ticker := time.NewTicker(rotationPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
			m.logger.Info("Certificate rotation period reached, triggering rotation")
			select {
			case m.rotationCh <- struct{}{}:
				// Signal sent successfully
			default:
				// Channel is full, rotation already pending
			}
		}
	}
}

// GetCertificateExpiration returns the expiration time of the current certificate
func (m *Manager) GetCertificateExpiration() (time.Time, error) {
	m.certMu.RLock()
	defer m.certMu.RUnlock()

	if m.cert == nil {
		return time.Time{}, fmt.Errorf("no certificate loaded")
	}

	leaf, err := x509.ParseCertificate(m.cert.Certificate[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return leaf.NotAfter, nil
}

// IsEnabled returns whether mTLS is enabled
func (m *Manager) IsEnabled() bool {
	return m.config.Enabled
}
