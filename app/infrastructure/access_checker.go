package infrastructure

import (
	"encoding/json"
	"fmt"
	"io"
	"moonbrain/app/tools"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

type AccessChecker struct {
	checkURL   *string
	checkToken *string
	httpClient http.Client
}

type AccessInfo struct {
	Email      string  `json:"email"`
	IsActive   bool    `json:"isActive"`
	SpaceLimit float64 `json:"spaceLimit"`
}

type AccessCheckRequestError struct {
	code int
}

func (e *AccessCheckRequestError) Error() string {
	return "access check request error, got status code " + strconv.FormatInt(int64(e.code), 10)
}

func (a *AccessChecker) checkAvailability(info AccessInfo, usedSpace int64) error {
	if !info.IsActive {
		return fmt.Errorf("user %s is not active", info.Email)
	}

	usedSpaceMb := tools.ConvertBytes2Megabyte(usedSpace)
	if (info.SpaceLimit - usedSpaceMb) < 0 {
		return fmt.Errorf("user %s has no space left, %v/%v are used", info.Email, usedSpaceMb, info.SpaceLimit)
	}

	return nil
}

func (a *AccessChecker) getRemoteInfo(userEmail string) (*AccessInfo, error) {
	fullCheckURL := *a.checkURL + "/" + userEmail

	req, err := http.NewRequest(http.MethodGet, fullCheckURL, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "Application/json")

	if a.checkToken != nil {
		req.Header.Add("Authorization", "Token "+*a.checkToken)
	}

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

func (a *AccessChecker) Check(userEmail string, usedSpace int64, errCh chan<- error) {
	// TODO: add cache here
	if a.checkURL == nil {
		log.Info().Msg("access checker: check: no check url provided")
		errCh <- nil
		return
	}

	accessInfo, err := a.getRemoteInfo(userEmail)

	if err != nil {
		errCh <- err
		return
	}

	errCh <- a.checkAvailability(*accessInfo, usedSpace)
}

func NewAccessChecker(httpClient http.Client, checkURL *string, checkToken *string) *AccessChecker {
	return &AccessChecker{
		checkURL,
		checkToken,
		httpClient,
	}
}
