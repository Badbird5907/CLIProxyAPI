# Daybreak Blue account routing

This fork supports the exact upstream model ID `gpt-daybreak-blue-latest`
for explicitly enabled Codex OAuth accounts. It does not infer Daybreak access
from a subscription tier or grant access at OpenAI.

Add this top-level field to the existing OAuth auth JSON for the approved account
in your configured `auth-dir`, preserving all existing credential fields:

```json
"daybreak_blue": true
```

Do not enable this field on accounts without access. The auth-file watcher
reloads changes. Request `gpt-daybreak-blue-latest` normally; only enabled
accounts are eligible. If the enabled account is unavailable, the request fails
rather than falling back to an account without the opt-in. Multiple enabled
accounts use the normal credential selection policy.

Existing per-account `excluded_models`, `model_aliases`, and `prefix` settings
still apply. For example, an account with `"prefix": "security"` also exposes
`security/gpt-daybreak-blue-latest`. The upstream model ID remains
`gpt-daybreak-blue-latest`.

Setting the flag to false or removing it unregisters the model for that account.
Token refresh preserves the flag. Other Codex models and the handling of
`cyber_policy` errors are unchanged. The Codex client catalog advertises
`daybreak_blue` as this alias's access program.

This feature covers locally managed Codex OAuth credentials. CLIProxyAPIHome
dispatch requires corresponding account/model configuration in Home.

## Upstream synchronization

The fork's **Sync upstream** GitHub Actions workflow runs daily at 08:23 UTC and
can be run manually. It merges `router-for-me/CLIProxyAPI`'s `main` into the fork's
`main`, tests all packages, builds the server, and pushes only on success.
Conflicts, test failures, or concurrent branch changes stop the workflow without
overwriting the fork. Resolve those failures manually; no force-push is used.

GitHub Actions must be enabled in the fork. Repository rules can prevent the
workflow token from pushing; in that case the workflow fails visibly. Upstream
sync does not update binaries installed on your machine.
