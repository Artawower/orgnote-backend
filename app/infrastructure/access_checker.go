package infrastructure

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AccessChecker struct {
	checkURL   *string
	httpClient http.Client
}

type AccessInfo struct {
	Email      string `json:"email"`
	IsActive   bool   `json:"isActive"`
	SpaceLimit int    `json:"spaceLimit"`
}

type AccessCheckRequestError struct {
	code int
}

func (e *AccessCheckRequestError) Error() string {
	return "access check request error, got status code " + string(rune(e.code))
}

func (a *AccessChecker) checkAvailability(info AccessInfo, occupiedSpace int) error {
	if !info.IsActive {
		return fmt.Errorf("user %s is not active", info.Email)
	}
	if (info.SpaceLimit - occupiedSpace) < 0 {
		return fmt.Errorf("user %s has no space left", info.Email)
	}

	return nil
}

func (a *AccessChecker) getRemoteInfo(userEmail string) (*AccessInfo, error) {
	req, err := http.NewRequest(http.MethodGet, *a.checkURL, nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("email", userEmail)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Content-Type", "Application/json")

	res, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, &AccessCheckRequestError{res.StatusCode}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data := &AccessInfo{}
	err = json.Unmarshal(body, &data)

	return data, err

}

func (a *AccessChecker) Check(userEmail string, occupiedSpace int, errCh chan<- error) {
	// TODO: master cache here
	if a.checkURL == nil {
		errCh <- nil
		return
	}

	accessInfo, err := a.getRemoteInfo(userEmail)

	if err != nil {
		errCh <- err
		return
	}

	errCh <- a.checkAvailability(*accessInfo, occupiedSpace)
}

func NewAccessChecker(httpClient http.Client, checkURL *string) *AccessChecker {
	return &AccessChecker{
		checkURL:   checkURL,
		httpClient: httpClient,
	}
}
