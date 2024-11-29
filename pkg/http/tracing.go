package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

const (
	uuidPattern             = `[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`
	numberPattern           = `\d+`
	urlPathSeparator        = "/"
	urlVersionPrefix        = "v"
	pathID                  = ":id"
	statusCodeDivisor       = 100
	providerRequestHit      = "provider.request.hit"
	providerRequestDuration = "provider.request.duration"
)

type rawResponseBody struct {
	RawBody string `json:"raw_body"`
}

type Transport struct {
	metrics      statsd.ClientInterface
	roundTripper http.RoundTripper
	jsonRedactor *JSONRedactor

	debugModeEnabled bool
	filteredKeys     []string
	serviceName      string
	providerName     string
}

type TransportOptions func(*Transport)

func WithDebugMode(debugModeEnabled bool) TransportOptions {
	return func(transport *Transport) {
		transport.debugModeEnabled = debugModeEnabled
	}
}

func WithFilteredKeys(filteredKeys []string) TransportOptions {
	return func(transport *Transport) {
		transport.filteredKeys = filteredKeys
	}
}

func WithProviderName(providerName string) TransportOptions {
	return func(transport *Transport) {
		transport.providerName = providerName
	}
}

func WithServiceName(serviceName string) TransportOptions {
	return func(transport *Transport) {
		transport.serviceName = serviceName
	}
}

func NewTransport(
	roundTripper http.RoundTripper,
	metrics statsd.ClientInterface,
	opts ...TransportOptions,
) *Transport {
	transport := &Transport{
		metrics: metrics,
	}

	for _, opt := range opts {
		opt(transport)
	}

	transport.jsonRedactor = NewJSONRedactor(WithKeysToHide(transport.filteredKeys), WithFilterString(defaultFilterString))
	transport.roundTripper = httptrace.WrapRoundTripper(
		roundTripper,
		httptrace.RTWithServiceName(transport.serviceName),
		httptrace.RTWithResourceNamer(func(req *http.Request) string { return sanitizePath(req.URL.Path) }),
	)

	return transport
}

func (trp *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	defer func() {
		if req.Body != nil {
			_ = req.Body.Close()
		}
	}()

	startedAt := time.Now().UTC()

	reqLog := log.Ctx(req.Context()).With().
		Str("host", req.URL.Host).
		Str("method", req.Method).
		Str("path", req.URL.Path).
		Str("service", trp.serviceName).
		Str("uri", req.URL.RequestURI()).
		Logger()

	if trp.isRequestEligibleForLogging(req) {
		reqBody, respErr := io.ReadAll(req.Body)

		if respErr != nil {
			reqLog.Error().Err(respErr).Msg("cannot read request body")

			return nil, respErr
		}

		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		requestJSON := compactJSON(reqBody)
		requestJSON = trp.jsonRedactor.HideJSONKeys(requestJSON)

		reqLog = reqLog.With().RawJSON("request_body", requestJSON).Logger()
	}

	reqLog.Debug().Msgf("started HTTP request %s", req.URL.RequestURI())

	responseLogLevel := zerolog.DebugLevel

	var resp *http.Response
	defer func() {
		finishedAt := time.Now().UTC()

		trp.sendProviderMetrics(req, resp, finishedAt.Sub(startedAt))

		reqLog.WithLevel(responseLogLevel).
			TimeDiff("duration", finishedAt, startedAt).
			Msgf("finished HTTP request %s", req.URL.RequestURI())
	}()

	resp, err := trp.roundTripper.RoundTrip(req)
	if err != nil {
		reqLog.Error().Err(err).Msg("http.RoundTripper error")

		return resp, err
	}

	if resp != nil {
		reqLog = reqLog.With().
			Str("status", resp.Status).
			Int("status_code", resp.StatusCode).
			Logger()
	}

	if trp.isResponseEligibleForLogging(resp) {
		if resp.StatusCode >= http.StatusBadRequest {
			responseLogLevel = zerolog.ErrorLevel
		}

		respBody, respErr := io.ReadAll(resp.Body)

		if respErr != nil {
			reqLog.Error().Err(respErr).Msg("cannot read response body")

			return resp, respErr
		}

		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))

		responseJSON := compactJSON(respBody)
		responseJSON = trp.jsonRedactor.HideJSONKeys(responseJSON)

		reqLog = reqLog.With().RawJSON("response_body", responseJSON).Logger()
	}

	return resp, err
}

func (trp *Transport) isRequestEligibleForLogging(req *http.Request) bool {
	if req == nil {
		return false
	}

	if req.Body == nil {
		return false
	}

	if req.Method == http.MethodGet {
		return false
	}

	return trp.debugModeEnabled
}

func (trp *Transport) isResponseEligibleForLogging(resp *http.Response) bool {
	if resp == nil {
		return false
	}

	if resp.Body == nil {
		return false
	}

	isError := resp.StatusCode >= http.StatusBadRequest

	return trp.debugModeEnabled || isError
}

func (trp *Transport) sendProviderMetrics(
	req *http.Request,
	resp *http.Response,
	duration time.Duration,
) {
	endpoint := fmt.Sprintf("%s %s", req.Method, sanitizePath(req.URL.Path))

	statusCode := http.StatusServiceUnavailable
	if resp != nil {
		statusCode = resp.StatusCode
	}

	tags := []string{
		fmt.Sprintf("endpoint:%s", endpoint),
		fmt.Sprintf("response_code:%d", statusCode),
		fmt.Sprintf("http.status_class:%s", stringHTTPStatus(statusCode)),
		fmt.Sprintf("provider.name:%s", trp.providerName),
	}

	_ = trp.metrics.Incr(providerRequestHit, tags, 1)
	_ = trp.metrics.Histogram(providerRequestDuration, duration.Seconds(), tags, 1)
}

func stringHTTPStatus(num int) string {
	prefix := num / statusCodeDivisor

	return fmt.Sprintf("%dxx", prefix)
}

func sanitizePath(value string) string {
	uuidRegex := regexp.MustCompile(uuidPattern)
	numberRegex := regexp.MustCompile(numberPattern)

	parts := strings.Split(value, urlPathSeparator)

	for i, part := range parts {
		switch {
		case part == "", strings.HasPrefix(part, urlVersionPrefix):
			continue
		case uuidRegex.MatchString(part):
			parts[i] = uuidRegex.ReplaceAllString(part, pathID)
		case numberRegex.MatchString(part):
			parts[i] = numberRegex.ReplaceAllString(part, pathID)
		default:
			continue
		}
	}

	return strings.Join(parts, urlPathSeparator)
}
