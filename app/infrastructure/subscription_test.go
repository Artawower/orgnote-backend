package infrastructure

import (
	"errors"
	"net/http"
	"net/http/httptest"
	subscription "orgnote/app/infrastructure/generated"
	"testing"

	cache "github.com/Code-Hex/go-generics-cache"
)

func newTestSubscriptionAPI(t *testing.T, handler http.HandlerFunc) *SubscriptionAPI {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	token := "token"
	api, err := NewSubscription(http.Client{}, &server.URL, &token, cache.New[string, SubscriptionInfo], 1)
	if err != nil {
		t.Fatalf("failed to create subscription api: %v", err)
	}

	return api
}

func TestActivateSubscription_ReturnsResponseBody(t *testing.T) {
	api := newTestSubscriptionAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"spaceLimit":1024}`))
	})

	info, err := api.ActivateSubscription(subscription.SubscriptionActivation{Key: "key", ExternalId: "42"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info == nil || info.SpaceLimit == nil || *info.SpaceLimit != 1024 {
		t.Fatalf("expected space limit 1024, got %v", info)
	}
}

func TestActivateSubscription_ReturnsErrorWhenResponseBodyMissing(t *testing.T) {
	api := newTestSubscriptionAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	_, err := api.ActivateSubscription(subscription.SubscriptionActivation{Key: "key", ExternalId: "42"})

	if err == nil {
		t.Fatal("expected error for missing activation response body")
	}
}

func TestActivateSubscription_ReturnsErrorWhenStatusIsCreated(t *testing.T) {
	api := newTestSubscriptionAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	_, err := api.ActivateSubscription(subscription.SubscriptionActivation{Key: "key", ExternalId: "42"})

	if err == nil {
		t.Fatal("expected error for unsupported activation response status")
	}
}

func TestActivateSubscription_ReturnsInvalidTokenError(t *testing.T) {
	api := newTestSubscriptionAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := api.ActivateSubscription(subscription.SubscriptionActivation{Key: "key", ExternalId: "42"})

	if !errors.Is(err, ErrorInvalidToken) {
		t.Fatalf("expected invalid token error, got %v", err)
	}
}
