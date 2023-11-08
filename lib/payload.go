package lib

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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

// func (ps PayloadState) MarshalBSON() ([]byte, error) {
// 	return bson.Marshal(ps.String())
// }

func (ps *PayloadState) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*ps = ToPayloadState(value)
	return nil
}

// func (ps *PayloadState) UnmarshalBSON(data []byte) error {
// 	var value string
// 	if err := bson.Unmarshal(data, &value); err != nil {
// 		return err
// 	}
// 	*ps = ToPayloadState(value)
// 	return nil
// }

type PayloadItem struct {
	ID      uuid.UUID    `json:"id"`
	Time    time.Time    `json:"time"`
	Message *string      `json:"message,omitempty"`
	State   PayloadState `json:"state"`
}

func (pi PayloadItem) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.D{
		{"id", pi.ID.String()},
	})
}

// func (pi *PayloadItem) MarshalBSON() ([]byte, error) {
// 	return bson.Marshal(map[string]any{
// 		"id": pi.ID.String(),
// 		"time": pi.Time,

// 	})
// }
