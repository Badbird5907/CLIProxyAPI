package auth

import (
	"context"
	"net/http"
	"testing"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/registry"
	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/executor"
)

func TestCodexDaybreakRequiresExplicitAccountOptIn(t *testing.T) {
	for _, value := range []any{nil, false, "true", true} {
		auth := &Auth{Provider: "codex", Metadata: map[string]any{"auth_kind": "oauth", "daybreak_blue": value}}
		// A future shared catalog must not enable other accounts accidentally.
		models := ApplyCodexDaybreakModels(auth, []*registry.ModelInfo{{ID: "ordinary"}, {ID: CodexDaybreakBlueModelID}})
		want := value == true
		if got := len(models) == 2; got != want {
			t.Fatalf("opt-in %v: registered = %v, want %v", value, got, want)
		}
		_, ok := lookupCodexOAuthModelCapability(auth, CodexDaybreakBlueModelID)
		if ok != want {
			t.Fatalf("opt-in %v: execution capability = %v, want %v", value, ok, want)
		}
	}
}

func TestCodexDaybreakNeverFallsBackToUnapprovedAccount(t *testing.T) {
	m := NewManager(nil, nil, nil)
	executor := &daybreakRecordingExecutor{authFallbackExecutor: &authFallbackExecutor{id: "codex", executeErrors: map[string]error{}}}
	m.RegisterExecutor(executor)
	reg := registry.GetGlobalRegistry()
	for _, id := range []string{"daybreak-unapproved", "daybreak-approved"} {
		auth := &Auth{ID: id, Provider: "codex", Metadata: map[string]any{"auth_kind": "oauth", "daybreak_blue": id == "daybreak-approved"}}
		reg.RegisterClient(id, "codex", ApplyCodexDaybreakModels(auth, []*registry.ModelInfo{{ID: "ordinary"}}))
		t.Cleanup(func() { reg.UnregisterClient(id) })
		if _, err := m.Register(context.Background(), auth); err != nil {
			t.Fatal(err)
		}
	}
	request := cliproxyexecutor.Request{Model: CodexDaybreakBlueModelID}
	if _, err := m.Execute(context.Background(), []string{"codex"}, request, cliproxyexecutor.Options{}); err != nil {
		t.Fatal(err)
	}
	executor.executeErrors["daybreak-approved"] = &Error{HTTPStatus: http.StatusTooManyRequests, Message: "quota exhausted"}
	_, err := m.Execute(context.Background(), []string{"codex"}, request, cliproxyexecutor.Options{})
	if err == nil {
		t.Fatal("expected approved account quota failure")
	}
	if calls := executor.ExecuteCalls(); len(calls) != 2 || calls[0] != "daybreak-approved" || calls[1] != "daybreak-approved" {
		t.Fatalf("credential attempts = %v, want only approved account", calls)
	}
	for _, model := range executor.models {
		if model != CodexDaybreakBlueModelID {
			t.Fatalf("upstream model rewritten to %q", model)
		}
	}
}

type daybreakRecordingExecutor struct {
	*authFallbackExecutor
	models []string
}

func (e *daybreakRecordingExecutor) Execute(ctx context.Context, auth *Auth, req cliproxyexecutor.Request, opts cliproxyexecutor.Options) (cliproxyexecutor.Response, error) {
	e.models = append(e.models, req.Model)
	return e.authFallbackExecutor.Execute(ctx, auth, req, opts)
}
