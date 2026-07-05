package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BenDowswell/go-kv-store/kvstore"
)

func TestGetHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	getHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	if rr.Body.String() != "Hello World!" {
		t.Errorf("expected Hello World!, got %q", rr.Body.String())
	}
}

func TestGetMissingKeyReturns404(t *testing.T) {
	store := kvstore.NewKVStore()

	req := httptest.NewRequest(http.MethodGet, "/kv/name", nil)
	req.SetPathValue("key", "name")

	rr := httptest.NewRecorder()

	getKey(rr, req, store)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

func TestGetExistingKeyReturnsValue(t *testing.T) {
	store := kvstore.NewKVStore()
	store.Set("name", "Ben")

	req := httptest.NewRequest(http.MethodGet, "/kv/name", nil)
	req.SetPathValue("key", "name")

	rr := httptest.NewRecorder()

	getKey(rr, req, store)

	// expect 200
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	// expect body contains Ben
	if rr.Body.String() != "Ben\n" {
		t.Errorf("expected body %q, got %q", "Ben\n", rr.Body.String())
	}

}
func TestPutStoresNewValue(t *testing.T) {
	store := kvstore.NewKVStore()

	body := strings.NewReader(`{"value":"Ben"}`)
	req := httptest.NewRequest(http.MethodPut, "/kv/name", body)
	req.SetPathValue("key", "name")
	rr := httptest.NewRecorder()

	updateKey(rr, req, store)

	value, ok := store.Get("name")
	if !ok || value != "Ben" {
		t.Errorf("expected name to be Ben and ok true, got value=%q ok=%v", value, ok)
	}

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rr.Code)
	}
}
func TestPutOverwritesExistingValue(t *testing.T) {
	store := kvstore.NewKVStore()
	store.Set("name", "Ben")
	body := strings.NewReader(`{"value":"Tony"}`)
	req := httptest.NewRequest(http.MethodPut, "/kv/name", body)
	req.SetPathValue("key", "name")
	rr := httptest.NewRecorder()

	updateKey(rr, req, store)

	value, ok := store.Get("name")
	if !ok || value != "Tony" {
		t.Errorf("expected name to be Tony and ok true, got value=%q ok=%v", value, ok)
	}

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rr.Code)
	}

}

func TestDeleteRemovesExistingKey(t *testing.T) {
	store := kvstore.NewKVStore()
	store.Set("name", "Ben")

	req := httptest.NewRequest(http.MethodDelete, "/kv/name", nil)
	req.SetPathValue("key", "name")
	rr := httptest.NewRecorder()
	deleteKey(rr, req, store)

	// expect 204
	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rr.Code)
	}

	_, ok := store.Get("name")
	if ok {
		t.Errorf("expected key to be deleted")
	}

}

func TestDeleteMissingKeyReturnsSuccess(t *testing.T) {
	store := kvstore.NewKVStore()
	req := httptest.NewRequest(http.MethodDelete, "/kv/name", nil)
	req.SetPathValue("key", "name")
	rr := httptest.NewRecorder()

	deleteKey(rr, req, store)

	// expect 204
	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rr.Code)
	}

}

func TestPutInvalidJSONReturns400(t *testing.T) {
	store := kvstore.NewKVStore()
	body := strings.NewReader(`{"value":`)
	req := httptest.NewRequest(http.MethodPut, "/kv/name", body)
	req.SetPathValue("key", "name")
	rr := httptest.NewRecorder()

	updateKey(rr, req, store)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestHealthIntegration(t *testing.T) {
	store := kvstore.NewKVStore()

	handler := Routes(store)
	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/health")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if string(body) != "Hello World!" {
		t.Errorf("expected Hello World!, got %q", string(body))
	}
}

func TestEndToEndIntegration(t *testing.T) {
	store := kvstore.NewKVStore()

	handler := Routes(store)
	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	body := strings.NewReader(`{"value":"Ben"}`)
	req, err := http.NewRequest(http.MethodPut, testServer.URL+"/kv/name", body)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Put a value in
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("PUT request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}

	// Get value out
	getName, err := http.Get(testServer.URL + "/kv/name")
	if err != nil {
		t.Fatalf("failed to get key: %v", err)
	}
	defer getName.Body.Close()
	bodyBytes, err := io.ReadAll(getName.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	if getName.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", getName.StatusCode)
	}

	if string(bodyBytes) != "Ben\n" {
		t.Errorf("expected Ben, got %q", string(bodyBytes))
	}

	// Delete value
	deleteReq, err := http.NewRequest(http.MethodDelete, testServer.URL+"/kv/name", nil)
	if err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}
	delResp, err := client.Do(deleteReq)
	if err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}
	defer delResp.Body.Close()

	if delResp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", delResp.StatusCode)
	}

	// GET but expect 404
	getName404, err := http.Get(testServer.URL + "/kv/name")
	if err != nil {
		t.Fatalf("failed to get key: %v", err)
	}
	defer getName404.Body.Close()

	if getName404.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", getName404.StatusCode)
	}
}
