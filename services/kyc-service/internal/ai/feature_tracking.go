package ai

import (
    "context"
    "time"
    "github.com/go-redis/redis/v8"
)

type FeatureImportance struct {
    ModelName    string                 `json:"model_name"`
    Version      string                 `json:"version"`
    Timestamp    time.Time             `json:"timestamp"`
    Features     map[string]float64    `json:"features"`
    GlobalShap   map[string]float64    `json:"global_shap"`
    LocalShap    map[string]float64    `json:"local_shap"`
}

type FeatureTracker struct {
    redis *redis.Client
}

func NewFeatureTracker(redis *redis.Client) *FeatureTracker {
    return &FeatureTracker{redis: redis}
}

func (ft *FeatureTracker) TrackFeatureImportance(ctx context.Context, importance FeatureImportance) error {
    key := ft.getFeatureKey(importance.ModelName, importance.Version)
    return ft.redis.Set(ctx, key, importance, 24*time.Hour).Err()
}

func (ft *FeatureTracker) GetFeatureImportance(ctx context.Context, modelName, version string) (*FeatureImportance, error) {
    key := ft.getFeatureKey(modelName, version)
    data, err := ft.redis.Get(ctx, key).Result()
    if err != nil {
        return nil, err
    }

    var importance FeatureImportance
    if err := json.Unmarshal([]byte(data), &importance); err != nil {
        return nil, err
    }

    return &importance, nil
}

func (ft *FeatureTracker) getFeatureKey(modelName, version string) string {
    return fmt.Sprintf("feature_importance:%s:%s", modelName, version)
}