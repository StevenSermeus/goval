package networking_test

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/StevenSermeus/goval/src/types"
)

type MockConn struct {
	ReadBuffer  bytes.Buffer
	WriteBuffer bytes.Buffer
}

func (m *MockConn) Read(b []byte) (n int, err error) {
	return m.ReadBuffer.Read(b)
}

func (m *MockConn) Write(b []byte) (n int, err error) {
	return m.WriteBuffer.Write(b)
}

func (m *MockConn) Close() error {
	return nil
}

func (m *MockConn) LocalAddr() net.Addr {
	return nil
}

func (m *MockConn) RemoteAddr() net.Addr {
	return nil
}

func (m *MockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *MockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestClient_Receive(t *testing.T) {
	tests := []struct {
		name           string
		mockReadBuffer string
		expected       interface{}
		expectError    bool
	}{
		{
			name:           "Valid string message",
			mockReadBuffer: "+Hello, World!\n\r\n\r",
			expected:       types.CommandInfo{Command: "Hello, World!", ValueType: "string"},
			expectError:    false,
		},
		{
			name:           "Valid integer message",
			mockReadBuffer: ":1234\n\r\n\r",
			expected:       types.CommandInfo{Command: "1234", ValueType: "int"},
			expectError:    false,
		},
		{
			name:           "Invalid message format",
			mockReadBuffer: "Invalid message\n\r\n\r",
			expected:       types.CommandInfo{},
			expectError:    true,
		},
		{
			name:           "Empty message",
			mockReadBuffer: "\n\r\n\r",
			expected:       types.CommandInfo{},
			expectError:    true,
		},
		{
			name:           "End of connection",
			mockReadBuffer: "",
			expected:       "",
			expectError:    true,
		},
		{
			name:           "Valid error message",
			mockReadBuffer: "-Error message\n\r\n\r",
			expected:       types.CommandInfo{Command: "Error message", ValueType: "error"},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConn := &MockConn{}
			mockConn.ReadBuffer.WriteString(tt.mockReadBuffer)

			client := &types.Client{
				Conn: mockConn,
			}

			result, err := client.Receive()
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
			if result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}
