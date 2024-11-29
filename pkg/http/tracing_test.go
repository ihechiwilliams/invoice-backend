package http

import (
	"bytes"
	"context"
	"errors"
	"invoice-backend/pkg/http/mocks"
	"io"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockRoundTripper struct {
	mock.Mock
}

func (m *mockRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	args := m.Called(request)

	return args.Get(0).(*http.Response), args.Error(1)
}

type loggingTransportSuite struct {
	suite.Suite
	context      context.Context
	logger       *zerolog.Logger
	out          *bytes.Buffer
	roundTripper *mockRoundTripper
	transport    *Transport
	metrics      *mocks.StatsdInterface
}

func (s *loggingTransportSuite) SetupTest() {
	s.roundTripper = &mockRoundTripper{}
	s.out = bytes.NewBufferString("")
	l := zerolog.New(s.out).With().Str("id", "XYZ1337").Logger()
	s.logger = &l
	s.context = l.WithContext(context.Background())
	s.metrics = mocks.NewStatsdInterface(s.T())

	s.transport = NewTransport(s.roundTripper, s.metrics, WithServiceName("abc-service"))
}

func (s *loggingTransportSuite) TearDownTest() {
	s.roundTripper.AssertExpectations(s.T())
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(loggingTransportSuite))
}

func (s *loggingTransportSuite) TestLoggingTransport_ValidJSON() {
	req, _ := http.NewRequestWithContext(s.context, http.MethodGet, "https://example.com/xyz?something=xyt", http.NoBody)
	resp := &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{ "foo": "bar" }`)),
	}

	s.metrics.On("Incr", providerRequestHit, mock.Anything, mock.Anything).Return(nil)
	s.metrics.On("Histogram", providerRequestDuration, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.roundTripper.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(resp, nil)

	resp, err := s.transport.RoundTrip(req)

	respBody, _ := io.ReadAll(resp.Body)

	s.NoError(err)
	s.Equal("OK", resp.Status)
	s.Equal(200, resp.StatusCode)
	s.Equal(`{ "foo": "bar" }`, string(respBody))
	s.Contains(s.out.String(), `"host":"example.com"`)
	s.Contains(s.out.String(), `"path":"/xyz"`)
	s.Contains(s.out.String(), `"level":"debug"`)
	s.Contains(s.out.String(), `"service":"abc-service"`)
	s.Contains(s.out.String(), `"uri":"/xyz?something=xyt"`)
	s.Contains(s.out.String(), `"id":"XYZ1337"`)
	s.Contains(s.out.String(), `"status":"OK"`)
	s.Contains(s.out.String(), `"status_code":200`)
}

func (s *loggingTransportSuite) TestLoggingTransport_ErroneousResponse() {
	req, _ := http.NewRequestWithContext(s.context, http.MethodPost, "https://example.com/xyz?something=xyt", http.NoBody)
	resp := &http.Response{
		Status:     http.StatusText(http.StatusBadRequest),
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewBufferString(`{ "foo": "bar" }`)),
	}

	s.metrics.On("Incr", providerRequestHit, mock.Anything, mock.Anything).Return(nil)
	s.metrics.On("Histogram", providerRequestDuration, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.roundTripper.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(resp, nil)

	resp, err := s.transport.RoundTrip(req)

	respBody, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	s.NoError(err)
	s.Equal("Bad Request", resp.Status)
	s.Equal(400, resp.StatusCode)
	s.Equal(`{ "foo": "bar" }`, string(respBody))
	s.Contains(s.out.String(), `"host":"example.com"`)
	s.Contains(s.out.String(), `"path":"/xyz"`)
	s.Contains(s.out.String(), `"level":"debug"`)
	s.Contains(s.out.String(), `"service":"abc-service"`)
	s.Contains(s.out.String(), `"response_body":{"foo":"bar"}`)
	s.Contains(s.out.String(), `"uri":"/xyz?something=xyt"`)
	s.Contains(s.out.String(), `"id":"XYZ1337"`)
	s.Contains(s.out.String(), `"status":"Bad Request"`)
	s.Contains(s.out.String(), `"status_code":400`)
}

func (s *loggingTransportSuite) TestLoggingTransport_InvalidJSON() {
	req, _ := http.NewRequestWithContext(s.context, http.MethodGet, "https://example.com/xyz?something=qwe", http.NoBody)
	resp := &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`NOT JSON`)),
	}

	s.metrics.On("Incr", providerRequestHit, mock.Anything, mock.Anything).Return(nil)
	s.metrics.On("Histogram", providerRequestDuration, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.roundTripper.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(resp, nil)

	resp, err := s.transport.RoundTrip(req)

	respBody, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	s.NoError(err)
	s.Equal("OK", resp.Status)
	s.Equal(200, resp.StatusCode)
	s.Equal(`NOT JSON`, string(respBody))

	s.Contains(s.out.String(), `"host":"example.com"`)
	s.Contains(s.out.String(), `"path":"/xyz"`)
	s.Contains(s.out.String(), `"level":"debug"`)
	s.Contains(s.out.String(), `"service":"abc-service"`)
	s.Contains(s.out.String(), `"uri":"/xyz?something=qwe"`)
	s.Contains(s.out.String(), `"id":"XYZ1337"`)
}

func (s *loggingTransportSuite) TestLoggingTransport_RoundTripErr() {
	req, _ := http.NewRequestWithContext(s.context, http.MethodGet, "https://example.com/xyz?something=qwe", http.NoBody)

	var retResp *http.Response

	s.roundTripper.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(retResp, errors.New("error"))
	s.metrics.On("Incr", providerRequestHit, mock.Anything, mock.Anything).Return(nil)
	s.metrics.On("Histogram", providerRequestDuration, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	resp, err := s.transport.RoundTrip(req)

	s.EqualError(err, "error")
	s.Nil(resp)

	s.Contains(s.out.String(), `"host":"example.com"`)
	s.Contains(s.out.String(), `"path":"/xyz"`)
	s.Contains(s.out.String(), `"level":"debug"`)
	s.Contains(s.out.String(), `"service":"abc-service"`)
	s.Contains(s.out.String(), `"uri":"/xyz?something=qwe"`)
	s.Contains(s.out.String(), `"id":"XYZ1337"`)
	s.Contains(s.out.String(), `"http.RoundTripper error`)
}

func Test_sanitizePath(t *testing.T) {
	type args struct {
		value string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "uuid",
			args: args{
				value: "/v1/accounts/1fa47c64-99f1-42f1-b346-2c9c6584daf8",
			},
			want: "/v1/accounts/:id",
		},
		{
			name: "numbers",
			args: args{
				value: "/v1/accounts/9081209",
			},
			want: "/v1/accounts/:id",
		},
		{
			name: "more numbers",
			args: args{
				value: "/internal/kyc/status/7532",
			},
			want: "/internal/kyc/status/:id",
		},
		{
			name: "mixed",
			args: args{
				value: "/v1/accounts/41293812937/path/8283a925-b1ed-45be-9f89-cc0a59b68f5c",
			},
			want: "/v1/accounts/:id/path/:id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, sanitizePath(tt.args.value), "sanitizePath(%v)", tt.args.value)
		})
	}
}
