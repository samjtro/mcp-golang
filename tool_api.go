package mcp_golang

import (
	"encoding/json"
	"fmt"
)

// Definition for a tool the client can call.
type Tool struct {
	// A human-readable description of the tool.
	Description *string `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// A JSON Schema object defining the expected parameters for the tool.
	InputSchema ToolInputSchema `json:"inputSchema" yaml:"inputSchema" mapstructure:"inputSchema"`

	// The name of the tool.
	Name string `json:"name" yaml:"name" mapstructure:"name"`
}

// A JSON Schema object defining the expected parameters for the tool.
type ToolInputSchema struct {
	// Properties corresponds to the JSON schema field "properties".
	Properties ToolInputSchemaProperties `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// Required corresponds to the JSON schema field "required".
	Required []string `json:"required,omitempty" yaml:"required,omitempty" mapstructure:"required,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}

type ToolInputSchemaProperties map[string]map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ToolInputSchema) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in ToolInputSchema: required")
	}
	type Plain ToolInputSchema
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ToolInputSchema(plain)
	return nil
}

// An optional notification from the server to the client, informing it that the
// list of tools it offers has changed. This may be issued by servers without any
// previous subscription from the client.
type ToolListChangedNotification struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"`

	// Params corresponds to the JSON schema field "params".
	Params *ToolListChangedNotificationParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`
}

type ToolListChangedNotificationParams struct {
	// This parameter name is reserved by MCP to allow clients and servers to attach
	// additional metadata to their notifications.
	Meta ToolListChangedNotificationParamsMeta `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	AdditionalProperties interface{} `mapstructure:",remain"`
}

// This parameter name is reserved by MCP to allow clients and servers to attach
// additional metadata to their notifications.
type ToolListChangedNotificationParamsMeta map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ToolListChangedNotification) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["method"]; raw != nil && !ok {
		return fmt.Errorf("field method in ToolListChangedNotification: required")
	}
	type Plain ToolListChangedNotification
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ToolListChangedNotification(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Tool) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if _, ok := raw["inputSchema"]; raw != nil && !ok {
		return fmt.Errorf("field inputSchema in Tool: required")
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in Tool: required")
	}
	type Plain Tool
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Tool(plain)
	return nil
}

// This is a union type of all the different ToolResponse that can be sent back to the client.
// We allow creation through constructors only to make sure that the ToolResponse is valid.
type ToolResponse struct {
	Content []*Content
}

func NewToolResponse(content ...*Content) *ToolResponse {
	return &ToolResponse{
		Content: content,
	}
}
