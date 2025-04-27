package session

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// Store defines a generic session store interface
type Store interface {
    Save(ctx context.Context, sessionID, key string, value any) error
    Load(ctx context.Context, sessionID, key string) ([]byte, error)
    LoadJSON(ctx context.Context, sessionID, key string, dest any) error
    Delete(ctx context.Context, sessionID, key string) error
    Clear(ctx context.Context, sessionID string) error

    HSet(ctx context.Context, sessionID string, values map[string]string) error
    HGetAll(ctx context.Context, sessionID string) (map[string]string, error)
    HGet(ctx context.Context, sessionID, key string) (string, error)
    Expire(ctx context.Context, sessionID string, expiration time.Duration) error
}

// SessionManager handles Redis-backed sessions
type SessionManager struct {
    RedisClient *redis.Client
}

// NewSessionManager creates a new Redis-backed SessionManager
func NewSessionManager(redisClient *redis.Client) *SessionManager {
    return &SessionManager{
        RedisClient: redisClient,
    }
}

// Save saves a Go value into the session
func (sm *SessionManager) Save(ctx context.Context, sessionID, key string, value any) error {
    fullKey := "session:" + sessionID

    jsonValue, err := json.Marshal(value)
    if err != nil {
        return fmt.Errorf("failed to marshal session value: %w", err)
    }

    if err := sm.RedisClient.HSet(ctx, fullKey, key, jsonValue).Err(); err != nil {
        return fmt.Errorf("failed to save session data: %w", err)
    }

    return nil
}

// Load loads a raw value (as []byte) from the session
func (sm *SessionManager) Load(ctx context.Context, sessionID, key string) ([]byte, error) {
    fullKey := "session:" + sessionID

    data, err := sm.RedisClient.HGet(ctx, fullKey, key).Result()
    if err == redis.Nil {
        return nil, fmt.Errorf("session key %q not found", key)
    }
    if err != nil {
        return nil, fmt.Errorf("failed to load session data: %w", err)
    }

    return []byte(data), nil
}

// LoadJSON unmarshals a Go value from the session
func (sm *SessionManager) LoadJSON(ctx context.Context, sessionID, key string, dest any) error {
    raw, err := sm.Load(ctx, sessionID, key)
    if err != nil {
        return err
    }

    if err := json.Unmarshal(raw, dest); err != nil {
        return fmt.Errorf("failed to unmarshal session value: %w", err)
    }

    return nil
}

// Delete deletes a key from the session
func (sm *SessionManager) Delete(ctx context.Context, sessionID, key string) error {
    fullKey := "session:" + sessionID

    if err := sm.RedisClient.HDel(ctx, fullKey, key).Err(); err != nil {
        return fmt.Errorf("failed to delete session key %q: %w", key, err)
    }

    return nil
}

// Clear deletes the entire session
func (sm *SessionManager) Clear(ctx context.Context, sessionID string) error {
    fullKey := "session:" + sessionID

    if err := sm.RedisClient.Del(ctx, fullKey).Err(); err != nil {
        return fmt.Errorf("failed to clear session: %w", err)
    }

    return nil
}

// HSet sets multiple fields in the session
func (sm *SessionManager) HSet(ctx context.Context, sessionID string, values map[string]string) error {
    fullKey := "session:" + sessionID
    return sm.RedisClient.HSet(ctx, fullKey, values).Err()
}

// HGetAll gets all fields from the session
func (sm *SessionManager) HGetAll(ctx context.Context, sessionID string) (map[string]string, error) {
    fullKey := "session:" + sessionID
    return sm.RedisClient.HGetAll(ctx, fullKey).Result()
}

// HGet gets a single field from the session
func (sm *SessionManager) HGet(ctx context.Context, sessionID, key string) (string, error) {
    fullKey := "session:" + sessionID
    return sm.RedisClient.HGet(ctx, fullKey, key).Result()
}

// Expire sets an expiration time for the session
func (sm *SessionManager) Expire(ctx context.Context, sessionID string, expiration time.Duration) error {
    fullKey := "session:" + sessionID
    return sm.RedisClient.Expire(ctx, fullKey, expiration).Err()
}
