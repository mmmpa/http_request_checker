package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSlackChallenge(t *testing.T) {
	mux := build()

	challengeBody := `{ "type": "url_verification", "challenge": "challenge" }`
	notChallengeBody := `{ "type": "other", "challenge": "challenge" }`
	expectedRequest := func() *http.Request { return httptest.NewRequest(http.MethodPost, "/", strings.NewReader(challengeBody)) }
	notExpectedRequest := func() *http.Request { return httptest.NewRequest(http.MethodPost, "/", strings.NewReader(notChallengeBody)) }

	isChallengeResponse := func(t *testing.T, w *httptest.ResponseRecorder) {
		assert.Equal(t, w.Code, 200)
		assert.Equal(t, w.Header().Get("content-type"), "text/plain")
		assert.Equal(t, w.Body.String(), "challenge")
	}

	isNotChallengeResponse := func(t *testing.T, w *httptest.ResponseRecorder) {
		assert.Equal(t, w.Code, 200)
		assert.Equal(t, w.Header().Get("content-type"), "")
		assert.NotEqual(t, w.Body.String(), "challenge")
		assert.Contains(t, w.Body.String(), "challenge")
	}

	t.Run("enabled challenge", func(t *testing.T) {
		disabledSlackChallenge = false

		t.Run("with valid request", func(t *testing.T) {
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, expectedRequest())

			isChallengeResponse(t, w)
		})

		t.Run("with invalid request", func(t *testing.T) {
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, notExpectedRequest())

			isNotChallengeResponse(t, w)
		})
	})

	t.Run("disabled challenge", func(t *testing.T) {
		disabledSlackChallenge = true

		t.Run("with valid request", func(t *testing.T) {
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, expectedRequest())

			isNotChallengeResponse(t, w)
		})

		t.Run("with invalid request", func(t *testing.T) {
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, notExpectedRequest())

			isNotChallengeResponse(t, w)
		})
	})
}
