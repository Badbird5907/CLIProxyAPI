package auth

import (
	"strings"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/registry"
)

// CodexDaybreakBlueModelID is the upstream alias, not a local model rewrite.
const CodexDaybreakBlueModelID = "gpt-daybreak-blue-latest"

// ApplyCodexDaybreakModels requires an explicit per-account opt-in, even if a
// future shared catalog includes Daybreak in a subscription tier's model list.
func ApplyCodexDaybreakModels(auth *Auth, models []*registry.ModelInfo) []*registry.ModelInfo {
	enabled := false
	if auth != nil && strings.EqualFold(auth.Provider, "codex") && auth.AuthKind() == AuthKindOAuth {
		enabled, _ = auth.Metadata["daybreak_blue"].(bool)
	}
	out := make([]*registry.ModelInfo, 0, len(models)+1)
	var daybreak *registry.ModelInfo
	for _, model := range models {
		if model != nil && strings.EqualFold(model.ID, CodexDaybreakBlueModelID) {
			daybreak = model
			continue
		}
		out = append(out, model)
	}
	if enabled {
		if daybreak == nil {
			// The alias can change underlying models. Do not invent fixed limits or
			// borrow capabilities from an unrelated mainline model.
			daybreak = &registry.ModelInfo{
				ID: CodexDaybreakBlueModelID, Object: "model", OwnedBy: "openai",
				Type: "codex", DisplayName: "Daybreak Blue", UserDefined: true,
			}
		}
		out = append(out, daybreak)
	}
	return out
}
