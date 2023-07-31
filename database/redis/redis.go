package redis

import (
	"context"
	"encoding/json"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/types"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func GetClient(Addr string, Password string, DB int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password, // no password set
		DB:       DB,       // use default DB // TODO
	})
}

func SaveToRedis(client *redis.Client, key string, track types.TrackInfo) error {
	ctx := context.Background()

	jsonData, err := json.Marshal(track)
	if err != nil {
		return err
	}

	err = client.Set(ctx, key, jsonData, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetFromRedis(client *redis.Client, key int) (types.TrackInfo, error) {
	ctx := context.Background()

	data, err := client.Get(ctx, strconv.Itoa(key)).Result()
	if err != nil {
		return types.TrackInfo{}, err
	}

	var lastTrackInfo types.TrackInfo
	err = json.Unmarshal([]byte(data), &lastTrackInfo)
	if err != nil {
		return types.TrackInfo{}, err
	}

	return lastTrackInfo, nil
}
