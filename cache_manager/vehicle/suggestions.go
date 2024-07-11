package vehicle

import (
	"car-comparison-service/caching"
	"car-comparison-service/config"
	"car-comparison-service/db/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

const suggestedKeySuffix = "suggestions_set"

type Vehicle struct {
	ctx       context.Context
	vehicleId *uuid.UUID
}

func (v *Vehicle) createKey() string {
	return fmt.Sprintf("%s-%s", v.vehicleId, suggestedKeySuffix)
}

func (v *Vehicle) SetVehicleSuggestionsDetails(res []model.VehicleSuggestionResult) error {
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	err = caching.GetRedisClient().SetWithExpiry(v.ctx, v.createKey(), string(data), config.RedisConf().RedisKeyExpiryTimeout)
	return err
}

func (v *Vehicle) GetVehicleSuggestionsDetails() ([]model.VehicleSuggestionResult, error) {
	res, err := caching.GetRedisClient().Get(v.ctx, v.createKey())
	if err != nil {
		return nil, err
	}
	var result []model.VehicleSuggestionResult
	err = json.Unmarshal([]byte(res), &result)
	return result, err
}

func CreateSuggestionVehicle(ctx context.Context, vehicleId *uuid.UUID) *Vehicle {
	return &Vehicle{
		ctx:       ctx,
		vehicleId: vehicleId,
	}
}
