package goupnp

import (
	"context"
	"net/http"
	"testing"
)

func TestWithHTTPClient(t *testing.T) {
	// Override the default HTTP client for testing purposes, to make sure
	// that we're returning that value (if it was changed) and not
	// http.DefaultClient.
	oldClient := HTTPClientDefault
	t.Cleanup(func() {
		HTTPClientDefault = oldClient
	})
	HTTPClientDefault = &http.Client{}

	customClient := &http.Client{}
	lastClient := &http.Client{}

	ctx := t.Context()
	tests := []struct {
		name string
		ctx  context.Context
		want *http.Client
	}{
		{
			name: "no context value returns HTTPClientDefault",
			ctx:  ctx,
			want: HTTPClientDefault,
		},
		{
			name: "context with custom client returns that client",
			ctx:  WithHTTPClient(ctx, customClient),
			want: customClient,
		},
		{
			name: "context with nil client returns HTTPClientDefault",
			ctx:  WithHTTPClient(ctx, nil),
			want: HTTPClientDefault,
		},
		{
			name: "nested contexts preserve last client",
			ctx:  WithHTTPClient(WithHTTPClient(ctx, customClient), lastClient),
			want: lastClient,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := httpClient(tt.ctx)
			if got != tt.want {
				t.Errorf("httpClient() = %p, want %p", got, tt.want)
			}
		})
	}
}
