package service

import (
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestGenerateSessionHash_UsesImplicitStickyKeyWhenNoExplicitSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest("POST", "/v1/chat/completions", nil)
	c.Set("api_key", &APIKey{ID: 2})

	svc := &OpenAIGatewayService{}
	body := []byte(`{"model":"gpt-5.5","messages":[{"role":"user","content":"hello"}]}`)

	hash1 := svc.GenerateSessionHash(c, body)
	require.NotEmpty(t, hash1)

	body2 := []byte(`{"model":"gpt-5.5","messages":[{"role":"user","content":"completely different"}]}`)
	hash2 := svc.GenerateSessionHash(c, body2)
	require.Equal(t, hash1, hash2, "implicit sticky should ignore prompt drift for same api key + model")
}

func TestGenerateSessionHash_ExplicitSessionOverridesImplicitSticky(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest("POST", "/v1/chat/completions", nil)
	c.Set("api_key", &APIKey{ID: 2})
	c.Request.Header.Set("conversation_id", "chat-001")

	svc := &OpenAIGatewayService{}
	body := []byte(`{"model":"gpt-5.5","messages":[{"role":"user","content":"hello"}]}`)

	hash := svc.GenerateSessionHash(c, body)
	require.NotEmpty(t, hash)
	implicitHash := DeriveSessionHashFromSeed(buildOpenAIImplicitStickyKey(2, "gpt-5.5"))
	require.NotEqual(t, implicitHash, hash, "explicit session should take priority over implicit sticky")
}

func TestGenerateSessionHash_ImplicitStickyDisabledFallsBackToContentSeed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest("POST", "/v1/chat/completions", nil)
	c.Set("api_key", &APIKey{ID: 2})

	svc := &OpenAIGatewayService{cfg: testConfigWithImplicitSticky(false, 1800)}
	body1 := []byte(`{"model":"gpt-5.5","messages":[{"role":"user","content":"hello"}]}`)
	body2 := []byte(`{"model":"gpt-5.5","messages":[{"role":"user","content":"different"}]}`)

	hash1 := svc.GenerateSessionHash(c, body1)
	hash2 := svc.GenerateSessionHash(c, body2)
	require.NotEmpty(t, hash1)
	require.NotEmpty(t, hash2)
	require.NotEqual(t, hash1, hash2, "when implicit sticky is disabled, content seed should vary with request body")
}

func testConfigWithImplicitSticky(enabled bool, ttl int) *config.Config {
	cfg := &config.Config{}
	cfg.Gateway.OpenAIWS.ImplicitStickyEnabled = enabled
	cfg.Gateway.OpenAIWS.ImplicitStickyTTLSeconds = ttl
	return cfg
}
