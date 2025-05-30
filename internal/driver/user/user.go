package user_drver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/FlyKarlik/effectiveMobile/internal/domain"
	"github.com/FlyKarlik/effectiveMobile/pkg/logger"
)

type IUserDriver interface {
	GetUserAge(ctx context.Context, name string) (int, error)
	GetUserNationality(ctx context.Context, name string) (string, error)
	GetUserSex(ctx context.Context, name string) (domain.SexEnum, error)
}

type userDriver struct {
	logger logger.Logger
}

func New(logger logger.Logger) *userDriver {
	return &userDriver{
		logger: logger,
	}
}

func (u *userDriver) GetUserAge(ctx context.Context, name string) (int, error) {
	const method = "GetUserAge"
	const layer = "driver"
	u.logger.Debug(layer, method, "started", "name", name)

	var userAge struct {
		Age int `json:"age"`
	}

	URL := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	u.logger.Debug(layer, method, "making request", "url", URL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		u.logger.Error(layer, method, "request failed", err, "url", URL)
		return 0, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		u.logger.Error(layer, method, "failed to client do request", err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		u.logger.Error(layer, method, "bad response status", err, "status", resp.StatusCode)
		return 0, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		u.logger.Error(layer, method, "failed to read response body", err)
		return 0, err
	}

	if err := json.Unmarshal(body, &userAge); err != nil {
		u.logger.Error(layer, method, "failed to unmarshal response", err, "body", string(body))
		return 0, err
	}

	u.logger.Debug(layer, method, "successfully completed", "name", name, "age", userAge.Age)
	return userAge.Age, nil
}

func (u *userDriver) GetUserNationality(ctx context.Context, name string) (string, error) {
	const method = "GetUserNationality"
	const layer = "driver"
	u.logger.Debug(layer, method, "started", "name", name)

	var userNationality struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}

	URL := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	u.logger.Debug(layer, method, "making request", "url", URL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		u.logger.Error(layer, method, "request failed", err, "url", URL)
		return "", err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		u.logger.Error(layer, method, "failed to client do request", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		u.logger.Error(layer, method, "bad response status", err, "status", resp.StatusCode)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		u.logger.Error(layer, method, "failed to read response body", err)
		return "", err
	}

	if err := json.Unmarshal(body, &userNationality); err != nil {
		u.logger.Error(layer, method, "failed to unmarshal response", err, "body", string(body))
		return "", err
	}

	if len(userNationality.Country) == 0 {
		u.logger.Warn(layer, method, "no country found", nil, "name", name)
		return "", nil
	}

	nationality := userNationality.Country[0].CountryID
	u.logger.Debug(layer, method, "successfully completed", "name", name, "nationality", nationality)
	return nationality, nil
}

func (u *userDriver) GetUserSex(ctx context.Context, name string) (domain.SexEnum, error) {
	const method = "GetUserSex"
	const layer = "driver"
	u.logger.Debug(layer, method, "started", "name", name)

	var userGender struct {
		Gender string `json:"gender"`
	}

	URL := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	u.logger.Debug(layer, method, "making request", "url", URL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		u.logger.Error("driver", method, "request failed", err, "url", URL)
		return "", err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		u.logger.Error(layer, method, "failed to client do request", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		u.logger.Error(layer, method, "bad response status", err, "status", resp.StatusCode)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		u.logger.Error(layer, method, "failed to read response body", err)
		return "", err
	}

	if err := json.Unmarshal(body, &userGender); err != nil {
		u.logger.Error(layer, method, "failed to unmarshal response", err, "body", string(body))
		return "", err
	}

	sex := getSexEnum(userGender.Gender)
	u.logger.Debug(layer, method, "successfully completed", "name", name, "sex", sex)
	return sex, nil
}

func getSexEnum(s string) domain.SexEnum {
	switch s {
	case "male":
		return domain.MaleSexEnum
	case "female":
		return domain.FemaleSexEnum
	}
	return ""
}
