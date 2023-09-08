package lib

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PayloadStateTestCase_ToString struct {
	Input    PayloadState
	Expected string
}

type PayloadStateTestCase_FromString struct {
	Input    string
	Expected PayloadState
}

var toStringTestcases = []PayloadStateTestCase_ToString{
	{
		Input:    Unknown,
		Expected: "unknown",
	},
	{
		Input:    Pending,
		Expected: "pending",
	},
	{
		Input:    Failed,
		Expected: "failed",
	},
	{
		Input:    Complete,
		Expected: "complete",
	},
	{
		Input:    -1,
		Expected: "unknown",
	},
	{
		Input:    4,
		Expected: "unknown",
	},
}

var fromStringTestcases = []PayloadStateTestCase_FromString{
	{
		Input:    "unknown",
		Expected: Unknown,
	},
	{
		Input:    "pending",
		Expected: Pending,
	},
	{
		Input:    "failed",
		Expected: Failed,
	},
	{
		Input:    "complete",
		Expected: Complete,
	},
	{
		Input:    "",
		Expected: Unknown,
	},
	{
		Input:    "foobar",
		Expected: Unknown,
	},
}

func TestPayloadState_ToString(t *testing.T) {
	for _, testcase := range toStringTestcases {
		t.Run(fmt.Sprintf("%d=>%s", testcase.Input, testcase.Expected), func(t *testing.T) {
			actual := PayloadState(testcase.Input).String()
			assert.Equal(t,
				testcase.Expected,
				actual,
			)
		})
	}
}

func TestPayloadState_FromString(t *testing.T) {
	for _, testcase := range fromStringTestcases {
		uppercaseInput := strings.ToUpper(testcase.Input)
		t.Run(fmt.Sprintf("%s=>%s", uppercaseInput, testcase.Expected), func(t *testing.T) {
			actual := ToPayloadState(uppercaseInput)
			assert.Equal(t,
				testcase.Expected,
				actual,
			)
		})
		lowercaseInput := strings.ToLower(testcase.Input)
		t.Run(fmt.Sprintf("%s=>%s", lowercaseInput, testcase.Expected), func(t *testing.T) {
			actual := ToPayloadState(lowercaseInput)
			assert.Equal(t,
				testcase.Expected,
				actual,
			)
		})
	}
}

func TestPayloadState_Marshal(t *testing.T) {
	for _, testcase := range toStringTestcases {
		t.Run(fmt.Sprintf("%d=>%s", testcase.Input, testcase.Expected), func(t *testing.T) {
			actualBytes, err := json.Marshal(testcase.Input)
			require.NoError(t, err)
			assert.Equal(t,
				`"`+testcase.Expected+`"`,
				string(actualBytes),
			)
		})
	}
}

func TestPayloadState_Unmarshal(t *testing.T) {
	for _, testcase := range fromStringTestcases {
		uppercaseInput := `"` + strings.ToUpper(testcase.Input) + `"`
		t.Run(fmt.Sprintf("%s=>%s", uppercaseInput, testcase.Expected), func(t *testing.T) {
			actual := PayloadState(0)
			err := json.Unmarshal([]byte(uppercaseInput), &actual)
			require.NoError(t, err)
			assert.Equal(t,
				testcase.Expected,
				actual,
			)
		})
		lowercaseInput := strings.ToLower(testcase.Input)
		t.Run(fmt.Sprintf("%s=>%s", lowercaseInput, testcase.Expected), func(t *testing.T) {
			actual := PayloadState(0)
			err := json.Unmarshal([]byte(lowercaseInput), &actual)
			require.NoError(t, err)
			assert.Equal(t,
				testcase.Expected,
				actual,
			)
		})
	}
}
