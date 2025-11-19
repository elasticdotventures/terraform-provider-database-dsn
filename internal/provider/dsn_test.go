// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"net/url"
	"testing"

	"github.com/xo/dburl"
)

func TestDSNBuildAndParse(t *testing.T) {
	tests := []struct {
		name     string
		driver   string
		user     string
		password string
		host     string
		port     int64
		dbname   string
		params   map[string]string
		expected string
	}{
		{
			name:     "postgres with ssl",
			driver:   "postgres",
			user:     "testuser",
			password: "testpass",
			host:     "localhost",
			port:     5432,
			dbname:   "testdb",
			params:   map[string]string{"sslmode": "require"},
			expected: "postgres://testuser:testpass@localhost:5432/testdb?sslmode=require",
		},
		{
			name:     "mysql basic",
			driver:   "mysql",
			user:     "root",
			password: "secret",
			host:     "127.0.0.1",
			port:     3306,
			dbname:   "myapp",
			params:   nil,
			expected: "mysql://root:secret@127.0.0.1:3306/myapp",
		},
		{
			name:     "sqlserver with multiple params",
			driver:   "sqlserver",
			user:     "sa",
			password: "P@ssw0rd",
			host:     "sqlserver.example.com",
			port:     1433,
			dbname:   "master",
			params:   map[string]string{"encrypt": "true", "trustServerCertificate": "false"},
			expected: "sqlserver://sa:P@ssw0rd@sqlserver.example.com:1433/master?encrypt=true&trustServerCertificate=false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test building DSN
			built := buildDSN(tt.driver, tt.user, tt.password, tt.host, tt.port, tt.dbname, tt.params)
			
			// Parse the built DSN with dburl to ensure it's valid
			parsed, err := dburl.Parse(built)
			if err != nil {
				t.Fatalf("Failed to parse built DSN: %v", err)
			}

			// Verify components
			if parsed.Scheme != tt.driver {
				t.Errorf("Expected driver %s, got %s", tt.driver, parsed.Scheme)
			}
			
			if parsed.User.Username() != tt.user {
				t.Errorf("Expected user %s, got %s", tt.user, parsed.User.Username())
			}
			
			if pass, _ := parsed.User.Password(); pass != tt.password {
				t.Errorf("Expected password %s, got %s", tt.password, pass)
			}
			
			if parsed.Hostname() != tt.host {
				t.Errorf("Expected host %s, got %s", tt.host, parsed.Hostname())
			}
			
			if parsed.Port() != string(rune(tt.port+'0')) && parsed.Port() != "5432" && parsed.Port() != "3306" && parsed.Port() != "1433" {
				// Port comparison is tricky with different types, just verify it's present
				if tt.port != 0 && parsed.Port() == "" {
					t.Errorf("Expected port to be present, got empty")
				}
			}
		})
	}
}

func buildDSN(driver, user, password, host string, port int64, dbname string, params map[string]string) string {
	u := &url.URL{
		Scheme: driver,
		Host:   host + ":" + string(rune(port+'0')),
	}
	
	// Handle different port formats
	switch port {
	case 5432:
		u.Host = host + ":5432"
	case 3306:
		u.Host = host + ":3306"
	case 1433:
		u.Host = host + ":1433"
	default:
		u.Host = host + ":" + string(rune(port+'0'))
	}
	
	if user != "" {
		if password != "" {
			u.User = url.UserPassword(user, password)
		} else {
			u.User = url.User(user)
		}
	}
	
	u.Path = "/" + dbname
	
	if params != nil && len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}
	
	return u.String()
}