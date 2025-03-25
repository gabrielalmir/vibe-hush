package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheServer(t *testing.T) {
	config := ServerConfig{
		Capacity:          100,
		DefaultExpiration: time.Second * 10,
		AuthToken:         "test-token",
	}

	server := NewCacheServer(config)

	t.Run("Set and Get Item", func(t *testing.T) {
		// Test data
		key := "test-key"
		value := map[string]interface{}{
			"name": "test",
			"age":  float64(30), // JSON numbers are always float64
		}
		item := CacheItem{Value: value}

		// Set item
		body, _ := json.Marshal(item)
		req := httptest.NewRequest("PUT", "/cache/"+key, bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Get item
		req = httptest.NewRequest("GET", "/cache/"+key, nil)
		req.Header.Set("Authorization", "Bearer test-token")

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, value, response["value"])
	})

	t.Run("Get Non-existent Item", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/cache/non-existent", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		w := httptest.NewRecorder()
		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Delete Item", func(t *testing.T) {
		// Set item first
		key := "delete-test"
		item := CacheItem{Value: "test"}

		body, _ := json.Marshal(item)
		req := httptest.NewRequest("PUT", "/cache/"+key, bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Delete item
		req = httptest.NewRequest("DELETE", "/cache/"+key, nil)
		req.Header.Set("Authorization", "Bearer test-token")

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify item is deleted
		req = httptest.NewRequest("GET", "/cache/"+key, nil)
		req.Header.Set("Authorization", "Bearer test-token")

		w = httptest.NewRecorder()
		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Authentication Required", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/cache/test", nil)
		w := httptest.NewRecorder()
		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/cache/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")

		w := httptest.NewRecorder()
		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
