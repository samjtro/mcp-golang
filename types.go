// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package mcp_golang

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type CallToolResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CallToolResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["content"]; raw != nil && !ok {
		return fmt.Errorf("field content in CallToolResult: required")
	}
	type Plain CallToolResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CallToolResult(plain)
	return nil
}

// This notification can be sent by either side to indicate that it is cancelling a
// previously-issued request.
//
// The request SHOULD still be in-flight, but due to communication latency, it is
// always possible that this notification MAY arrive after the request has already
// finished.
//
// This notification indicates that the result will be unused, so any associated
// processing SHOULD cease.
//
// A client MUST NOT attempt to cancel its `initialize` request.
type CancelledNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params CancelledNotificationParams `json:"params" yaml:"params" mapstructure:"params"`
}

type CancelledNotificationParams struct {
	// An optional string describing the reason for the cancellation. This MAY be
	// logged or presented to the user.
	Reason *string `json:"reason,omitempty" yaml:"reason,omitempty" mapstructure:"reason,omitempty"`

	// The ID of the request to cancel.
	//
	// This MUST correspond to the ID of a request previously issued in the same
	// direction.
	RequestId RequestId `json:"requestId" yaml:"requestId" mapstructure:"requestId"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CancelledNotificationParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["requestId"]; raw != nil && !ok {
		return fmt.Errorf("field requestId in CancelledNotificationParams: required")
	}
	type Plain CancelledNotificationParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CancelledNotificationParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CancelledNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in CancelledNotification: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in CancelledNotification: required")
	}
	type Plain CancelledNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CancelledNotification(plain)
	return nil
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type CompleteResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CompleteResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["completion"]; raw != nil && !ok {
		return fmt.Errorf("field completion in CompleteResult: required")
	}
	type Plain CompleteResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CompleteResult(plain)
	return nil
}


// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type CreateMessageResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CreateMessageResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["content"]; raw != nil && !ok {
		return fmt.Errorf("field content in CreateMessageResult: required")
	}
	if _, ok := raw["model"]; raw != nil && !ok {
		return fmt.Errorf("field model in CreateMessageResult: required")
	}
	if _, ok := raw["role"]; raw != nil && !ok {
		return fmt.Errorf("field role in CreateMessageResult: required")
	}
	type Plain CreateMessageResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CreateMessageResult(plain)
	return nil
}

// An opaque token used to represent a cursor for pagination.
type Cursor string

type EmbeddedResourceAnnotations struct {
	// Describes who the intended customer of this object or data is.
	//
	// It can include multiple entries to indicate content useful for multiple
	// audiences (e.g., `["user", "assistant"]`).
	Audience []Role `json:"audience,omitempty" yaml:"audience,omitempty" mapstructure:"audience,omitempty"`

	// Describes how important this data is for operating the server.
	//
	// A value of 1 means "most important," and indicates that the data is
	// effectively required, while 0 means "least important," and indicates that
	// the data is entirely optional.
	Priority *float64 `json:"priority,omitempty" yaml:"priority,omitempty" mapstructure:"priority,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *EmbeddedResourceAnnotations) UnmarshalJSON(b []byte) error {
	type Plain EmbeddedResourceAnnotations
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Priority != nil && 1 < *plain.Priority {
		return fmt.Errorf("field %s: must be <= %v", "priority", 1)
	}
	if plain.Priority != nil && 0 > *plain.Priority {
		return fmt.Errorf("field %s: must be >= %v", "priority", 0)
	}
	*j = EmbeddedResourceAnnotations(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *EmbeddedResource) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["resource"]; raw != nil && !ok {
		return fmt.Errorf("field resource in EmbeddedResource: required")
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in EmbeddedResource: required")
	}
	type Plain EmbeddedResource
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = EmbeddedResource(plain)
	return nil
}

type ImageContentAnnotations struct {
	// Describes who the intended customer of this object or data is.
	//
	// It can include multiple entries to indicate content useful for multiple
	// audiences (e.g., `["user", "assistant"]`).
	Audience []Role `json:"audience,omitempty" yaml:"audience,omitempty" mapstructure:"audience,omitempty"`

	// Describes how important this data is for operating the server.
	//
	// A value of 1 means "most important," and indicates that the data is
	// effectively required, while 0 means "least important," and indicates that
	// the data is entirely optional.
	Priority *float64 `json:"priority,omitempty" yaml:"priority,omitempty" mapstructure:"priority,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ImageContentAnnotations) UnmarshalJSON(b []byte) error {
	type Plain ImageContentAnnotations
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Priority != nil && 1 < *plain.Priority {
		return fmt.Errorf("field %s: must be <= %v", "priority", 1)
	}
	if plain.Priority != nil && 0 > *plain.Priority {
		return fmt.Errorf("field %s: must be >= %v", "priority", 0)
	}
	*j = ImageContentAnnotations(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ImageContent) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["data"]; raw != nil && !ok {
		return fmt.Errorf("field data in ImageContent: required")
	}
	if _, ok := raw["mimeType"]; raw != nil && !ok {
		return fmt.Errorf("field mimeType in ImageContent: required")
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in ImageContent: required")
	}
	type Plain ImageContent
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ImageContent(plain)
	return nil
}

// Describes the name and version of an MCP implementation.
type Implementation struct {
	// Name corresponds to the JSON schema field "name".
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Version corresponds to the JSON schema field "version".
	Version string `json:"version" yaml:"version" mapstructure:"version"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Implementation) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in Implementation: required")
	}
	if _, ok := raw["version"]; raw != nil && !ok {
		return fmt.Errorf("field version in Implementation: required")
	}
	type Plain Implementation
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Implementation(plain)
	return nil
}

// Hints to use for model selection.
//
// Keys not declared here are currently left unspecified by the spec and are up
// to the client to interpret.
type ModelHint struct {
	// A hint for a model name.
	//
	// The client SHOULD treat this as a substring of a model name; for example:
	//  - `claude-3-5-sonnet` should match `claude-3-5-sonnet-20241022`
	//  - `sonnet` should match `claude-3-5-sonnet-20241022`,
	// `claude-3-sonnet-20240229`, etc.
	//  - `claude` should match any Claude model
	//
	// The client MAY also map the string to a different provider's model name or a
	// different model family, as long as it fills a similar niche; for example:
	//  - `gemini-1.5-flash` could match `claude-3-haiku-20240307`
	Name *string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
}

// The server's preferences for model selection, requested of the client during
// sampling.
//
// Because LLMs can vary along multiple dimensions, choosing the "best" model is
// rarely straightforward.  Different models excel in different areas—some are
// faster but less capable, others are more capable but more expensive, and so
// on. This interface allows servers to express their priorities across multiple
// dimensions to help clients make an appropriate selection for their use case.
//
// These preferences are always advisory. The client MAY ignore them. It is also
// up to the client to decide how to interpret these preferences and how to
// balance them against other considerations.
type ModelPreferences struct {
	// How much to prioritize cost when selecting a model. A value of 0 means cost
	// is not important, while a value of 1 means cost is the most important
	// factor.
	CostPriority *float64 `json:"costPriority,omitempty" yaml:"costPriority,omitempty" mapstructure:"costPriority,omitempty"`

	// Optional hints to use for model selection.
	//
	// If multiple hints are specified, the client MUST evaluate them in order
	// (such that the first match is taken).
	//
	// The client SHOULD prioritize these hints over the numeric priorities, but
	// MAY still use the priorities to select from ambiguous matches.
	Hints []ModelHint `json:"hints,omitempty" yaml:"hints,omitempty" mapstructure:"hints,omitempty"`

	// How much to prioritize intelligence and capabilities when selecting a
	// model. A value of 0 means intelligence is not important, while a value of 1
	// means intelligence is the most important factor.
	IntelligencePriority *float64 `json:"intelligencePriority,omitempty" yaml:"intelligencePriority,omitempty" mapstructure:"intelligencePriority,omitempty"`

	// How much to prioritize sampling speed (latency) when selecting a model. A
	// value of 0 means speed is not important, while a value of 1 means speed is
	// the most important factor.
	SpeedPriority *float64 `json:"speedPriority,omitempty" yaml:"speedPriority,omitempty" mapstructure:"speedPriority,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ModelPreferences) UnmarshalJSON(b []byte) error {
	type Plain ModelPreferences
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.CostPriority != nil && 1 < *plain.CostPriority {
		return fmt.Errorf("field %s: must be <= %v", "costPriority", 1)
	}
	if plain.CostPriority != nil && 0 > *plain.CostPriority {
		return fmt.Errorf("field %s: must be >= %v", "costPriority", 0)
	}
	if plain.IntelligencePriority != nil && 1 < *plain.IntelligencePriority {
		return fmt.Errorf("field %s: must be <= %v", "intelligencePriority", 1)
	}
	if plain.IntelligencePriority != nil && 0 > *plain.IntelligencePriority {
		return fmt.Errorf("field %s: must be >= %v", "intelligencePriority", 0)
	}
	if plain.SpeedPriority != nil && 1 < *plain.SpeedPriority {
		return fmt.Errorf("field %s: must be <= %v", "speedPriority", 1)
	}
	if plain.SpeedPriority != nil && 0 > *plain.SpeedPriority {
		return fmt.Errorf("field %s: must be >= %v", "speedPriority", 0)
	}
	*j = ModelPreferences(plain)
	return nil
}

type Notification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *NotificationParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type NotificationParams struct {
	// This parameter name is reserved by MCP to allow clients and servers to attach
	// additional metadata to their notifications.
	Meta NotificationParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

// This parameter name is reserved by MCP to allow clients and servers to attach
// additional metadata to their notifications.
type NotificationParamsMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Notification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in Notification: required")
	}
	type Plain Notification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Notification(plain)
	return nil
}

type PaginatedRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *PaginatedRequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type PaginatedRequestParams struct {
	// An opaque token representing the current pagination position.
	// If provided, the server should return results starting after this cursor.
	Cursor *string `json:"cursor,omitempty" yaml:"cursor,omitempty" mapstructure:"cursor,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PaginatedRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in PaginatedRequest: required")
	}
	type Plain PaginatedRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PaginatedRequest(plain)
	return nil
}

type PaginatedResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta PaginatedResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// An opaque token representing the pagination position after the last returned
	// result.
	// If present, there may be more results available.
	NextCursor *string `json:"nextCursor,omitempty" yaml:"nextCursor,omitempty" mapstructure:"nextCursor,omitempty"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type PaginatedResultMeta map[string]interface{}

// A ping, issued by either the server or the client, to check that the other party
// is still alive. The receiver must promptly respond, or else may be disconnected.
type PingRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *PingRequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type PingRequestParams struct {
	// Meta corresponds to the JSON schema field "_meta".
	Meta *PingRequestParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

type PingRequestParamsMeta struct {
	// If specified, the caller is requesting out-of-band progress notifications for
	// this request (as represented by notifications/progress). The value of this
	// parameter is an opaque token that will be attached to any subsequent
	// notifications. The receiver is not obligated to provide these notifications.
	ProgressToken *ProgressToken `json:"progressToken,omitempty" yaml:"progressToken,omitempty" mapstructure:"progressToken,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PingRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in PingRequest: required")
	}
	type Plain PingRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PingRequest(plain)
	return nil
}

// An out-of-band notification used to inform the receiver of a progress update for
// a long-running request.
type ProgressNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params ProgressNotificationParams `json:"params" yaml:"params" mapstructure:"params"`
}

type ProgressNotificationParams struct {
	// The progress thus far. This should increase every time progress is made, even
	// if the total is unknown.
	Progress float64 `json:"progress" yaml:"progress" mapstructure:"progress"`

	// The progress token which was given in the initial request, used to associate
	// this notification with the request that is proceeding.
	ProgressToken ProgressToken `json:"progressToken" yaml:"progressToken" mapstructure:"progressToken"`

	// Total number of items to process (or total progress required), if known.
	Total *float64 `json:"total,omitempty" yaml:"total,omitempty" mapstructure:"total,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ProgressNotificationParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["progress"]; raw != nil && !ok {
		return fmt.Errorf("field progress in ProgressNotificationParams: required")
	}
	if _, ok := raw["progressToken"]; raw != nil && !ok {
		return fmt.Errorf("field progressToken in ProgressNotificationParams: required")
	}
	type Plain ProgressNotificationParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ProgressNotificationParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ProgressNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ProgressNotification: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in ProgressNotification: required")
	}
	type Plain ProgressNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ProgressNotification(plain)
	return nil
}

// A progress token, used to associate progress notifications with the original
// request.
type ProgressToken int

// A prompt or prompt template that the server offers.
type Prompt struct {
	// A list of arguments to use for templating the prompt.
	Arguments []PromptArgument `json:"arguments,omitempty" yaml:"arguments,omitempty" mapstructure:"arguments,omitempty"`

	// An optional description of what this prompt provides
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// The name of the prompt or prompt template.
	Name string `json:"name" yaml:"name" mapstructure:"name"`
}

// Describes an argument that a prompt can accept.
type PromptArgument struct {
	// A human-readable description of the argument.
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// The name of the argument.
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Whether this argument must be provided.
	Required *bool `json:"required,omitempty" yaml:"required,omitempty" mapstructure:"required,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PromptArgument) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in PromptArgument: required")
	}
	type Plain PromptArgument
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PromptArgument(plain)
	return nil
}

// An optional notification from the server to the client, informing it that the
// list of prompts it offers has changed. This may be issued by servers without any
// previous subscription from the client.
type PromptListChangedNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *PromptListChangedNotificationParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type PromptListChangedNotificationParams struct {
	// This parameter name is reserved by MCP to allow clients and servers to attach
	// additional metadata to their notifications.
	Meta PromptListChangedNotificationParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

// This parameter name is reserved by MCP to allow clients and servers to attach
// additional metadata to their notifications.
type PromptListChangedNotificationParamsMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PromptListChangedNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in PromptListChangedNotification: required")
	}
	type Plain PromptListChangedNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PromptListChangedNotification(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PromptMessage) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["content"]; raw != nil && !ok {
		return fmt.Errorf("field content in PromptMessage: required")
	}
	if _, ok := raw["role"]; raw != nil && !ok {
		return fmt.Errorf("field role in PromptMessage: required")
	}
	type Plain PromptMessage
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PromptMessage(plain)
	return nil
}

// Identifies a prompt.
type PromptReference struct {
	// The name of the prompt or prompt template
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Type corresponds to the JSON schema field "type".
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *PromptReference) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in PromptReference: required")
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in PromptReference: required")
	}
	type Plain PromptReference
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = PromptReference(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Prompt) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in Prompt: required")
	}
	type Plain Prompt
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Prompt(plain)
	return nil
}

// The server's response to a resources/read request from the client.
type ReadResourceResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta ReadResourceResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// Contents corresponds to the JSON schema field "contents".
	Contents []interface{} `json:"contents" yaml:"contents" mapstructure:"contents"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type ReadResourceResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ReadResourceResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["contents"]; raw != nil && !ok {
		return fmt.Errorf("field contents in ReadResourceResult: required")
	}
	type Plain ReadResourceResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ReadResourceResult(plain)
	return nil
}

type Request struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *RequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

// A uniquely identifying ID for a request in JSON-RPC.
type RequestId int

type RequestParams struct {
	// Meta corresponds to the JSON schema field "_meta".
	Meta *RequestParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

type RequestParamsMeta struct {
	// If specified, the caller is requesting out-of-band progress notifications for
	// this request (as represented by notifications/progress). The value of this
	// parameter is an opaque token that will be attached to any subsequent
	// notifications. The receiver is not obligated to provide these notifications.
	ProgressToken *ProgressToken `json:"progressToken,omitempty" yaml:"progressToken,omitempty" mapstructure:"progressToken,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Request) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in Request: required")
	}
	type Plain Request
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Request(plain)
	return nil
}

// A known resource that the server is capable of reading.
type Resource struct {
	// Annotations corresponds to the JSON schema field "annotations".
	Annotations *ResourceAnnotations `json:"annotations,omitempty" yaml:"annotations,omitempty" mapstructure:"annotations,omitempty"`

	// A description of what this resource represents.
	//
	// This can be used by clients to improve the LLM's understanding of available
	// resources. It can be thought of like a "hint" to the model.
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// The MIME type of this resource, if known.
	MimeType *string `json:"mimeType,omitempty" yaml:"mimeType,omitempty" mapstructure:"mimeType,omitempty"`

	// A human-readable name for this resource.
	//
	// This can be used by clients to populate UI elements.
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// The URI of this resource.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

type ResourceAnnotations struct {
	// Describes who the intended customer of this object or data is.
	//
	// It can include multiple entries to indicate content useful for multiple
	// audiences (e.g., `["user", "assistant"]`).
	Audience []Role `json:"audience,omitempty" yaml:"audience,omitempty" mapstructure:"audience,omitempty"`

	// Describes how important this data is for operating the server.
	//
	// A value of 1 means "most important," and indicates that the data is
	// effectively required, while 0 means "least important," and indicates that
	// the data is entirely optional.
	Priority *float64 `json:"priority,omitempty" yaml:"priority,omitempty" mapstructure:"priority,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceAnnotations) UnmarshalJSON(b []byte) error {
	type Plain ResourceAnnotations
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Priority != nil && 1 < *plain.Priority {
		return fmt.Errorf("field %s: must be <= %v", "priority", 1)
	}
	if plain.Priority != nil && 0 > *plain.Priority {
		return fmt.Errorf("field %s: must be >= %v", "priority", 0)
	}
	*j = ResourceAnnotations(plain)
	return nil
}

// The contents of a specific resource or sub-resource.
type ResourceContents struct {
	// The MIME type of this resource, if known.
	MimeType *string `json:"mimeType,omitempty" yaml:"mimeType,omitempty" mapstructure:"mimeType,omitempty"`

	// The URI of this resource.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceContents) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in ResourceContents: required")
	}
	type Plain ResourceContents
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ResourceContents(plain)
	return nil
}

// This parameter name is reserved by MCP to allow clients and servers to attach
// additional metadata to their notifications.
type ResourceListChangedNotificationParamsMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceListChangedNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ResourceListChangedNotification: required")
	}
	type Plain ResourceListChangedNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ResourceListChangedNotification(plain)
	return nil
}

// A reference to a resource or resource template definition.
type ResourceReference struct {
	// Type corresponds to the JSON schema field "type".
	Type string `json:"type" yaml:"type" mapstructure:"type"`

	// The URI or URI template of the resource.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceReference) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in ResourceReference: required")
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in ResourceReference: required")
	}
	type Plain ResourceReference
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ResourceReference(plain)
	return nil
}

type Result struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta ResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type ResultMeta map[string]interface{}

var enumValues_Role = []interface{}{
	"assistant",
	"user",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Role) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_Role {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_Role, v)
	}
	*j = Role(v)
	return nil
}

// Represents a root directory or file that the server can operate on.
type Root struct {
	// An optional name for the root. This can be used to provide a human-readable
	// identifier for the root, which may be useful for display purposes or for
	// referencing the root in other parts of the application.
	Name *string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`

	// The URI identifying the root. This *must* start with file:// for now.
	// This restriction may be relaxed in future versions of the protocol to allow
	// other URI schemes.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Root) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in Root: required")
	}
	type Plain Root
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Root(plain)
	return nil
}

// This parameter name is reserved by MCP to allow clients and servers to attach
// additional metadata to their notifications.
type RootsListChangedNotificationParamsMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RootsListChangedNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in RootsListChangedNotification: required")
	}
	type Plain RootsListChangedNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = RootsListChangedNotification(plain)
	return nil
}

// Describes a message issued to or received from an LLM API.
type SamplingMessage struct {
	// Content corresponds to the JSON schema field "content".
	Content interface{} `json:"content" yaml:"content" mapstructure:"content"`

	// Role corresponds to the JSON schema field "role".
	Role Role `json:"role" yaml:"role" mapstructure:"role"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SamplingMessage) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["content"]; raw != nil && !ok {
		return fmt.Errorf("field content in SamplingMessage: required")
	}
	if _, ok := raw["role"]; raw != nil && !ok {
		return fmt.Errorf("field role in SamplingMessage: required")
	}
	type Plain SamplingMessage
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = SamplingMessage(plain)
	return nil
}

type TextContentAnnotations struct {
	// Describes who the intended customer of this object or data is.
	//
	// It can include multiple entries to indicate content useful for multiple
	// audiences (e.g., `["user", "assistant"]`).
	Audience []Role `json:"audience,omitempty" yaml:"audience,omitempty" mapstructure:"audience,omitempty"`

	// Describes how important this data is for operating the server.
	//
	// A value of 1 means "most important," and indicates that the data is
	// effectively required, while 0 means "least important," and indicates that
	// the data is entirely optional.
	Priority *float64 `json:"priority,omitempty" yaml:"priority,omitempty" mapstructure:"priority,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TextContentAnnotations) UnmarshalJSON(b []byte) error {
	type Plain TextContentAnnotations
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Priority != nil && 1 < *plain.Priority {
		return fmt.Errorf("field %s: must be <= %v", "priority", 1)
	}
	if plain.Priority != nil && 0 > *plain.Priority {
		return fmt.Errorf("field %s: must be >= %v", "priority", 0)
	}
	*j = TextContentAnnotations(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TextContent) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["text"]; raw != nil && !ok {
		return fmt.Errorf("field text in TextContent: required")
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in TextContent: required")
	}
	type Plain TextContent
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = TextContent(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TextResourceContents) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["text"]; raw != nil && !ok {
		return fmt.Errorf("field text in TextResourceContents: required")
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in TextResourceContents: required")
	}
	type Plain TextResourceContents
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = TextResourceContents(plain)
	return nil
}