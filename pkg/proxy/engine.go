package proxy

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/The-17/agentsecrets/pkg/keyring"
)

// CallRequest is the input to the engine — used by both MCP and HTTP paths.
type CallRequest struct {
	TargetURL  string            // full URL e.g. https://api.stripe.com/v1/charges
	Method     string            // GET, POST, PUT, PATCH, DELETE
	Headers    map[string]string // extra headers to forward (non-auth)
	Body       []byte            // raw request body (optional)
	Injections []Injection       // what to inject and where
	AgentID    string            // optional, for audit logging
}

// Injection describes one credential to inject.
type Injection struct {
	Style     string // "bearer", "basic", "header", "query", "body", "form"
	Target    string // header name, query param (depends on style)
	SecretKey string // keyring key name e.g. "STRIPE_SECRET_KEY"
}

// CallResult is the output from the engine.
type CallResult struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

// SecretResolver is a function that retrieves a secret value by key name.
// This allows the engine to be tested with a mock keyring.
type SecretResolver func(key string) (string, error)

// Engine coordinates keyring lookup, injection, forwarding, and auditing.
type Engine struct {
	ProjectID     string
	Audit         *AuditLogger
	Client        *http.Client
	ResolveSecret SecretResolver
}

// NewEngine creates an engine wired to the real keyring for the given project.
func NewEngine(projectID string) (*Engine, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required — run 'agentsecrets project use <name>' first")
	}

	audit, err := NewAuditLogger("")
	if err != nil {
		// Audit logger is non-critical — log to stderr but continue
		audit = nil
	}

	return &Engine{
		ProjectID: projectID,
		Audit:     audit,
		Client: &http.Client{
			Timeout: DefaultTimeout,
		},
		ResolveSecret: func(key string) (string, error) {
			return keyring.GetSecret(projectID, key)
		},
	}, nil
}

// Execute runs the full proxy pipeline: resolve secrets → inject → forward → audit.
func (e *Engine) Execute(req CallRequest) (*CallResult, error) {
	// --- Validate ---
	if req.TargetURL == "" {
		return nil, fmt.Errorf("target URL is required")
	}
	if len(req.Injections) == 0 {
		return nil, fmt.Errorf("at least one injection is required — specify how to authenticate (e.g. bearer, header, query)")
	}

	method := strings.ToUpper(req.Method)
	if method == "" {
		method = "GET"
	}

	// --- Build outbound request ---
	var bodyReader *bytes.Reader
	if len(req.Body) > 0 {
		bodyReader = bytes.NewReader(req.Body)
	} else {
		bodyReader = bytes.NewReader(nil)
	}

	outbound, err := http.NewRequest(method, req.TargetURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// Copy any extra headers
	for k, v := range req.Headers {
		outbound.Header.Set(k, v)
	}

	// --- Resolve secrets and inject ---
	secretKeys := make([]string, 0, len(req.Injections))
	authStyles := make([]string, 0, len(req.Injections))

	for _, inj := range req.Injections {
		cred, err := e.ResolveSecret(inj.SecretKey)
		if err != nil {
			return nil, fmt.Errorf("secret '%s' not found in keychain — use list_secrets to see available keys, or add it with 'agentsecrets secrets set %s=VALUE'", inj.SecretKey, inj.SecretKey)
		}

		if err := Inject(outbound, cred, inj); err != nil {
			return nil, fmt.Errorf("injection failed for %s (%s): %w", inj.SecretKey, inj.Style, err)
		}

		secretKeys = append(secretKeys, inj.SecretKey)
		authStyles = append(authStyles, inj.Style)
	}

	// --- Forward ---
	result, err := Forward(e.Client, outbound)
	if err != nil {
		return nil, err
	}

	// --- Audit ---
	if e.Audit != nil {
		_ = e.Audit.Log(AuditEvent{
			Timestamp:  time.Now().UTC(),
			SecretKeys: secretKeys,
			AgentID:    req.AgentID,
			Method:     method,
			TargetURL:  req.TargetURL,
			AuthStyles: authStyles,
			StatusCode: result.StatusCode,
			DurationMs: result.Duration.Milliseconds(),
		})
	}

	// --- Build response ---
	headers := make(map[string][]string)
	for k, v := range result.Headers {
		headers[k] = v
	}

	return &CallResult{
		StatusCode: result.StatusCode,
		Headers:    headers,
		Body:       result.Body,
	}, nil
}
