//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/eve-an/estimated/internal/api/dto"
	"github.com/eve-an/estimated/internal/config"
	"github.com/stretchr/testify/assert"
)

func buildJsonVotes(t *testing.T, value int, timestamp string) string {
	t.Helper()
	return fmt.Sprintf(`{"value": %d, "timestamp": "%s"}`, value, timestamp)
}

func addVote(t *testing.T, url string, value int, timestamp string) *http.Response {
	t.Helper()

	r := strings.NewReader(buildJsonVotes(t, value, timestamp))
	response, err := http.Post(url, "application/json", r)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	return response
}

func TestSuccessfullAddAndGet(t *testing.T) {
	cfg := &config.Config{
		SessionChannelSize: 1,
		ServerTimeout:      config.Duration(time.Second * 5),
	}

	handler, err := InitializeTestApp(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	votesURL := ts.URL + "/api/v1/votes"

	addVote(t, votesURL, 1, "2025-05-04T11:33:49.554689+02:00")
	addVote(t, votesURL, 2, "2025-05-04T11:33:50.553932+02:00")
	addVote(t, votesURL, 3, "2025-05-04T11:33:51.554181+02:00")
	addVote(t, votesURL, 5, "2025-05-04T11:33:52.554181+02:00")
	addVote(t, votesURL, 8, "2025-05-04T11:33:53.554181+02:00")

	resp, err := http.Get(ts.URL + "/api/v1/votes")
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, string(body))

	apiResponse := struct {
		Data []dto.VoteResponseDTO `json:"data"`
	}{}

	err = json.Unmarshal(body, &apiResponse)
	assert.NoError(t, err)

	slices.SortFunc(apiResponse.Data, func(a, b dto.VoteResponseDTO) int {
		return a.Timestamp.Compare(b.Timestamp)
	})

	assert.Equal(t, []dto.VoteResponseDTO{
		{Value: 1, Timestamp: mustParseTime(t, "2025-05-04T11:33:49.554689+02:00")},
		{Value: 2, Timestamp: mustParseTime(t, "2025-05-04T11:33:50.553932+02:00")},
		{Value: 3, Timestamp: mustParseTime(t, "2025-05-04T11:33:51.554181+02:00")},
		{Value: 5, Timestamp: mustParseTime(t, "2025-05-04T11:33:52.554181+02:00")},
		{Value: 8, Timestamp: mustParseTime(t, "2025-05-04T11:33:53.554181+02:00")},
	}, apiResponse.Data)
}

func mustParseTime(t *testing.T, rawTime string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339, rawTime)
	assert.NoError(t, err)

	return parsed
}
