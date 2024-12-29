package mcp_golang

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/invopop/jsonschema"
	"github.com/metoro-io/mcp-golang/internal/datastructures"
	"github.com/metoro-io/mcp-golang/internal/protocol"
	"github.com/metoro-io/mcp-golang/internal/tools"
	"github.com/metoro-io/mcp-golang/transport"
	"reflect"
	"sort"
	"strings"
)

// Here we define the actual MCP server that users will create and run
// A server can be passed a number of handlers to handle requests from clients
// Additionally it can be parametrized by a transport. This transport will be used to actually send and receive messages.
// So for example if the stdio transport is used, the server will read from stdin and write to stdout
// If the SSE transport is used, the server will send messages over an SSE connection and receive messages from HTTP POST requests.

// The interface that we're looking to support is something like [gin](https://github.com/gin-gonic/gin)s interface

type toolResponseSent struct {
	Response *ToolResponse
	Error    error
}

// Custom JSON marshaling for ToolResponse
func (c toolResponseSent) MarshalJSON() ([]byte, error) {
	if c.Error != nil {
		errorText := c.Error.Error()
		c.Response = NewToolResponse(NewTextContent(errorText))
	}
	return json.Marshal(struct {
		Content []*Content `json:"content" yaml:"content" mapstructure:"content"`
		IsError bool       `json:"isError" yaml:"isError" mapstructure:"isError"`
	}{
		Content: c.Response.Content,
		IsError: c.Error != nil,
	})
}

// Custom JSON marshaling for ToolResponse
func (c resourceResponseSent) MarshalJSON() ([]byte, error) {
	if c.Error != nil {
		errorText := c.Error.Error()
		c.Response = NewResourceResponse(NewTextEmbeddedResource(c.Uri, errorText, "text/plain"))
	}
	return json.Marshal(c.Response)
}

type resourceResponseSent struct {
	Response *ResourceResponse
	Uri      string
	Error    error
}

func newResourceResponseSentError(err error) *resourceResponseSent {
	return &resourceResponseSent{
		Error: err,
	}
}

// newToolResponseSent creates a new toolResponseSent
func newResourceResponseSent(response *ResourceResponse) *resourceResponseSent {
	return &resourceResponseSent{
		Response: response,
	}
}

type promptResponseSent struct {
	Response *PromptResponse
	Error    error
}

func newPromptResponseSentError(err error) *promptResponseSent {
	return &promptResponseSent{
		Error: err,
	}
}

// newToolResponseSent creates a new toolResponseSent
func newPromptResponseSent(response *PromptResponse) *promptResponseSent {
	return &promptResponseSent{
		Response: response,
	}
}

// Custom JSON marshaling for PromptResponse
func (c promptResponseSent) MarshalJSON() ([]byte, error) {
	if c.Error != nil {
		errorText := c.Error.Error()
		c.Response = NewPromptResponse("error", NewPromptMessage(NewTextContent(errorText), RoleUser))
	}
	return json.Marshal(c.Response)
}

type Server struct {
	isRunning          bool
	transport          transport.Transport
	protocol           *protocol.Protocol
	paginationLimit    *int
	tools              *datastructures.SyncMap[string, *tool]
	prompts            *datastructures.SyncMap[string, *prompt]
	resources          *datastructures.SyncMap[string, *resource]
	serverInstructions *string
	serverName         string
	serverVersion      string
}

type prompt struct {
	Name              string
	Description       string
	Handler           func(baseGetPromptRequestParamsArguments) *promptResponseSent
	PromptInputSchema *promptSchema
}

type tool struct {
	Name            string
	Description     string
	Handler         func(baseCallToolRequestParams) *toolResponseSent
	ToolInputSchema *jsonschema.Schema
}

type resource struct {
	Name        string
	Description string
	Uri         string
	mimeType    string
	Handler     func() *resourceResponseSent
}

type ServerOptions func(*Server)

// The server's response to a tool call.
//
// Any errors that originate from the tool SHOULD be reported inside the result
// object, with `isError` set to true, _not_ as an MCP protocol-level error
// response. Otherwise, the LLM would not be able to see that an error occurred
// and self-correct.
//
// However, any errors in _finding_ the tool, an error indicating that the
// server does not support tool calls, or any other exceptional conditions,
// should be reported as an MCP error response.
type CallToolResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta CallToolResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// Content corresponds to the JSON schema field "content".
	Content []interface{} `json:"content" yaml:"content" mapstructure:"content"`

	// Whether the tool call ended in an error.
	//
	// If not set, this is assumed to be false (the call was successful).
	IsError *bool `json:"isError,omitempty" yaml:"isError,omitempty" mapstructure:"isError,omitempty"`
}

// The server's response to a completion/complete request
type CompleteResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta CompleteResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// Completion corresponds to the JSON schema field "completion".
	Completion CompleteResultCompletion `json:"completion" yaml:"completion" mapstructure:"completion"`
}

type CompleteResultCompletion struct {
	// Indicates whether there are additional completion options beyond those provided
	// in the current response, even if the exact total is unknown.
	HasMore *bool `json:"hasMore,omitempty" yaml:"hasMore,omitempty" mapstructure:"hasMore,omitempty"`

	// The total number of completion options available. This can exceed the number of
	// values actually sent in the response.
	Total *int `json:"total,omitempty" yaml:"total,omitempty" mapstructure:"total,omitempty"`

	// An array of completion values. Must not exceed 100 items.
	Values []string `json:"values" yaml:"values" mapstructure:"values"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CompleteResultCompletion) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["values"]; raw != nil && !ok {
		return fmt.Errorf("field values in CompleteResultCompletion: required")
	}
	type Plain CompleteResultCompletion
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CompleteResultCompletion(plain)
	return nil
}

// The server's response to a prompts/get request from the client.
type GetPromptResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta GetPromptResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// An optional description for the prompt.
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// Messages corresponds to the JSON schema field "messages".
	Messages []PromptMessage `json:"messages" yaml:"messages" mapstructure:"messages"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type GetPromptResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *GetPromptResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["messages"]; raw != nil && !ok {
		return fmt.Errorf("field messages in GetPromptResult: required")
	}
	type Plain GetPromptResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = GetPromptResult(plain)
	return nil
}

// After receiving an initialize request from the client, the server sends this
// response.
type InitializeResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta InitializeResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// Capabilities corresponds to the JSON schema field "capabilities".
	Capabilities serverCapabilities `json:"capabilities" yaml:"capabilities" mapstructure:"capabilities"`

	// Instructions describing how to use the server and its features.
	//
	// This can be used by clients to improve the LLM's understanding of available
	// tools, resources, etc. It can be thought of like a "hint" to the model. For
	// example, this information MAY be added to the system prompt.
	Instructions *string `json:"instructions,omitempty" yaml:"instructions,omitempty" mapstructure:"instructions,omitempty"`

	// The version of the Model Context Protocol that the server wants to use. This
	// may not match the version that the client requested. If the client cannot
	// support this version, it MUST disconnect.
	ProtocolVersion string `json:"protocolVersion" yaml:"protocolVersion" mapstructure:"protocolVersion"`

	// ServerInfo corresponds to the JSON schema field "serverInfo".
	ServerInfo Implementation `json:"serverInfo" yaml:"serverInfo" mapstructure:"serverInfo"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type InitializeResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *InitializeResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["capabilities"]; raw != nil && !ok {
		return fmt.Errorf("field capabilities in InitializeResult: required")
	}
	if _, ok := raw["protocolVersion"]; raw != nil && !ok {
		return fmt.Errorf("field protocolVersion in InitializeResult: required")
	}
	if _, ok := raw["serverInfo"]; raw != nil && !ok {
		return fmt.Errorf("field serverInfo in InitializeResult: required")
	}
	type Plain InitializeResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = InitializeResult(plain)
	return nil
}

// The server's response to a prompts/list request from the client.
type ListPromptsResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta ListPromptsResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// An opaque token representing the pagination position after the last returned
	// result.
	// If present, there may be more results available.
	NextCursor *string `json:"nextCursor,omitempty" yaml:"nextCursor,omitempty" mapstructure:"nextCursor,omitempty"`

	// Prompts corresponds to the JSON schema field "prompts".
	Prompts []Prompt `json:"prompts" yaml:"prompts" mapstructure:"prompts"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type ListPromptsResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListPromptsResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["prompts"]; raw != nil && !ok {
		return fmt.Errorf("field prompts in ListPromptsResult: required")
	}
	type Plain ListPromptsResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListPromptsResult(plain)
	return nil
}

// The server's response to a resources/templates/list request from the client.
type ListResourceTemplatesResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta ListResourceTemplatesResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// An opaque token representing the pagination position after the last returned
	// result.
	// If present, there may be more results available.
	NextCursor *string `json:"nextCursor,omitempty" yaml:"nextCursor,omitempty" mapstructure:"nextCursor,omitempty"`

	// ResourceTemplates corresponds to the JSON schema field "resourceTemplates".
	ResourceTemplates []ResourceTemplate `json:"resourceTemplates" yaml:"resourceTemplates" mapstructure:"resourceTemplates"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type ListResourceTemplatesResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListResourceTemplatesResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["resourceTemplates"]; raw != nil && !ok {
		return fmt.Errorf("field resourceTemplates in ListResourceTemplatesResult: required")
	}
	type Plain ListResourceTemplatesResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListResourceTemplatesResult(plain)
	return nil
}

// The server's response to a resources/list request from the client.
type ListResourcesResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta ListResourcesResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// An opaque token representing the pagination position after the last returned
	// result.
	// If present, there may be more results available.
	NextCursor *string `json:"nextCursor,omitempty" yaml:"nextCursor,omitempty" mapstructure:"nextCursor,omitempty"`

	// Resources corresponds to the JSON schema field "resources".
	Resources []Resource `json:"resources" yaml:"resources" mapstructure:"resources"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type ListResourcesResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListResourcesResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["resources"]; raw != nil && !ok {
		return fmt.Errorf("field resources in ListResourcesResult: required")
	}
	type Plain ListResourcesResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListResourcesResult(plain)
	return nil
}

// Sent from the server to request a list of root URIs from the client. Roots allow
// servers to ask for specific directories or files to operate on. A common example
// for roots is providing a set of repositories or directories a server should
// operate
// on.
//
// This request is typically used when the server needs to understand the file
// system
// structure or access specific locations that the client has permission to read
// from.
type ListRootsRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *ListRootsRequestParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type ListRootsRequestParams struct {
	// Meta corresponds to the JSON schema field "_meta".
	Meta *ListRootsRequestParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

type ListRootsRequestParamsMeta struct {
	// If specified, the caller is requesting out-of-band progress notifications for
	// this request (as represented by notifications/progress). The value of this
	// parameter is an opaque token that will be attached to any subsequent
	// notifications. The receiver is not obligated to provide these notifications.
	ProgressToken *ProgressToken `json:"progressToken,omitempty" yaml:"progressToken,omitempty" mapstructure:"progressToken,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListRootsRequest) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ListRootsRequest: required")
	}
	type Plain ListRootsRequest
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListRootsRequest(plain)
	return nil
}

// The server's response to a tools/list request from the client.
type ListToolsResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta ListToolsResultMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// An opaque token representing the pagination position after the last returned
	// result.
	// If present, there may be more results available.
	NextCursor *string `json:"nextCursor,omitempty" yaml:"nextCursor,omitempty" mapstructure:"nextCursor,omitempty"`

	// Tools corresponds to the JSON schema field "tools".
	Tools []Tool `json:"tools" yaml:"tools" mapstructure:"tools"`
}

// This result property is reserved by the protocol to allow clients and servers to
// attach additional metadata to their responses.
type ListToolsResultMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ListToolsResult) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["tools"]; raw != nil && !ok {
		return fmt.Errorf("field tools in ListToolsResult: required")
	}
	type Plain ListToolsResult
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ListToolsResult(plain)
	return nil
}

// An optional notification from the server to the client, informing it that the
// list of resources it can read from has changed. This may be issued by servers
// without any previous subscription from the client.
type ResourceListChangedNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *ResourceListChangedNotificationParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type ResourceListChangedNotificationParams struct {
	// This parameter name is reserved by MCP to allow clients and servers to attach
	// additional metadata to their notifications.
	Meta ResourceListChangedNotificationParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

// A template description for resources available on the server.
type ResourceTemplate struct {
	// Annotations corresponds to the JSON schema field "annotations".
	Annotations *ResourceTemplateAnnotations `json:"annotations,omitempty" yaml:"annotations,omitempty" mapstructure:"annotations,omitempty"`

	// A description of what this template is for.
	//
	// This can be used by clients to improve the LLM's understanding of available
	// resources. It can be thought of like a "hint" to the model.
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// The MIME type for all resources that match this template. This should only be
	// included if all resources matching this template have the same type.
	MimeType *string `json:"mimeType,omitempty" yaml:"mimeType,omitempty" mapstructure:"mimeType,omitempty"`

	// A human-readable name for the type of resource this template refers to.
	//
	// This can be used by clients to populate UI elements.
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// A URI template (according to RFC 6570) that can be used to construct resource
	// URIs.
	UriTemplate string `json:"uriTemplate" yaml:"uriTemplate" mapstructure:"uriTemplate"`
}

type ResourceTemplateAnnotations struct {
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
func (j *ResourceTemplateAnnotations) UnmarshalJSON(b []byte) error {
	type Plain ResourceTemplateAnnotations
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
	*j = ResourceTemplateAnnotations(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceTemplate) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in ResourceTemplate: required")
	}
	if _, ok := raw["uriTemplate"]; raw != nil && !ok {
		return fmt.Errorf("field uriTemplate in ResourceTemplate: required")
	}
	type Plain ResourceTemplate
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ResourceTemplate(plain)
	return nil
}

type ServerNotification interface{}

type ServerRequest interface{}

type ServerResult interface{}

// A notification from the server to the client, informing it that a resource has
// changed and may need to be read again. This should only be sent if the client
// previously sent a resources/subscribe request.
type ResourceUpdatedNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params ResourceUpdatedNotificationParams `json:"params" yaml:"params" mapstructure:"params"`
}

type ResourceUpdatedNotificationParams struct {
	// The URI of the resource that has been updated. This might be a sub-resource of
	// the one that the client actually subscribed to.
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceUpdatedNotificationParams) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in ResourceUpdatedNotificationParams: required")
	}
	type Plain ResourceUpdatedNotificationParams
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ResourceUpdatedNotificationParams(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ResourceUpdatedNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ResourceUpdatedNotification: required")
	}
	if _, ok := raw["params"]; raw != nil && !ok {
		return fmt.Errorf("field params in ResourceUpdatedNotification: required")
	}
	type Plain ResourceUpdatedNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ResourceUpdatedNotification(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Resource) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in Resource: required")
	}
	if _, ok := raw["uri"]; raw != nil && !ok {
		return fmt.Errorf("field uri in Resource: required")
	}
	type Plain Resource
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Resource(plain)
	return nil
}

func WithProtocol(protocol *protocol.Protocol) ServerOptions {
	return func(s *Server) {
		s.protocol = protocol
	}
}

// Beware: As of 2024-12-13, it looks like Claude does not support pagination yet
func WithPaginationLimit(limit int) ServerOptions {
	return func(s *Server) {
		s.paginationLimit = &limit
	}
}

func NewServer(transport transport.Transport, options ...ServerOptions) *Server {
	server := &Server{
		protocol:  protocol.NewProtocol(nil),
		transport: transport,
		tools:     new(datastructures.SyncMap[string, *tool]),
		prompts:   new(datastructures.SyncMap[string, *prompt]),
		resources: new(datastructures.SyncMap[string, *resource]),
	}
	for _, option := range options {
		option(server)
	}
	return server
}

// RegisterTool registers a new tool with the server
func (s *Server) RegisterTool(name string, description string, handler any) error {
	err := validateToolHandler(handler)
	if err != nil {
		return err
	}
	inputSchema := createJsonSchemaFromHandler(handler)

	s.tools.Store(name, &tool{
		Name:            name,
		Description:     description,
		Handler:         createWrappedToolHandler(handler),
		ToolInputSchema: inputSchema,
	})

	return s.sendToolListChangedNotification()
}

func (s *Server) sendToolListChangedNotification() error {
	if !s.isRunning {
		return nil
	}
	return s.protocol.Notification("notifications/tools/list_changed", nil)
}

func (s *Server) CheckToolRegistered(name string) bool {
	_, ok := s.tools.Load(name)
	return ok
}

func (s *Server) DeregisterTool(name string) error {
	s.tools.Delete(name)
	return s.sendToolListChangedNotification()
}

func (s *Server) RegisterResource(uri string, name string, description string, mimeType string, handler any) error {
	err := validateResourceHandler(handler)
	if err != nil {
		panic(err)
	}
	s.resources.Store(uri, &resource{
		Name:        name,
		Description: description,
		Uri:         uri,
		mimeType:    mimeType,
		Handler:     createWrappedResourceHandler(handler),
	})
	return s.sendResourceListChangedNotification()
}

func (s *Server) sendResourceListChangedNotification() error {
	if !s.isRunning {
		return nil
	}
	return s.protocol.Notification("notifications/resources/list_changed", nil)
}

func (s *Server) CheckResourceRegistered(uri string) bool {
	_, ok := s.resources.Load(uri)
	return ok
}

func (s *Server) DeregisterResource(uri string) error {
	s.resources.Delete(uri)
	return s.sendResourceListChangedNotification()
}

func createWrappedResourceHandler(userHandler any) func() *resourceResponseSent {
	handlerValue := reflect.ValueOf(userHandler)
	return func() *resourceResponseSent {
		// Call the handler with no arguments
		output := handlerValue.Call([]reflect.Value{})

		if len(output) != 2 {
			return newResourceResponseSentError(fmt.Errorf("handler must return exactly two values, got %d", len(output)))
		}

		if !output[0].CanInterface() {
			return newResourceResponseSentError(fmt.Errorf("handler must return a struct, got %s", output[0].Type().Name()))
		}
		promptR := output[0].Interface()
		if !output[1].CanInterface() {
			return newResourceResponseSentError(fmt.Errorf("handler must return an error, got %s", output[1].Type().Name()))
		}
		errorOut := output[1].Interface()
		if errorOut == nil {
			return newResourceResponseSent(promptR.(*ResourceResponse))
		}
		return newResourceResponseSentError(errorOut.(error))
	}
}

// We just want to check that handler takes no arguments and returns a ResourceResponse and an error
func validateResourceHandler(handler any) error {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()
	if handlerType.NumIn() != 0 {
		return fmt.Errorf("handler must take no arguments, got %d", handlerType.NumIn())
	}
	if handlerType.NumOut() != 2 {
		return fmt.Errorf("handler must return exactly two values, got %d", handlerType.NumOut())
	}
	//if handlerType.Out(0) != reflect.TypeOf((*ResourceResponse)(nil)).Elem() {
	//	return fmt.Errorf("handler must return ResourceResponse, got %s", handlerType.Out(0).Name())
	//}
	//if handlerType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
	//	return fmt.Errorf("handler must return error, got %s", handlerType.Out(1).Name())
	//}
	return nil
}

func (s *Server) RegisterPrompt(name string, description string, handler any) error {
	err := validatePromptHandler(handler)
	if err != nil {
		return err
	}
	promptSchema := createPromptSchemaFromHandler(handler)
	s.prompts.Store(name, &prompt{
		Name:              name,
		Description:       description,
		Handler:           createWrappedPromptHandler(handler),
		PromptInputSchema: promptSchema,
	})

	return s.sendPromptListChangedNotification()
}

func (s *Server) sendPromptListChangedNotification() error {
	if !s.isRunning {
		return nil
	}
	return s.protocol.Notification("notifications/prompts/list_changed", nil)
}

func (s *Server) CheckPromptRegistered(name string) bool {
	_, ok := s.prompts.Load(name)
	return ok
}

func (s *Server) DeregisterPrompt(name string) error {
	s.prompts.Delete(name)
	return s.sendPromptListChangedNotification()
}

func createWrappedPromptHandler(userHandler any) func(baseGetPromptRequestParamsArguments) *promptResponseSent {
	handlerValue := reflect.ValueOf(userHandler)
	handlerType := handlerValue.Type()
	argumentType := handlerType.In(0)
	return func(arguments baseGetPromptRequestParamsArguments) *promptResponseSent {
		// Instantiate a struct of the type of the arguments
		if !reflect.New(argumentType).CanInterface() {
			return newPromptResponseSentError(fmt.Errorf("arguments must be a struct"))
		}
		unmarshaledArguments := reflect.New(argumentType).Interface()

		// Unmarshal the JSON into the correct type
		err := json.Unmarshal(arguments.Arguments, &unmarshaledArguments)
		if err != nil {
			return newPromptResponseSentError(fmt.Errorf("failed to unmarshal arguments: %w", err))
		}

		// Need to dereference the unmarshaled arguments
		of := reflect.ValueOf(unmarshaledArguments)
		if of.Kind() != reflect.Ptr || !of.Elem().CanInterface() {
			return newPromptResponseSentError(fmt.Errorf("arguments must be a struct"))
		}
		// Call the handler with the typed arguments
		output := handlerValue.Call([]reflect.Value{of.Elem()})

		if len(output) != 2 {
			return newPromptResponseSentError(fmt.Errorf("handler must return exactly two values, got %d", len(output)))
		}

		if !output[0].CanInterface() {
			return newPromptResponseSentError(fmt.Errorf("handler must return a struct, got %s", output[0].Type().Name()))
		}
		promptR := output[0].Interface()
		if !output[1].CanInterface() {
			return newPromptResponseSentError(fmt.Errorf("handler must return an error, got %s", output[1].Type().Name()))
		}
		errorOut := output[1].Interface()
		if errorOut == nil {
			return newPromptResponseSent(promptR.(*PromptResponse))
		}
		return newPromptResponseSentError(errorOut.(error))
	}
}

// Get the argument and iterate over the fields, we pull description from the jsonschema description tag
// We pull required from the jsonschema required tag
// Example:
// type Content struct {
// Title       string  `json:"title" jsonschema:"description=The title to submit,required"`
// Description *string `json:"description" jsonschema:"description=The description to submit"`
// }
// Then we get the jsonschema for the struct where Title is a required field and Description is an optional field
func createPromptSchemaFromHandler(handler any) *promptSchema {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()
	argumentType := handlerType.In(0)

	promptSchema := promptSchema{
		Arguments: make([]promptSchemaArgument, argumentType.NumField()),
	}

	for i := 0; i < argumentType.NumField(); i++ {
		field := argumentType.Field(i)
		fieldName := field.Name

		jsonSchemaTags := strings.Split(field.Tag.Get("jsonschema"), ",")
		var description *string
		var required = false
		for _, tag := range jsonSchemaTags {
			if strings.HasPrefix(tag, "description=") {
				s := strings.TrimPrefix(tag, "description=")
				description = &s
			}
			if tag == "required" {
				required = true
			}
		}

		promptSchema.Arguments[i] = promptSchemaArgument{
			Name:        fieldName,
			Description: description,
			Required:    &required,
		}
	}
	return &promptSchema
}

// A prompt can only take a struct with fields of type string or *string as the argument
func validatePromptHandler(handler any) error {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()
	argumentType := handlerType.In(0)

	if argumentType.Kind() != reflect.Struct {
		return fmt.Errorf("argument must be a struct")
	}

	for i := 0; i < argumentType.NumField(); i++ {
		field := argumentType.Field(i)
		isValid := false
		if field.Type.Kind() == reflect.String {
			isValid = true
		}
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.String {
			isValid = true
		}
		if !isValid {
			return fmt.Errorf("all fields of the struct must be of type string or *string, found %s", field.Type.Kind())
		}
	}
	return nil
}

// Creates a full JSON schema from a user provided handler by introspecting the arguments
func createJsonSchemaFromHandler(handler any) *jsonschema.Schema {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()
	argumentType := handlerType.In(0)
	inputSchema := jsonSchemaReflector.ReflectFromType(argumentType)
	return inputSchema
}

// This takes a user provided handler and returns a wrapped handler which can be used to actually answer requests
// Concretely, it will deserialize the arguments and call the user provided handler and then serialize the response
// If the handler returns an error, it will be serialized and sent back as a tool error rather than a protocol error
func createWrappedToolHandler(userHandler any) func(baseCallToolRequestParams) *toolResponseSent {
	handlerValue := reflect.ValueOf(userHandler)
	handlerType := handlerValue.Type()
	argumentType := handlerType.In(0)
	return func(arguments baseCallToolRequestParams) *toolResponseSent {
		// Instantiate a struct of the type of the arguments
		if !reflect.New(argumentType).CanInterface() {
			return newToolResponseSentError(fmt.Errorf("arguments must be a struct"))
		}
		unmarshaledArguments := reflect.New(argumentType).Interface()

		// Unmarshal the JSON into the correct type
		err := json.Unmarshal(arguments.Arguments, &unmarshaledArguments)
		if err != nil {
			return newToolResponseSentError(fmt.Errorf("failed to unmarshal arguments: %w", err))
		}

		// Need to dereference the unmarshaled arguments
		of := reflect.ValueOf(unmarshaledArguments)
		if of.Kind() != reflect.Ptr || !of.Elem().CanInterface() {
			return newToolResponseSentError(fmt.Errorf("arguments must be a struct"))
		}
		// Call the handler with the typed arguments
		output := handlerValue.Call([]reflect.Value{of.Elem()})

		if len(output) != 2 {
			return newToolResponseSentError(fmt.Errorf("handler must return exactly two values, got %d", len(output)))
		}

		if !output[0].CanInterface() {
			return newToolResponseSentError(fmt.Errorf("handler must return a struct, got %s", output[0].Type().Name()))
		}
		tool := output[0].Interface()
		if !output[1].CanInterface() {
			return newToolResponseSentError(fmt.Errorf("handler must return an error, got %s", output[1].Type().Name()))
		}
		errorOut := output[1].Interface()
		if errorOut == nil {
			return newToolResponseSent(tool.(*ToolResponse))
		}
		return newToolResponseSentError(errorOut.(error))
	}
}

func (s *Server) Serve() error {
	if s.isRunning == true {
		return fmt.Errorf("server is already running")
	}
	pr := s.protocol
	pr.SetRequestHandler("ping", s.handlePing)
	pr.SetRequestHandler("initialize", s.handleInitialize)
	pr.SetRequestHandler("tools/list", s.handleListTools)
	pr.SetRequestHandler("tools/call", s.handleToolCalls)
	pr.SetRequestHandler("prompts/list", s.handleListPrompts)
	pr.SetRequestHandler("prompts/get", s.handlePromptCalls)
	pr.SetRequestHandler("resources/list", s.handleListResources)
	pr.SetRequestHandler("resources/read", s.handleResourceCalls)
	err := pr.Connect(s.transport)
	if err != nil {
		return err
	}
	s.protocol = pr
	s.isRunning = true
	return nil
}

func (s *Server) handleInitialize(_ *transport.BaseJSONRPCRequest, _ protocol.RequestHandlerExtra) (transport.JsonRpcBody, error) {
	return initializeResult{
		Meta:            nil,
		Capabilities:    s.generateCapabilities(),
		Instructions:    s.serverInstructions,
		ProtocolVersion: "2024-11-05",
		ServerInfo: implementation{
			Name:    s.serverName,
			Version: s.serverVersion,
		},
	}, nil
}

func (s *Server) handleListTools(request *transport.BaseJSONRPCRequest, _ protocol.RequestHandlerExtra) (transport.JsonRpcBody, error) {
	type toolRequestParams struct {
		Cursor *string `json:"cursor"`
	}
	var params toolRequestParams
	err := json.Unmarshal(request.Params, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	// Order by name for pagination
	var orderedTools []*tool
	s.tools.Range(func(k string, t *tool) bool {
		orderedTools = append(orderedTools, t)
		return true
	})
	sort.Slice(orderedTools, func(i, j int) bool {
		return orderedTools[i].Name < orderedTools[j].Name
	})

	startPosition := 0
	if params.Cursor != nil {
		// Base64 decode the cursor
		c, err := base64.StdEncoding.DecodeString(*params.Cursor)
		if err != nil {
			return nil, fmt.Errorf("failed to decode cursor: %w", err)
		}
		cString := string(c)
		// Iterate through the tools until we find an entry > the cursor
		found := false
		for i := 0; i < len(orderedTools); i++ {
			if orderedTools[i].Name > cString {
				startPosition = i
				found = true
				break
			}
		}
		if !found {
			startPosition = len(orderedTools)
		}
	}
	endPosition := len(orderedTools)
	if s.paginationLimit != nil {
		// Make sure we don't go out of bounds
		if len(orderedTools) > startPosition+*s.paginationLimit {
			endPosition = startPosition + *s.paginationLimit
		}
	}

	toolsToReturn := make([]tools.ToolRetType, 0)

	for i := startPosition; i < endPosition; i++ {
		toolsToReturn = append(toolsToReturn, tools.ToolRetType{
			Name:        orderedTools[i].Name,
			Description: &orderedTools[i].Description,
			InputSchema: orderedTools[i].ToolInputSchema,
		})
	}

	return tools.ToolsResponse{
		Tools: toolsToReturn,
		NextCursor: func() *string {
			if s.paginationLimit != nil && len(toolsToReturn) >= *s.paginationLimit {
				toString := base64.StdEncoding.EncodeToString([]byte(toolsToReturn[len(toolsToReturn)-1].Name))
				return &toString
			}
			return nil
		}(),
	}, nil
}

func (s *Server) handleToolCalls(req *transport.BaseJSONRPCRequest, _ protocol.RequestHandlerExtra) (transport.JsonRpcBody, error) {
	params := baseCallToolRequestParams{}
	// Instantiate a struct of the type of the arguments
	err := json.Unmarshal(req.Params, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	var toolToUse *tool
	s.tools.Range(func(k string, t *tool) bool {
		if k != params.Name {
			return true
		}
		toolToUse = t
		return false
	})

	if toolToUse == nil {
		return nil, fmt.Errorf("unknown tool: %s", req.Method)
	}
	return toolToUse.Handler(params), nil
}

func (s *Server) generateCapabilities() serverCapabilities {
	t := false
	return serverCapabilities{
		Tools: func() *serverCapabilitiesTools {
			return &serverCapabilitiesTools{
				ListChanged: &t,
			}
		}(),
		Prompts: func() *serverCapabilitiesPrompts {
			return &serverCapabilitiesPrompts{
				ListChanged: &t,
			}
		}(),
		Resources: func() *serverCapabilitiesResources {
			return &serverCapabilitiesResources{
				ListChanged: &t,
			}
		}(),
	}
}

func (s *Server) handleListPrompts(request *transport.BaseJSONRPCRequest, extra protocol.RequestHandlerExtra) (transport.JsonRpcBody, error) {
	type promptRequestParams struct {
		Cursor *string `json:"cursor"`
	}
	var params promptRequestParams
	err := json.Unmarshal(request.Params, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	// Order by name for pagination
	var orderedPrompts []*prompt
	s.prompts.Range(func(k string, p *prompt) bool {
		orderedPrompts = append(orderedPrompts, p)
		return true
	})
	sort.Slice(orderedPrompts, func(i, j int) bool {
		return orderedPrompts[i].Name < orderedPrompts[j].Name
	})

	startPosition := 0
	if params.Cursor != nil {
		// Base64 decode the cursor
		c, err := base64.StdEncoding.DecodeString(*params.Cursor)
		if err != nil {
			return nil, fmt.Errorf("failed to decode cursor: %w", err)
		}
		cString := string(c)
		// Iterate through the prompts until we find an entry > the cursor
		for i := 0; i < len(orderedPrompts); i++ {
			if orderedPrompts[i].Name > cString {
				startPosition = i
				break
			}
		}
	}
	endPosition := len(orderedPrompts)
	if s.paginationLimit != nil {
		// Make sure we don't go out of bounds
		if len(orderedPrompts) > startPosition+*s.paginationLimit {
			endPosition = startPosition + *s.paginationLimit
		}
	}

	promptsToReturn := make([]*promptSchema, 0)
	for i := startPosition; i < endPosition; i++ {
		schema := orderedPrompts[i].PromptInputSchema
		schema.Name = orderedPrompts[i].Name
		promptsToReturn = append(promptsToReturn, schema)
	}

	return listPromptsResult{
		Prompts: promptsToReturn,
		NextCursor: func() *string {
			if s.paginationLimit != nil && len(promptsToReturn) >= *s.paginationLimit {
				toString := base64.StdEncoding.EncodeToString([]byte(promptsToReturn[len(promptsToReturn)-1].Name))
				return &toString
			}
			return nil
		}(),
	}, nil
}

func (s *Server) handleListResources(request *transport.BaseJSONRPCRequest, extra protocol.RequestHandlerExtra) (transport.JsonRpcBody, error) {
	type resourceRequestParams struct {
		Cursor *string `json:"cursor"`
	}
	var params resourceRequestParams
	err := json.Unmarshal(request.Params, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	// Order by URI for pagination
	var orderedResources []*resource
	s.resources.Range(func(k string, r *resource) bool {
		orderedResources = append(orderedResources, r)
		return true
	})
	sort.Slice(orderedResources, func(i, j int) bool {
		return orderedResources[i].Uri < orderedResources[j].Uri
	})

	startPosition := 0
	if params.Cursor != nil {
		// Base64 decode the cursor
		c, err := base64.StdEncoding.DecodeString(*params.Cursor)
		if err != nil {
			return nil, fmt.Errorf("failed to decode cursor: %w", err)
		}
		cString := string(c)
		// Iterate through the resources until we find an entry > the cursor
		for i := 0; i < len(orderedResources); i++ {
			if orderedResources[i].Uri > cString {
				startPosition = i
				break
			}
		}
	}
	endPosition := len(orderedResources)
	if s.paginationLimit != nil {
		// Make sure we don't go out of bounds
		if len(orderedResources) > startPosition+*s.paginationLimit {
			endPosition = startPosition + *s.paginationLimit
		}
	}

	resourcesToReturn := make([]*resourceSchema, 0)
	for i := startPosition; i < endPosition; i++ {
		r := orderedResources[i]
		resourcesToReturn = append(resourcesToReturn, &resourceSchema{
			Annotations: nil,
			Description: &r.Description,
			MimeType:    &r.mimeType,
			Name:        r.Name,
			Uri:         r.Uri,
		})
	}

	return listResourcesResult{
		Resources: resourcesToReturn,
		NextCursor: func() *string {
			if s.paginationLimit != nil && len(resourcesToReturn) >= *s.paginationLimit {
				toString := base64.StdEncoding.EncodeToString([]byte(resourcesToReturn[len(resourcesToReturn)-1].Uri))
				return &toString
			}
			return nil
		}(),
	}, nil
}

func (s *Server) handlePromptCalls(req *transport.BaseJSONRPCRequest, extra protocol.RequestHandlerExtra) (transport.JsonRpcBody, error) {
	params := baseGetPromptRequestParamsArguments{}
	// Instantiate a struct of the type of the arguments
	err := json.Unmarshal(req.Params, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	var promptToUse *prompt
	s.prompts.Range(func(k string, p *prompt) bool {
		if k != params.Name {
			return true
		}
		promptToUse = p
		return false
	})

	if promptToUse == nil {
		return nil, fmt.Errorf("unknown prompt: %s", req.Method)
	}
	return promptToUse.Handler(params), nil
}

func (s *Server) handleResourceCalls(req *transport.BaseJSONRPCRequest, extra protocol.RequestHandlerExtra) (transport.JsonRpcBody, error) {
	params := readResourceRequestParams{}
	// Instantiate a struct of the type of the arguments
	err := json.Unmarshal(req.Params, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	var resourceToUse *resource
	s.resources.Range(func(k string, r *resource) bool {
		if k != params.Uri {
			return true
		}
		resourceToUse = r
		return false
	})

	if resourceToUse == nil {
		return nil, fmt.Errorf("unknown prompt: %s", req.Method)
	}
	return resourceToUse.Handler(), nil
}

func (s *Server) handlePing(request *transport.BaseJSONRPCRequest, extra protocol.RequestHandlerExtra) (transport.JsonRpcBody, error) {
	return map[string]interface{}{}, nil
}

func validateToolHandler(handler any) error {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()

	if handlerType.NumIn() != 1 {
		return fmt.Errorf("handler must take exactly one argument, got %d", handlerType.NumIn())
	}

	if handlerType.NumOut() != 2 {
		return fmt.Errorf("handler must return exactly two values, got %d", handlerType.NumOut())
	}

	// Check that the output type is *tools.ToolResponse
	if handlerType.Out(0) != reflect.PointerTo(reflect.TypeOf(ToolResponse{})) {
		return fmt.Errorf("handler must return *tools.ToolResponse, got %s", handlerType.Out(0).Name())
	}

	// Check that the output type is error
	if handlerType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return fmt.Errorf("handler must return error, got %s", handlerType.Out(1).Name())
	}

	return nil
}

var (
	jsonSchemaReflector = jsonschema.Reflector{
		BaseSchemaID:               "",
		Anonymous:                  true,
		AssignAnchor:               false,
		AllowAdditionalProperties:  true,
		RequiredFromJSONSchemaTags: true,
		DoNotReference:             true,
		ExpandedStruct:             true,
		FieldNameTag:               "",
		IgnoredTypes:               nil,
		Lookup:                     nil,
		Mapper:                     nil,
		Namer:                      nil,
		KeyNamer:                   nil,
		AdditionalFields:           nil,
		CommentMap:                 nil,
	}
)
