package protocol

import "encoding/json"

type PollRequest struct {
	InstanceID    string              `json:"instance_id"`
	Name          string              `json:"name"`
	Executors     []string            `json:"executors"`
	Repositories  []string            `json:"repositories"`
	Models        map[string][]string `json:"models,omitempty"`
	FleetRelease  string              `json:"fleet_release,omitempty"`
	ReleaseState  string              `json:"release_state,omitempty"`
	AcceptingWork *bool               `json:"accepting_work,omitempty"`
}

type PollResponse struct {
	Run                  *RunSpec `json:"run,omitempty"`
	RequiredFleetRelease string   `json:"required_fleet_release,omitempty"`
}

type RunSpec struct {
	ID             string `json:"id"`
	JobID          string `json:"job_id"`
	Command        string `json:"command"`
	CommandHash    string `json:"command_hash"`
	Executor       string `json:"executor"`
	Model          string `json:"model,omitempty"`
	Repository     string `json:"repository"`
	RenderedPrompt string `json:"rendered_prompt"`
	TimeoutMillis  int64  `json:"timeout_millis"`
	LeaseToken     string `json:"lease_token"`
}

type Heartbeat struct {
	InstanceID string `json:"instance_id"`
	LeaseToken string `json:"lease_token"`
}

type Completion struct {
	InstanceID string          `json:"instance_id"`
	LeaseToken string          `json:"lease_token"`
	State      string          `json:"state"`
	ExitCode   int             `json:"exit_code"`
	Error      string          `json:"error,omitempty"`
	Result     json.RawMessage `json:"result,omitempty"`
	Events     string          `json:"events,omitempty"`
}
