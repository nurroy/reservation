package redis

import (
	pg "bitbucket.org/pinang-och-go/pkg/v1/postgres"
	"strconv"
)

const (
	APPID_KEY = `APPID`
)

// StoreAppId
func StoreAppId() error {
	// select last application id
	appId, err := pg.SelectPublicPinangLastApplicationID()
	if err != nil {
		return err
	}

	// delete application id from redis
	client.Del(APPID_KEY)

	// store application id to redis
	err = client.Set(APPID_KEY, appId+1, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// Create Application Id
func CreateAppId() (*string, error) {
	// get application id from redis
	appId, err := client.Get(APPID_KEY).Result()
	if err != nil {
		return nil, err
	}

	// delete application id from redis
	client.Del(APPID_KEY)

	appID, err := strconv.Atoi(appId)
	if err != nil {
		return nil, err
	}

	// store application id to redis
	err = client.Set(APPID_KEY, appID+1, 0).Err()
	if err != nil {
		return nil, err
	}

	return &appId, nil
}
