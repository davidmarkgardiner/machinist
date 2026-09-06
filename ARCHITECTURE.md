# Architecture

Machinist owns process execution, not orchestration.

- `config.toml` defines portable named commands, optional prompt templates, timeouts,
  triggers, and server settings.
- `worker.toml` defines approved executor argument arrays and logical repository paths.
- `internal/runner` starts one process in one repository, writes the prompt to stdin,
  streams both output channels, records artifacts and token usage, and terminates the
  process tree on timeout or cancellation.
- `internal/controlplane` stores one job and one run, leases it to a capable worker,
  rejects stale completions, and exposes authenticated APIs and the web UI.
- `internal/managedworker` resolves only worker-owned executor and repository names.

An optional fleet-release policy gates new leases without coupling running work to an
updater. Workers advertise an immutable release plus an accepting-work bit on every poll.
The control plane records and displays both values, returns an existing lease before
applying admission gates, and only gives new work to a compatible, accepting worker.

Each job has exactly one run. The database enforces this with a unique `runs.job_id`.
Terminal state comes only from the process result. There is no internal stage model.
