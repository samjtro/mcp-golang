package mcp_golang

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type LoggingLevel string

const LoggingLevelAlert LoggingLevel = "alert"
const LoggingLevelCritical LoggingLevel = "critical"
const LoggingLevelDebug LoggingLevel = "debug"
const LoggingLevelEmergency LoggingLevel = "emergency"
const LoggingLevelError LoggingLevel = "error"
const LoggingLevelInfo LoggingLevel = "info"
const LoggingLevelNotice LoggingLevel = "notice"
const LoggingLevelWarning LoggingLevel = "warning"

var enumValues_LoggingLevel = []interface{}{
	"alert",
	"critical",
	"debug",
	"emergency",
	"error",
	"info",
	"notice",
	"warning",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *LoggingLevel) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_LoggingLevel {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_LoggingLevel, v)
	}
	*j = LoggingLevel(v)
	return nil
}

// Notification of a log message passed from server to client. If no
// logging/setLevel request has been sent from the client, the server MAY decide
// which messages to send automatically.
type LoggingMessageNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params LoggingMessageNotificationParams `json:"params" yaml:"params" mapstructure:"params"`
}

type LoggingMessageNotificationParams struct {
	// The data to be logged, such as a string message or an object. Any JSON
	// serializable type is allowed here.
	Data interface{} `json:"data" yaml:"data" mapstructure:"data"`

	// The severity of this log message.
	Level LoggingLevel `json:"level" yaml:"level" mapstructure:"level"`

	// An optional name of the logger issuing this message.
	Logger *string `json:"logger,omitempty" yaml:"logger,omitempty" mapstructure:"logger,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *LoggingMessageNotificationParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["data"]; raw != nil && !ok {
		return fmt.Errorf("field data in LoggingMessageNotificationParams: required")
	}
	if _, ok := raw["level"]; raw != nil && !ok {
		return fmt.Errorf("field level in LoggingMessageNotificationParams: required")
	}
	type Plain LoggingMessageNotificationParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = LoggingMessageNotificationParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *LoggingMessageNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in LoggingMessageNotification: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in LoggingMessageNotification: required")
	}
	type Plain LoggingMessageNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = LoggingMessageNotification(plain)
	return nil
}
