package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	subscription "orgnote/app/infrastructure/generated"
	"orgnote/app/tools"
	"strconv"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/rs/zerolog/log"
	"github.com/thoas/go-funk"
)

type SubscriptionAPI struct {
	checkURL      *string
	checkToken    *string
	client        *subscription.ClientWithResponses
	cache         *cache.Cache[string, SubscriptionInfo]
	cacheLifeTime int
}

type SubscriptionInfo struct {
	Email      string  `json:"email"`
	IsActive   bool    `json:"isActive"`
	SpaceLimit float64 `json:"spaceLimit"`
}

type SubscriptionRequestError struct {
	code int
}

func (e *SubscriptionRequestError) Error() string {
	return "access check request error, got status code " + strconv.FormatInt(int64(e.code), 10)
}

func (a *SubscriptionAPI) checkAvailability(info SubscriptionInfo, usedSpace int64) error {
	if !info.IsActive {
		return fmt.Errorf("user %s is not active", info.Email)
	}

	usedSpaceMb := tools.ConvertBytes2Megabyte(usedSpace)
	spaceLimit := tools.ConvertBytes2Megabyte(int64(info.SpaceLimit))
	if (spaceLimit - usedSpaceMb) < 0 {
		return fmt.Errorf("user %s has no space left, %v/%v are used", info.Email, usedSpaceMb, spaceLimit)
	}

	return nil
}

func (a *SubscriptionAPI) getRemoteInfo(provider string, externalID string) (*SubscriptionInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	res, err := a.client.SubscriptionInfoRetrieve(ctx, provider, externalID, a.addAuthHeader)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, &SubscriptionRequestError{res.StatusCode}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data := &SubscriptionInfo{}
	err = json.Unmarshal(body, &data)

	return data, err

}

func (a *SubscriptionAPI) getInfo(provider string, externalID string) (*SubscriptionInfo, error) {
	key := provider + externalID
	cachedInfo, ok := a.cache.Get(key)

	if ok {
		return &cachedInfo, nil
	}

	accessInfo, err := a.getRemoteInfo(provider, externalID)
	if err != nil {
		return nil, err
	}

	a.cache.Set(externalID, *accessInfo, cache.WithExpiration(time.Duration(a.cacheLifeTime)*time.Minute))

	return accessInfo, err
}

func (a *SubscriptionAPI) Check(provider string, externalID string, usedSpace int64, errCh chan<- error) {
	// TODO: add cache here
	if a.checkURL == nil {
		errCh <- nil
		return
	}

	accessInfo, err := a.getInfo(provider, externalID)

	if err != nil {
		errCh <- err
		return
	}

	errCh <- a.checkAvailability(*accessInfo, usedSpace)
}

func (a *SubscriptionAPI) addAuthHeader(ctx context.Context, req *http.Request) error {
	req.Header.Add("Authorization", "Token "+*a.checkToken)
	return nil
}

var ErrorInvalidToken = fmt.Errorf("invalid activation token")

func (a *SubscriptionAPI) ActivateSubscription(data subscription.SubscriptionActivation) (*subscription.SubscriptionInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// TODO: master map errors
	rspns, err := a.client.SubscriptionActivationCreateWithResponse(ctx, data)

	if err != nil {
		log.Info().Msgf("subscription: activate subscription: subscription request: %v", err)
		return nil, fmt.Errorf("subscription: subscribe: %v", err)
	}

	if rspns.HTTPResponse.StatusCode == http.StatusNotFound {
		return nil, ErrorInvalidToken
	}

	if !funk.Contains([]int{http.StatusOK, http.StatusCreated}, rspns.HTTPResponse.StatusCode) {
		return nil, fmt.Errorf("subscription: subscribe: got status code %v", rspns.HTTPResponse.StatusCode)
	}

	return rspns.JSON200, nil
}

type cacheFactory[K comparable, V any] func(...cache.Option[K, V]) *cache.Cache[K, V]

func NewSubscription(
	httpClient http.Client,
	checkURL *string,
	checkToken *string,
	cacheFactory cacheFactory[string, SubscriptionInfo],
	cacheLifeTime int,
) (*SubscriptionAPI, error) {
	// TODO: master use as dependency
	client, err := subscription.NewClientWithResponses(*checkURL)

	if err != nil {
		return nil, fmt.Errorf("subscription: new subscription: init client: %v", err)
	}

	return &SubscriptionAPI{
		checkURL,
		checkToken,
		client,
		cacheFactory(),
		cacheLifeTime,
	}, nil
}
