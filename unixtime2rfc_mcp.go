package unixtime2rfc

import (
	"context"
	_ "embed" // Required for go:embed
	"errors"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed description.txt
var toolDescription string

// UnixTimeInput defines the parameters for the "unixtime2rfc" tool.
// Pointers are used to distinguish between a field being absent and having a zero value (epoch).
type UnixTimeInput struct {
	UnixTime   *int64 `json:"unixtime,omitempty"`
	UnixTimeMs *int64 `json:"unixtimeMs,omitempty"`
	UnixTimeUs *int64 `json:"unixtimeUs,omitempty"`
	Layout     string `json:"layout,omitempty"`
}

// UnixTimeOutput defines the output of the "unixtime2rfc" tool.
type UnixTimeOutput struct {
	FormattedTime string `json:"formattedTime"`
}

// ErrNoTimestampProvided is returned when no valid timestamp (unixtime, unixtimeMs, or unixtimeUs) is provided.
var ErrNoTimestampProvided = errors.New("no timestamp provided")

// ProcessTimeInput selects the appropriate timestamp based on precedence and converts it to a time.Time object.
func ProcessTimeInput(input UnixTimeInput) (time.Time, error) {
	if input.UnixTime != nil {
		return UnixSecToTime(*input.UnixTime), nil
	}
	if input.UnixTimeMs != nil {
		return UnixMilliToTime(*input.UnixTimeMs), nil
	}
	if input.UnixTimeUs != nil {
		return UnixMicroToTime(*input.UnixTimeUs), nil
	}
	return time.Time{}, ErrNoTimestampProvided
}

// NewServer creates a new MCP server and returns it as an http.Handler.
func NewServer() (http.Handler, error) {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "unixtime2rfc",
		Version: "v1.3.0",
		Title:   "Unix Time to Formatted String Converter",
	}, nil)

	unixtimeTool := func(ctx context.Context, req *mcp.CallToolRequest, input UnixTimeInput) (
		*mcp.CallToolResult,
		UnixTimeOutput,
		error,
	) {
		timeVal, err := ProcessTimeInput(input)
		if err != nil {
			return nil, UnixTimeOutput{}, err
		}

		formattedString, err := FormatTime(timeVal, input.Layout)
		if err != nil {
			return nil, UnixTimeOutput{}, err
		}

		return nil, UnixTimeOutput{FormattedTime: formattedString}, nil
	}

	//nolint:exhaustruct
	mcp.AddTool(server, &mcp.Tool{
		Name:        "unixtime2formatted",
		Title:       "Convert Unix Time",
		Description: toolDescription,
	}, unixtimeTool)

	handler := mcp.NewStreamableHTTPHandler(
		func(req *http.Request) *mcp.Server { return server },
		//nolint:exhaustruct
		&mcp.StreamableHTTPOptions{Stateless: true},
	)

	return handler, nil
}
