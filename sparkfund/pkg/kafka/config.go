package kafka

// Config holds Kafka configuration settings
type Config struct {
	Brokers  []string
	GroupID  string
	Username string
	Password string
}

// NewConfig creates a new Kafka configuration with default values
func NewConfig() *Config {
	return &Config{
		Brokers: []string{"localhost:9092"},
		GroupID: "sparkfund-group",
	}
}

// WithBrokers sets the Kafka brokers
func (c *Config) WithBrokers(brokers []string) *Config {
	c.Brokers = brokers
	return c
}

// WithGroupID sets the consumer group ID
func (c *Config) WithGroupID(groupID string) *Config {
	c.GroupID = groupID
	return c
}

// WithCredentials sets the Kafka credentials
func (c *Config) WithCredentials(username, password string) *Config {
	c.Username = username
	c.Password = password
	return c
}
