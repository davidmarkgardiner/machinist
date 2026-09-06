# Configuration

Commands use an executor, optional prompt template, and timeout:

```toml
[commands.audit]
executor = "codex"
prompt_file = "prompts/audit.md"
timeout = "30m"

[commands.custom-workflow]
executor = "custom-workflow-script"
timeout = "2h"
```

Without `prompt_file`, the input prompt is sent unchanged. With a template, include
`{{machinist.prompt}}`. Executors and repositories remain worker-owned:

```toml
[executors.custom-workflow-script]
command = ["./scripts/custom-workflow.sh"]

[repositories.my-project]
path = "/absolute/path/to/my-project"
```

Managed triggers select one command with `command = "audit"`. Model selection remains
available when the executor command includes `{{machinist.model}}`.

## Governed fleet releases

A control plane can require workers to advertise a compatible immutable fleet release:

```toml
[server]
fleet_release_policy_file = "~/.machinist/server/fleet-release-policy.json"
```

The policy is deliberately outside the fleet repository, so rollout automation can move
through a dual-release canary window without changing the release being deployed:

```json
{
  "required": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
  "accepted": [
    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
    "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
  ]
}
```

Workers report the release through two machine-local files:

```toml
fleet_release_file = "~/.machinist/current-release/release"
drain_file = "~/.machinist/drain"
```

The worker reads its release marker once at process startup, so changing the marker or
switching the release symlink does not claim that already-running code has changed.
Restart the worker to advertise the new release. The drain marker remains dynamic.

When a release policy is active, a missing, unreadable, invalid, or unaccepted release
remains visible in worker status but receives no new lease. A present drain file has the
same admission effect. An existing valid lease continues to be returned to its worker so
an update never abandons work already in progress. Rollout tools should first create the
drain file and wait for status to show `accepting_work: false` with no active run; they can
then atomically switch `current-release`, restart the worker, and remove the drain file.
During a canary, keep both old and new releases in `accepted`; after every worker reports
the required release, narrow `accepted` to that release alone.

To disable admission gating while retaining the configured policy path, use
`{"required":"","accepted":[]}` and restart the control plane. Malformed policy JSON or
an accepted list that does not contain its required release fails startup rather than
silently admitting an unknown worker release.

## Migration

The `agents` table was renamed to `commands`. Move `[agents.NAME]` to `[commands.NAME]`
and replace `--agent` with `--command`.

The pipeline feature was removed. Replace a sequential pipeline with one executable script,
configure that script as an approved worker executor, and expose it through one command.
Legacy `[pipelines]` configuration fails with migration guidance. Pre-command databases are
recreated once because this release intentionally consolidates the schema before active use.
