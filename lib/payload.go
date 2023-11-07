package lib

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type PayloadState int

const (
	Unknown PayloadState = iota
	Pending
	Failed
	Complete
)

func (ps PayloadState) String() string {
	values := []string{"unknown", "pending", "failed", "complete"}
	if ps < 0 || int(ps) >= len(values) {
		return values[0]
	}
	return values[ps]
}

func ToPayloadState(value string) PayloadState {
	return PayloadState(max(slices.Index([]string{"unknown", "pending", "failed", "complete"}, strings.ToLower(value)), 0))
}

func (ps PayloadState) MarshalJSON() ([]byte, error) {
	return json.Marshal(ps.String())
}

func (ps *PayloadState) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*ps = ToPayloadState(value)
	return nil
}

type PayloadItem struct {
	ID       uuid.UUID    `json:"id"`
	Time     time.Time    `json:"time"`
	TaskName string       `json:"task_name"`
	Message  string       `json:"message,omitempty"`
	State    PayloadState `json:"state"`
}
