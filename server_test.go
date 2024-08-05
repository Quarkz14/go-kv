package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestGetAllKeys(t *testing.T) {
	t.Run("return all the keys stored", func(t *testing.T) {
		t.Parallel()
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		s := &Server{NewKeyValueStore()}
		s.store.PutOrCreateValue("key", "value")
		s.store.PutOrCreateValue("blah", "value")
		s.GetKeys(response, request)

		var jsonResponse map[string][]string

		json.NewDecoder(response.Body).Decode(&jsonResponse)

		got := jsonResponse["Keys"]
		want := []string{"key", "blah"}

		if ok := slices.Equal(got, want); ok != true {
			t.Errorf("Got: %s\n wanted: %s\n", got, want)
		}

		if response.Code != http.StatusOK {
			t.Errorf("Got: %d\n want: %d\n", response.Code, http.StatusOK)
		}
	})

	t.Run("return 404 if no keys are present", func(t *testing.T) {
		t.Parallel()
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		s := &Server{NewKeyValueStore()}
		s.GetKeys(response, request)

		if response.Code != http.StatusNotFound {
			t.Errorf("Got: %d\n want: %d\n", response.Code, http.StatusNotFound)
		}
	})
}

func TestPutOrUpdateValue(t *testing.T) {
	t.Run("update an existing key with a new value", func(t *testing.T) {
		t.Parallel()
		jsonBody := []byte(`{"value": "superdupersecret"}`)
		bodyReader := bytes.NewReader(jsonBody)

		request, _ := http.NewRequest(http.MethodPut, "/", bodyReader)
		response := httptest.NewRecorder()
		request.SetPathValue("key", "mysecret")

		s := &Server{NewKeyValueStore()}
		s.store.PutOrCreateValue("mysecret", "blah")
		s.PutOrCreateValue(response, request)

		if response.Code != http.StatusOK {
			t.Errorf("Got: %d \n want: %d\n", response.Code, http.StatusOK)
		}
	})
	t.Run("create a key with a value if key doesn't exist", func(t *testing.T) {
		t.Parallel()
		jsonBody := []byte(`{"value": "superdupersecret"}`)
		bodyReader := bytes.NewReader(jsonBody)

		request, _ := http.NewRequest(http.MethodPut, "/", bodyReader)
		response := httptest.NewRecorder()
		request.SetPathValue("key", "newsecret")
		s := &Server{NewKeyValueStore()}
		s.PutOrCreateValue(response, request)

		if response.Code != http.StatusOK {
			t.Errorf("Got: %d \n want: %d\n", response.Code, http.StatusOK)
		}
	})
	t.Run("put with wrong json body should return bad request", func(t *testing.T) {
		t.Parallel()
		jsonBody := []byte(`{"notvalue": "superdupersecret"}`)
		bodyReader := bytes.NewReader(jsonBody)

		request, _ := http.NewRequest(http.MethodPut, "/", bodyReader)
		response := httptest.NewRecorder()
		request.SetPathValue("key", "mysecret")
		s := &Server{NewKeyValueStore()}
		s.PutOrCreateValue(response, request)

		if response.Code != http.StatusBadRequest {
			t.Errorf("Got: %d \n want: %d\n", response.Code, http.StatusBadRequest)
		}
	})
}

func TestGetValue(t *testing.T) {
	t.Run("get a value when given a key", func(t *testing.T) {
		t.Parallel()
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		request.SetPathValue("key", "mysecret")

		s := &Server{NewKeyValueStore()}
		s.store.PutOrCreateValue("mysecret", "blah")

		s.GetValue(response, request)

		if response.Code != http.StatusOK {
			t.Errorf("Got: %d \n want: %d\n", response.Code, http.StatusOK)
		}

		var kv Get
		json.NewDecoder(response.Body).Decode(&kv)
		got := kv
		want := Get{
			Key:   "mysecret",
			Value: "blah",
		}
		if got != want {
			t.Errorf("Got: %+v\n want: %+v\n", got, want)
		}
	})

	t.Run("get a value when a key doesn't exist", func(t *testing.T) {
		t.Parallel()
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		request.SetPathValue("key", "mysecret")

		s := &Server{NewKeyValueStore()}

		s.GetValue(response, request)

		if response.Code != http.StatusNotFound {
			t.Errorf("Got: %d \n want: %d\n", response.Code, http.StatusNotFound)
		}
	})
}

func TestDeleteValue(t *testing.T) {
	t.Run("delete a value when given a key", func(t *testing.T) {
		t.Parallel()
		request, _ := http.NewRequest(http.MethodDelete, "/", nil)
		response := httptest.NewRecorder()

		request.SetPathValue("key", "mysecret")

		s := &Server{NewKeyValueStore()}
		s.store.PutOrCreateValue("mysecret", "blah")

		s.DeleteKeyValue(response, request)

		if response.Code != http.StatusOK {
			t.Errorf("Got: %d \n want: %d\n", response.Code, http.StatusOK)
		}
	})

	t.Run("deleting a value when a key doesn't exist should return 404", func(t *testing.T) {
		t.Parallel()
		request, _ := http.NewRequest(http.MethodDelete, "/", nil)
		response := httptest.NewRecorder()

		request.SetPathValue("key", "mysecret")

		s := &Server{NewKeyValueStore()}

		s.DeleteKeyValue(response, request)

		if response.Code != http.StatusNotFound {
			t.Errorf("Got: %d \n want: %d\n", response.Code, http.StatusNotFound)
		}
	})
}
