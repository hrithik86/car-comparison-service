package vehicle

import (
	"car-comparison-service/caching"
	"car-comparison-service/config"
	"car-comparison-service/db/model"
	"context"
	"encoding/json"
	"fmt"
)

const suggestedKeySuffix = "{suggestions_set}"

type IVehicleSuggestions interface {
	CreateKey(ctx context.Context, vehicleId string) string
	SetVehicleSuggestionsDetails(ctx context.Context, key string, res []model.VehicleSuggestionResult) error
	GetVehicleSuggestionsDetails(ctx context.Context, key string) ([]model.VehicleSuggestionResult, error)
}

func NewVehicleSuggestionsManager() IVehicleSuggestions {
	return &Suggestions{}
}

type Suggestions struct {
}

func (v *Suggestions) CreateKey(ctx context.Context, vehicleId string) string {
	return fmt.Sprintf("%s_%s", suggestedKeySuffix, vehicleId)
}

func (v *Suggestions) SetVehicleSuggestionsDetails(ctx context.Context, key string, res []model.VehicleSuggestionResult) error {
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	err = caching.GetRedisClient().SetWithExpiry(ctx, key, string(data), config.RedisConf().RedisKeyExpiryTimeout)
	return err
}

func (v *Suggestions) GetVehicleSuggestionsDetails(ctx context.Context, key string) ([]model.VehicleSuggestionResult, error) {
	res, err := caching.GetRedisClient().Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var result []model.VehicleSuggestionResult
	err = json.Unmarshal([]byte(res), &result)
	return result, err
}
