package ai

import (
    "context"
    "time"
    "github.com/go-redis/redis/v8"
)

type ModelVersion struct {
    Version     string    `json:"version"`
    DeployedAt  time.Time `json:"deployed_at"`
    Performance struct {
        Accuracy    float64 `json:"accuracy"`
        F1Score     float64 `json:"f1_score"`
        Precision   float64 `json:"precision"`
        Recall      float64 `json:"recall"`
    } `json:"performance"`
    Config map[string]interface{} `json:"config"`
}

type ModelVersioning struct {
    redis *redis.Client
}

func NewModelVersioning(redis *redis.Client) *ModelVersioning {
    return &ModelVersioning{redis: redis}
}

func (mv *ModelVersioning) RegisterVersion(ctx context.Context, modelName string, version ModelVersion) error {
    key := mv.getVersionKey(modelName, version.Version)
    return mv.redis.Set(ctx, key, version, 0).Err()
}

func (mv *ModelVersioning) GetCurrentVersion(ctx context.Context, modelName string) (*ModelVersion, error) {
    key := mv.getCurrentVersionKey(modelName)
    data, err := mv.redis.Get(ctx, key).Result()
    if err != nil {
        return nil, err
    }

    var version ModelVersion
    if err := json.Unmarshal([]byte(data), &version); err != nil {
        return nil, err
    }

    return &version, nil
}

func (mv *ModelVersioning) SetCurrentVersion(ctx context.Context, modelName, version string) error {
    key := mv.getCurrentVersionKey(modelName)
    return mv.redis.Set(ctx, key, version, 0).Err()
}

func (mv *ModelVersioning) getVersionKey(modelName, version string) string {
    return fmt.Sprintf("model:%s:version:%s", modelName, version)
}

func (mv *ModelVersioning) getCurrentVersionKey(modelName string) string {
    return fmt.Sprintf("model:%s:current_version", modelName)
}