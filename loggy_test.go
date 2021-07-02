package loggy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestExtractArgs(t *testing.T) {
	tests := map[string]struct {
		givenFields  []string
		givenCtx     context.Context
		expectedArgs []interface{}
	}{
		"Should return empty slice": {
			givenCtx: context.Background(),
		},
		"Should return request_id field and value": {
			givenFields:  []string{"request_id"},
			givenCtx:     context.WithValue(context.Background(), "request_id", "<request-id-value>"),
			expectedArgs: []interface{}{"request_id", "<request-id-value>"},
		},
		"Should return not return request_id field and value if not present": {
			givenFields: []string{"request_id"},
			givenCtx:    context.Background(),
		},
		"Should return request_id field and value and user field and value": {
			givenFields:  []string{"request_id", "user"},
			givenCtx:     context.WithValue(context.WithValue(context.Background(), "request_id", "<request-id-value>"), "user", "<user-value>"),
			expectedArgs: []interface{}{"request_id", "<request-id-value>", "user", "<user-value>"},
		},
		"Should only return user field and value if request_id not present": {
			givenFields:  []string{"request_id", "user"},
			givenCtx:     context.WithValue(context.Background(), "user", "<user-value>"),
			expectedArgs: []interface{}{"user", "<user-value>"},
		},
	}

	for n, tc := range tests {
		tc := tc
		t.Run(n, func(t *testing.T) {
			zaplogger, err := zap.NewProduction()
			require.NoError(t, err)

			sugaredLogger := zaplogger.Sugar()

			logger := New(sugaredLogger)
			logger = logger.WithFields(tc.givenFields...)

			actualArgs := logger.extractArgs(tc.givenCtx)
			require.ElementsMatch(t, tc.expectedArgs, actualArgs)
		})
	}
}
