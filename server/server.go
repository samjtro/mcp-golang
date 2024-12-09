package server

import (
	"encoding/json"
	"fmt"
	"github.com/invopop/jsonschema"
	internaltools "github.com/metoro-io/mcp-golang/internal/tools"
	protocol2 "github.com/metoro-io/mcp-golang/protocol"
	"github.com/metoro-io/mcp-golang/server/tools"
	transport2 "github.com/metoro-io/mcp-golang/transport"
	"reflect"
)

// Here we define the actual MCP server that users will create and run
// A server can be passed a number of handlers to handle requests from clients
// Additionally it can be parametrized by a transport. This transport will be used to actually send and receive messages.
// So for example if the stdio transport is used, the server will read from stdin and write to stdout
// If the SSE transport is used, the server will send messages over an SSE connection and receive messages from HTTP POST requests.

// The interface that we're looking to support is something like [gin](https://github.com/gin-gonic/gin)s interface

type Server struct {
	transport          transport2.Transport
	tools              map[string]*tool
	serverInstructions *string
	serverName         string
	serverVersion      string
}

type tool struct {
	Name            string
	Description     string
	Handler         func(transport2.BaseCallToolRequestParams) *tools.ToolResponseSent
	ToolInputSchema *jsonschema.Schema
}

func NewServer(transport transport2.Transport) *Server {
	return &Server{
		transport: transport,
		tools:     make(map[string]*tool),
	}
}

// RegisterTool registers a new tool with the server
func (s *Server) RegisterTool(name string, description string, handler any) error {
	err := validateHandler(handler)
	if err != nil {
		return err
	}
	inputSchema := createJsonSchemaFromHandler(handler)

	s.tools[name] = &tool{
		Name:            name,
		Description:     description,
		Handler:         createWrappedToolHandler(handler),
		ToolInputSchema: inputSchema,
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
func createWrappedToolHandler(userHandler any) func(transport2.BaseCallToolRequestParams) *tools.ToolResponseSent {
	handlerValue := reflect.ValueOf(userHandler)
	handlerType := handlerValue.Type()
	argumentType := handlerType.In(0)
	return func(arguments transport2.BaseCallToolRequestParams) *tools.ToolResponseSent {
		// Instantiate a struct of the type of the arguments
		if !reflect.New(argumentType).CanInterface() {
			return tools.NewToolResponseSentError(fmt.Errorf("arguments must be a struct"))
		}
		unmarshaledArguments := reflect.New(argumentType).Interface()

		// Unmarshal the JSON into the correct type
		err := json.Unmarshal(arguments.Arguments, &unmarshaledArguments)
		if err != nil {
			return tools.NewToolResponseSentError(fmt.Errorf("failed to unmarshal arguments: %w", err))
		}

		// Need to dereference the unmarshaled arguments
		of := reflect.ValueOf(unmarshaledArguments)
		if of.Kind() != reflect.Ptr || !of.Elem().CanInterface() {
			return tools.NewToolResponseSentError(fmt.Errorf("arguments must be a struct"))
		}
		// Call the handler with the typed arguments
		output := handlerValue.Call([]reflect.Value{of.Elem()})

		if len(output) != 2 {
			return tools.NewToolResponseSentError(fmt.Errorf("handler must return exactly two values, got %d", len(output)))
		}

		if !output[0].CanInterface() {
			return tools.NewToolResponseSentError(fmt.Errorf("handler must return a struct, got %s", output[0].Type().Name()))
		}
		tool := output[0].Interface()
		if !output[1].CanInterface() {
			return tools.NewToolResponseSentError(fmt.Errorf("handler must return an error, got %s", output[1].Type().Name()))
		}
		errorOut := output[1].Interface()
		if errorOut == nil {
			return tools.NewToolResponseSent(tool.(*tools.ToolResponse))
		}
		return tools.NewToolResponseSentError(errorOut.(error))
	}
}

func (s *Server) Serve() error {
	protocol := protocol2.NewProtocol(nil)
	protocol.SetRequestHandler("initialize", s.handleInitialize)
	protocol.SetRequestHandler("tools/list", s.handleListTools)
	protocol.SetRequestHandler("tools/call", s.handleToolCalls)
	return protocol.Connect(s.transport)
}

func (s *Server) handleInitialize(_ *transport2.BaseJSONRPCRequest, _ protocol2.RequestHandlerExtra) (interface{}, error) {
	return InitializeResult{
		Meta:            nil,
		Capabilities:    s.generateCapabilities(),
		Instructions:    s.serverInstructions,
		ProtocolVersion: "2024-11-05",
		ServerInfo: Implementation{
			Name:    s.serverName,
			Version: s.serverVersion,
		},
	}, nil
}

func (s *Server) handleListTools(_ *transport2.BaseJSONRPCRequest, _ protocol2.RequestHandlerExtra) (interface{}, error) {
	return internaltools.ToolsResponse{
		Tools: func() []internaltools.ToolRetType {
			var tools []internaltools.ToolRetType
			for _, tool := range s.tools {
				tools = append(tools, internaltools.ToolRetType{
					Name:        tool.Name,
					Description: &tool.Description,
					InputSchema: tool.ToolInputSchema,
				})
			}
			return tools
		}(),
	}, nil
}

func (s *Server) handleToolCalls(req *transport2.BaseJSONRPCRequest, _ protocol2.RequestHandlerExtra) (interface{}, error) {
	params := transport2.BaseCallToolRequestParams{}
	// Instantiate a struct of the type of the arguments
	err := json.Unmarshal(req.Params, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	for name, tool := range s.tools {
		if name != params.Name {
			continue
		}
		return tool.Handler(params), nil
	}
	return nil, fmt.Errorf("unknown tool: %s", req.Method)
}

func (s *Server) generateCapabilities() ServerCapabilities {
	f := false
	return ServerCapabilities{
		Tools: func() *ServerCapabilitiesTools {
			if s.tools == nil {
				return nil
			}
			return &ServerCapabilitiesTools{
				ListChanged: &f,
			}
		}(),
	}
}

func validateHandler(handler any) error {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()

	if handlerType.NumIn() != 1 {
		return fmt.Errorf("handler must take exactly one argument, got %d", handlerType.NumIn())
	}

	if handlerType.NumOut() != 2 {
		return fmt.Errorf("handler must return exactly two values, got %d", handlerType.NumOut())
	}

	// Check that the output type is *tools.ToolResponse
	if handlerType.Out(0) != reflect.PointerTo(reflect.TypeOf(tools.ToolResponse{})) {
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
		RequiredFromJSONSchemaTags: false,
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