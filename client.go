package mcp_golang

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Base for objects that include optional annotations for the client. The client
// can use annotations to inform how objects are used or displayed
type Annotated struct {
	// Annotations corresponds to the JSON schema field "annotations".
	Annotations *AnnotatedAnnotations `json:"annotations,omitempty" yaml:"annotations,omitempty" mapstructure:"annotations,omitempty"`
}

type AnnotatedAnnotations struct {
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
func (j *AnnotatedAnnotations) UnmarshalJSON(b []byte) error {
	type Plain AnnotatedAnnotations
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
	*j = AnnotatedAnnotations(plain)
	return nil
}

// Used by the client to invoke a tool provided by the server.
type CallToolRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params CallToolRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type CallToolRequestParams struct {
	// Arguments corresponds to the JSON schema field "arguments".
	Arguments CallToolRequestParamsArguments `json:"arguments,omitempty" yaml:"arguments,omitempty" mapstructure:"arguments,omitempty"`

	// Name corresponds to the JSON schema field "name".
	Name string `json:"name" yaml:"name" mapstructure:"name"`
}

type CallToolRequestParamsArguments map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CallToolRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in CallToolRequestParams: required")
	}
	type Plain CallToolRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CallToolRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CallToolRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in CallToolRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in CallToolRequest: required")
	}
	type Plain CallToolRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CallToolRequest(plain)
	return nil
}

// Used by the client to get a prompt provided by the server.
type GetPromptRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params GetPromptRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type GetPromptRequestParams struct {
	// Arguments to use for templating the prompt.
	Arguments GetPromptRequestParamsArguments `json:"arguments,omitempty" yaml:"arguments,omitempty" mapstructure:"arguments,omitempty"`

	// The name of the prompt or prompt template.
	Name string `json:"name" yaml:"name" mapstructure:"name"`
}

// Arguments to use for templating the prompt.
type GetPromptRequestParamsArguments map[string]string

// UnmarshalJSON implements json.Unmarshaler.
func (j *GetPromptRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in GetPromptRequestParams: required")
	}
	type Plain GetPromptRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = GetPromptRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *GetPromptRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in GetPromptRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in GetPromptRequest: required")
	}
	type Plain GetPromptRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = GetPromptRequest(plain)
	return nil
}

// This request is sent from the client to the server when it first connects,
// asking it to begin initialization.
type InitializeRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params InitializeRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type InitializeRequestParams struct {
	// Capabilities corresponds to the JSON schema field "capabilities".
	Capabilities ClientCapabilities `json:"capabilities" yaml:"capabilities" mapstructure:"capabilities"`

	// ClientInfo corresponds to the JSON schema field "clientInfo".
	ClientInfo Implementation `json:"clientInfo" yaml:"clientInfo" mapstructure:"clientInfo"`

	// The latest version of the Model Context Protocol that the client supports. The
	// client MAY decide to support older versions as well.
	ProtocolVersion string `json:"protocolVersion" yaml:"protocolVersion" mapstructure:"protocolVersion"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *InitializeRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["capabilities"]; raw != nil && !ok {
		return fmt.Errorf("field capabilities in InitializeRequestParams: required")
	}
	if _, ok := raw["clientInfo"]; raw != nil && !ok {
		return fmt.Errorf("field clientInfo in InitializeRequestParams: required")
	}
	if _, ok := raw["protocolVersion"]; raw != nil && !ok {
		return fmt.Errorf("field protocolVersion in InitializeRequestParams: required")
	}
	type Plain InitializeRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = InitializeRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *InitializeRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in InitializeRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in InitializeRequest: required")
	}
	type Plain InitializeRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = InitializeRequest(plain)
	return nil
}

// This notification is sent from the client to the server after initialization has
// finished.
type InitializedNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *InitializedNotificationParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type InitializedNotificationParams struct {
	// This parameter name is reserved by MCP to allow clients and servers to attach
	// additional metadata to their notifications.
	Meta InitializedNotificationParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

// This parameter name is reserved by MCP to allow clients and servers to attach
// additional metadata to their notifications.
type InitializedNotificationParamsMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *InitializedNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in InitializedNotification: required")
	}
	type Plain InitializedNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = InitializedNotification(plain)
	return nil
}

// Sent from the client to request a list of prompts and prompt templates the
// server has.
type ListPromptsRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *ListPromptsRequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type ListPromptsRequestParams struct {
	// An opaque token representing the current pagination position.
	// If provided, the server should return results starting after this cursor.
	Cursor *string `json:"cursor,omitempty" yaml:"cursor,omitempty" mapstructure:"cursor,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListPromptsRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ListPromptsRequest: required")
	}
	type Plain ListPromptsRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListPromptsRequest(plain)
	return nil
}

// Capabilities a client may support. Known capabilities are defined here, in this
// schema, but this is not a closed set: any client can define its own, additional
// capabilities.
type ClientCapabilities struct {
	// Experimental, non-standard capabilities that the client supports.
	Experimental ClientCapabilitiesExperimental `json:"experimental,omitempty" yaml:"experimental,omitempty" mapstructure:"experimental,omitempty"`

	// Present if the client supports listing roots.
	Roots *ClientCapabilitiesRoots `json:"roots,omitempty" yaml:"roots,omitempty" mapstructure:"roots,omitempty"`

	// Present if the client supports sampling from an LLM.
	Sampling ClientCapabilitiesSampling `json:"sampling,omitempty" yaml:"sampling,omitempty" mapstructure:"sampling,omitempty"`
}

// Experimental, non-standard capabilities that the client supports.
type ClientCapabilitiesExperimental map[string]map[string]interface{}

// Present if the client supports listing roots.
type ClientCapabilitiesRoots struct {
	// Whether the client supports notifications for changes to the roots list.
	ListChanged *bool `json:"listChanged,omitempty" yaml:"listChanged,omitempty" mapstructure:"listChanged,omitempty"`
}

// Present if the client supports sampling from an LLM.
type ClientCapabilitiesSampling map[string]interface{}

type ClientNotification interface{}

type ClientRequest interface{}

type ClientResult interface{}

// A request from the client to the server, to ask for completion options.
type CompleteRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params CompleteRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type CompleteRequestParams struct {
	// The argument's information
	Argument CompleteRequestParamsArgument `json:"argument" yaml:"argument" mapstructure:"argument"`

	// Ref corresponds to the JSON schema field "ref".
	Ref interface{} `json:"ref" yaml:"ref" mapstructure:"ref"`
}

// The argument's information
type CompleteRequestParamsArgument struct {
	// The name of the argument
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// The value of the argument to use for completion matching.
	Value string `json:"value" yaml:"value" mapstructure:"value"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CompleteRequestParamsArgument) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in CompleteRequestParamsArgument: required")
	}
	if _, ok := raw["value"]; raw != nil && !ok {
		return fmt.Errorf("field value in CompleteRequestParamsArgument: required")
	}
	type Plain CompleteRequestParamsArgument
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CompleteRequestParamsArgument(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CompleteRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["argument"]; raw != nil && !ok {
		return fmt.Errorf("field argument in CompleteRequestParams: required")
	}
	if _, ok := raw["ref"]; raw != nil && !ok {
		return fmt.Errorf("field ref in CompleteRequestParams: required")
	}
	type Plain CompleteRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CompleteRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CompleteRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in CompleteRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in CompleteRequest: required")
	}
	type Plain CompleteRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CompleteRequest(plain)
	return nil
}

// A request from the server to sample an LLM via the client. The client has full
// discretion over which model to select. The client should also inform the user
// before beginning sampling, to allow them to inspect the request (human in the
// loop) and decide whether to approve it.
type CreateMessageRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params CreateMessageRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type CreateMessageRequestParams struct {
	// A request to include context from one or more MCP servers (including the
	// caller), to be attached to the prompt. The client MAY ignore this request.
	IncludeContext *CreateMessageRequestParamsIncludeContext `json:"includeContext,omitempty" yaml:"includeContext,omitempty" mapstructure:"includeContext,omitempty"`

	// The maximum number of tokens to sample, as requested by the server. The client
	// MAY choose to sample fewer tokens than requested.
	MaxTokens int `json:"maxTokens" yaml:"maxTokens" mapstructure:"maxTokens"`

	// Messages corresponds to the JSON schema field "messages".
	Messages []SamplingMessage `json:"messages" yaml:"messages" mapstructure:"messages"`

	// Optional metadata to pass through to the LLM provider. The format of this
	// metadata is provider-specific.
	Metadata CreateMessageRequestParamsMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty" mapstructure:"metadata,omitempty"`

	// The server's preferences for which model to select. The client MAY ignore these
	// preferences.
	ModelPreferences *ModelPreferences `json:"modelPreferences,omitempty" yaml:"modelPreferences,omitempty" mapstructure:"modelPreferences,omitempty"`

	// StopSequences corresponds to the JSON schema field "stopSequences".
	StopSequences []string `json:"stopSequences,omitempty" yaml:"stopSequences,omitempty" mapstructure:"stopSequences,omitempty"`

	// An optional system prompt the server wants to use for sampling. The client MAY
	// modify or omit this prompt.
	SystemPrompt *string `json:"systemPrompt,omitempty" yaml:"systemPrompt,omitempty" mapstructure:"systemPrompt,omitempty"`

	// Temperature corresponds to the JSON schema field "temperature".
	Temperature *float64 `json:"temperature,omitempty" yaml:"temperature,omitempty" mapstructure:"temperature,omitempty"`
}

type CreateMessageRequestParamsIncludeContext string

const CreateMessageRequestParamsIncludeContextAllServers CreateMessageRequestParamsIncludeContext = "allServers"
const CreateMessageRequestParamsIncludeContextNone CreateMessageRequestParamsIncludeContext = "none"
const CreateMessageRequestParamsIncludeContextThisServer CreateMessageRequestParamsIncludeContext = "thisServer"

var enumValues_CreateMessageRequestParamsIncludeContext = []interface{}{
	"allServers",
	"none",
	"thisServer",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CreateMessageRequestParamsIncludeContext) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_CreateMessageRequestParamsIncludeContext {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_CreateMessageRequestParamsIncludeContext, v)
	}
	*j = CreateMessageRequestParamsIncludeContext(v)
	return nil
}

// Optional metadata to pass through to the LLM provider. The format of this
// metadata is provider-specific.
type CreateMessageRequestParamsMetadata map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CreateMessageRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["maxTokens"]; raw != nil && !ok {
		return fmt.Errorf("field maxTokens in CreateMessageRequestParams: required")
	}
	if _, ok := raw["messages"]; raw != nil && !ok {
		return fmt.Errorf("field messages in CreateMessageRequestParams: required")
	}
	type Plain CreateMessageRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CreateMessageRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CreateMessageRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in CreateMessageRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in CreateMessageRequest: required")
	}
	type Plain CreateMessageRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CreateMessageRequest(plain)
	return nil
}

// The client's response to a sampling/create_message request from the server. The
// client should inform the user before returning the sampled message, to allow
// them to inspect the response (human in the loop) and decide whether to allow the
// server to see it.
type CreateMessageResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta CreateMessageResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// Content corresponds to the JSON schema field "content".
	Content interface{} `json:"content" yaml:"content" mapstructure:"content"`

	// The name of the model that generated the message.
	Model string `json:"model" yaml:"model" mapstructure:"model"`

	// Role corresponds to the JSON schema field "role".
	Role Role `json:"role" yaml:"role" mapstructure:"role"`

	// The reason why sampling stopped, if known.
	StopReason *string `json:"stopReason,omitempty" yaml:"stopReason,omitempty" mapstructure:"stopReason,omitempty"`
}

// Sent from the client to request a list of resource templates the server has.
type ListResourceTemplatesRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *ListResourceTemplatesRequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type ListResourceTemplatesRequestParams struct {
	// An opaque token representing the current pagination position.
	// If provided, the server should return results starting after this cursor.
	Cursor *string `json:"cursor,omitempty" yaml:"cursor,omitempty" mapstructure:"cursor,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListResourceTemplatesRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ListResourceTemplatesRequest: required")
	}
	type Plain ListResourceTemplatesRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListResourceTemplatesRequest(plain)
	return nil
}

// Sent from the client to request a list of resources the server has.
type ListResourcesRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *ListResourcesRequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type ListResourcesRequestParams struct {
	// An opaque token representing the current pagination position.
	// If provided, the server should return results starting after this cursor.
	Cursor *string `json:"cursor,omitempty" yaml:"cursor,omitempty" mapstructure:"cursor,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListResourcesRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ListResourcesRequest: required")
	}
	type Plain ListResourcesRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListResourcesRequest(plain)
	return nil
}

// The client's response to a roots/list request from the server.
// This result contains an array of Root objects, each representing a root
// directory
// or file that the server can operate on.
type ListRootsResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta ListRootsResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// Roots corresponds to the JSON schema field "roots".
	Roots []Root `json:"roots" yaml:"roots" mapstructure:"roots"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type ListRootsResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListRootsResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["roots"]; raw != nil && !ok {
		return fmt.Errorf("field roots in ListRootsResult: required")
	}
	type Plain ListRootsResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListRootsResult(plain)
	return nil
}

// Sent from the client to request a list of tools the server has.
type ListToolsRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *ListToolsRequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type ListToolsRequestParams struct {
	// An opaque token representing the current pagination position.
	// If provided, the server should return results starting after this cursor.
	Cursor *string `json:"cursor,omitempty" yaml:"cursor,omitempty" mapstructure:"cursor,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListToolsRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ListToolsRequest: required")
	}
	type Plain ListToolsRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListToolsRequest(plain)
	return nil
}

// Sent from the client to the server, to read a specific resource URI.
type ReadResourceRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params ReadResourceRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}

type ReadResourceRequestParams struct {
	// The URI of the resource to read. The URI can use any protocol; it is up to the
	// server how to interpret it.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ReadResourceRequestParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in ReadResourceRequestParams: required")
	}
	type Plain ReadResourceRequestParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ReadResourceRequestParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ReadResourceRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ReadResourceRequest: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in ReadResourceRequest: required")
	}
	type Plain ReadResourceRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ReadResourceRequest(plain)
	return nil
}
