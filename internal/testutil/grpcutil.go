package testutil

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// TestGRPCServer represents a test gRPC server
type TestGRPCServer struct {
	Server   *grpc.Server
	Listener *bufconn.Listener
	T        *testing.T
}

// NewTestGRPCServer creates a new test gRPC server
func NewTestGRPCServer(t *testing.T) *TestGRPCServer {
	lis := bufconn.Listen(bufSize)
	server := grpc.NewServer()
	return &TestGRPCServer{
		Server:   server,
		Listener: lis,
		T:        t,
	}
}

// Start starts the test gRPC server
func (ts *TestGRPCServer) Start() {
	go func() {
		if err := ts.Server.Serve(ts.Listener); err != nil {
			ts.T.Fatalf("Server exited with error: %v", err)
		}
	}()
}

// Stop stops the test gRPC server
func (ts *TestGRPCServer) Stop() {
	ts.Server.Stop()
	ts.Listener.Close()
}

// GetTestConn creates a test gRPC connection
func (ts *TestGRPCServer) GetTestConn() *grpc.ClientConn {
	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return ts.Listener.Dial()
		}),
		grpc.WithInsecure(),
	)
	require.NoError(ts.T, err)
	return conn
}

// CreateTestContext creates a test context with timeout
func CreateTestContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return ctx
}
