package supportlink

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/config"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/logger"
)

var endpoint = "https://jxsn.dev/api/v1/support"

const (
	MaxDuration      = 2 * time.Hour
	pollWait         = 25
	maxBackoff       = time.Minute
	activityKept     = 20
	maxTextBytes     = 64 << 10
	maxFileBytes     = 8 << 20
	requestTimeout   = 2 * time.Minute
	snapshotInterval = time.Minute
)

var (
	ErrActive          = errors.New("a support connection is already open")
	ErrInvalidDuration = errors.New("the support connection duration must be between 5 minutes and 2 hours")
	errUnknownSession  = errors.New("the support service does not know this session")
	errRejected        = errors.New("the support service rejected this session")
)

type Activity struct {
	At      time.Time `json:"at"`
	Command string    `json:"command"`
	OK      bool      `json:"ok"`
}

type Status struct {
	Active    bool       `json:"active"`
	Connected bool       `json:"connected"`
	Claimed   bool       `json:"claimed"`
	ClaimCode string     `json:"claimCode,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	Activity  []Activity `json:"activity"`
	Error     string     `json:"error,omitempty"`
}

type session struct {
	id        string
	secret    string
	claimCode string
	expiresAt time.Time
	cancel    context.CancelFunc
	done      chan struct{}
}

type request struct {
	ID      string `json:"id"`
	Command string `json:"command"`
}

type pollResponse struct {
	Claimed  bool      `json:"claimed"`
	Requests []request `json:"requests"`
}

type file struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type response struct {
	RequestID string `json:"requestId"`
	OK        bool   `json:"ok"`
	Error     string `json:"error,omitempty"`
	Text      string `json:"text,omitempty"`
	Data      any    `json:"data,omitempty"`
	File      *file  `json:"file,omitempty"`
}

var (
	mu      sync.Mutex
	current *session
	status  = Status{Activity: []Activity{}}
	client  = &http.Client{Timeout: requestTimeout}
)

func Snapshot() Status {
	mu.Lock()
	defer mu.Unlock()
	snapshot := status
	snapshot.Activity = append([]Activity{}, status.Activity...)
	if snapshot.Claimed {
		snapshot.ClaimCode = ""
	}
	return snapshot
}

func Start(duration time.Duration) (Status, error) {
	if duration < 5*time.Minute || duration > MaxDuration {
		return Snapshot(), ErrInvalidDuration
	}
	id, err := randomHex(16)
	if err != nil {
		return Snapshot(), err
	}
	secret, err := randomHex(32)
	if err != nil {
		return Snapshot(), err
	}
	code, err := newClaimCode()
	if err != nil {
		return Snapshot(), err
	}

	mu.Lock()
	if current != nil {
		mu.Unlock()
		return Snapshot(), ErrActive
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(duration))
	expiresAt, _ := ctx.Deadline()
	current = &session{id: id, secret: secret, claimCode: code, expiresAt: expiresAt, cancel: cancel, done: make(chan struct{})}
	status = Status{Active: true, ClaimCode: code, ExpiresAt: &expiresAt, Activity: []Activity{}}
	active := current
	mu.Unlock()

	logger.Main.Info(fmt.Sprintf("Support connection opened until %s", expiresAt.Format(time.Kitchen)))
	go run(ctx, active)
	return Snapshot(), nil
}

func Stop() Status {
	mu.Lock()
	active := current
	mu.Unlock()
	if active == nil {
		return Snapshot()
	}
	active.cancel()
	<-active.done
	return Snapshot()
}

func run(ctx context.Context, active *session) {
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := active.call(closeCtx, http.MethodDelete, "", nil, nil); err != nil && !errors.Is(err, errUnknownSession) {
			logger.Main.Debug("Could not close the support session: " + err.Error())
		}
		cancel()
		mu.Lock()
		current = nil
		status.Active, status.Connected, status.Claimed, status.ClaimCode, status.ExpiresAt = false, false, false, "", nil
		mu.Unlock()
		logger.Main.Info("Support connection closed")
		close(active.done)
	}()

	backoff := 2 * time.Second
	registered := false
	claimed := false
	var lastSnapshot time.Time
	for ctx.Err() == nil {
		if !registered {
			if err := active.register(ctx); err != nil {
				if errors.Is(err, errRejected) {
					if claimed {
						logger.Main.Info("The support team ended the support connection")
					} else {
						setError(err)
					}
					return
				}
				setConnected(false, err)
				if !sleep(ctx, backoff) {
					return
				}
				backoff = min(backoff*2, maxBackoff)
				continue
			}
			registered = true
		}

		poll, err := active.poll(ctx)
		if ctx.Err() != nil {
			return
		}
		if errors.Is(err, errUnknownSession) {
			registered = false
			continue
		}
		if err != nil {
			setConnected(false, err)
			if !sleep(ctx, backoff) {
				return
			}
			backoff = min(backoff*2, maxBackoff)
			continue
		}
		backoff = 2 * time.Second
		if claimed && !poll.Claimed {
			logger.Main.Info("The support team ended the support connection")
			return
		}
		claimed = poll.Claimed
		mu.Lock()
		status.Connected, status.Claimed, status.Error = true, poll.Claimed, ""
		mu.Unlock()
		for _, req := range poll.Requests {
			active.answer(ctx, req)
		}
		if poll.Claimed && time.Since(lastSnapshot) >= snapshotInterval {
			lastSnapshot = time.Now()
			if err := active.call(ctx, http.MethodPut, "/snapshot", buildSnapshot(), nil); err != nil && ctx.Err() == nil {
				logger.Main.Debug("Could not send the support snapshot: " + err.Error())
			}
		}
	}
}

func (active *session) register(ctx context.Context) error {
	body := map[string]string{
		"secretHash": hashHex(active.secret),
		"claimHash":  hashHex(normalizeClaimCode(active.claimCode)),
		"expiresAt":  active.expiresAt.UTC().Format(time.RFC3339),
		"version":    config.GetVersion(),
	}
	var result struct {
		Claimed bool `json:"claimed"`
	}
	if err := active.call(ctx, http.MethodPut, "", body, &result); err != nil {
		if errors.Is(err, errUnknownSession) {
			return errRejected
		}
		return err
	}
	mu.Lock()
	status.Connected, status.Claimed, status.Error = true, result.Claimed, ""
	mu.Unlock()
	return nil
}

func (active *session) poll(ctx context.Context) (pollResponse, error) {
	var result pollResponse
	err := active.call(ctx, http.MethodGet, fmt.Sprintf("/requests?wait=%d", pollWait), nil, &result)
	return result, err
}

func (active *session) answer(ctx context.Context, req request) {
	reply := execute(req.Command)
	reply.RequestID = req.ID
	if len(reply.Text) > maxTextBytes {
		reply.Text = reply.Text[:maxTextBytes]
	}
	if reply.File != nil && len(reply.File.Data) > maxFileBytes {
		reply = response{RequestID: req.ID, Error: "the result is larger than 8 MB; ask the user to download it from the WebUI"}
	}
	if err := active.call(ctx, http.MethodPost, "/responses", reply, nil); err != nil {
		logger.Main.Warn("Could not send a support response: " + err.Error())
	}

	mu.Lock()
	status.Activity = append(status.Activity, Activity{At: time.Now(), Command: req.Command, OK: reply.OK})
	if len(status.Activity) > activityKept {
		status.Activity = status.Activity[len(status.Activity)-activityKept:]
	}
	mu.Unlock()
	logger.Main.Info(fmt.Sprintf("Support connection ran %q (ok: %v)", req.Command, reply.OK))
}

func (active *session) call(ctx context.Context, method, path string, body, target any) error {
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(encoded)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint+"/sessions/"+active.id+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+active.secret)
	req.Header.Set("User-Agent", "SSUI/"+config.GetVersion())
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusForbidden:
		return errUnknownSession
	case resp.StatusCode < 200 || resp.StatusCode >= 300:
		return fmt.Errorf("support service answered %s", resp.Status)
	case target == nil || resp.StatusCode == http.StatusNoContent:
		return nil
	}
	return json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(target)
}

func setConnected(connected bool, err error) {
	mu.Lock()
	defer mu.Unlock()
	status.Connected = connected
	if err != nil {
		status.Error = err.Error()
	}
}

func setError(err error) {
	mu.Lock()
	status.Error = err.Error()
	mu.Unlock()
	logger.Main.Warn("Support connection failed: " + err.Error())
}

func sleep(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func randomHex(size int) (string, error) {
	buffer := make([]byte, size)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}

func newClaimCode() (string, error) {
	buffer := make([]byte, 17)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	raw := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buffer)[:26]
	groups := []string{"SSUI"}
	for start := 0; start < len(raw); start += 5 {
		groups = append(groups, raw[start:min(start+5, len(raw))])
	}
	return strings.Join(groups, "-"), nil
}

func normalizeClaimCode(code string) string {
	code = strings.ToUpper(strings.TrimSpace(code))
	code = strings.TrimPrefix(code, "SSUI")
	return strings.NewReplacer("-", "", " ", "").Replace(code)
}

func hashHex(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
