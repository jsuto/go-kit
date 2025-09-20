package sessionutils

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "strconv"
    "time"

    "github.com/gofiber/fiber/v2"
)

// SessionMiddlewareConfig defines the config for the session middleware
type SessionMiddlewareConfig struct {
    Store                Store
    CookieName           string
    Secure               bool
    SessionDuration      time.Duration
    RegenerateAfter      time.Duration
}

// NewSessionMiddleware returns a Fiber middleware that handles
// session ID creation, TTL refreshing, and optional session ID rotation
func NewSessionMiddleware(config SessionMiddlewareConfig) fiber.Handler {
    return func(c *fiber.Ctx) error {
        ctx := context.Background()

        sessionID, isNew, err := GetOrCreateSessionID(c, config.CookieName, config.Secure, config.SessionDuration)
        if err != nil {
            return fiber.ErrInternalServerError
        }

        if isNew {
            // New session, set created_at
            err := config.Store.HSet(ctx, sessionID, map[string]string{
                "created_at": strconv.FormatInt(time.Now().Unix(), 10),
            })
            if err != nil {
                return fiber.ErrInternalServerError
            }
        } else {
            // Existing session, refresh TTL
            if err := config.Store.Expire(ctx, sessionID, config.SessionDuration); err != nil {
                return fiber.ErrInternalServerError
            }

            // Optionally rotate session ID
            createdAtStr, err := config.Store.HGet(ctx, sessionID, "created_at")
            if err == nil {
                createdAtUnix, err := strconv.ParseInt(createdAtStr, 10, 64)
                if err == nil {
                    createdAt := time.Unix(createdAtUnix, 0)
                    if time.Since(createdAt) > config.RegenerateAfter {
                        newSessionID, err := rotateSessionID(ctx, c, config, sessionID)
                        if err != nil {
                            return fiber.ErrInternalServerError
                        }
                        sessionID = newSessionID
                    }
                }
            }
        }

        // Store session ID in Fiber locals
        c.Locals("session_id", sessionID)

        return c.Next()
    }
}

// GetOrCreateSessionID checks if a session ID cookie exists, otherwise creates one
func GetOrCreateSessionID(c *fiber.Ctx, cookieName string, secure bool, sessionDuration time.Duration) (sessionID string, isNew bool, err error) {
    cookie := c.Cookies(cookieName)
    if cookie != "" {
        return cookie, false, nil
    }

    // Generate new session ID
    sessionID, err = generateSessionID()
    if err != nil {
        return "", false, err
    }

    // Set cookie
    c.Cookie(&fiber.Cookie{
        Name:     cookieName,
        Value:    sessionID,
        HTTPOnly: true,
        Secure:   secure,
        SameSite: "Lax",
        Path:     "/",
        Expires:  time.Now().Add(sessionDuration),
    })

    return sessionID, true, nil
}

// generateSessionID creates a new random session ID
func generateSessionID() (string, error) {
    bytes := make([]byte, 32)
    _, err := rand.Read(bytes)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

// rotateSessionID generates a new session ID, copies data from old session, and deletes old session
func rotateSessionID(ctx context.Context, c *fiber.Ctx, config SessionMiddlewareConfig, oldSessionID string) (string, error) {
    // Dump old session data
    oldData, err := config.Store.HGetAll(ctx, oldSessionID)
    if err != nil {
        return "", err
    }

    // Generate new session ID
    newSessionID, err := generateSessionID()
    if err != nil {
        return "", err
    }

    // Copy old data to new session
    oldData["created_at"] = strconv.FormatInt(time.Now().Unix(), 10)
    if err := config.Store.HSet(ctx, newSessionID, oldData); err != nil {
        return "", err
    }

    // Set TTL for new session
    if err := config.Store.Expire(ctx, newSessionID, config.SessionDuration); err != nil {
        return "", err
    }

    // Delete old session
    _ = config.Store.Clear(ctx, oldSessionID)

    // Set new cookie
    c.Cookie(&fiber.Cookie{
        Name:     config.CookieName,
        Value:    newSessionID,
        HTTPOnly: true,
        Secure:   config.Secure,
        SameSite: "Lax",
        Path:     "/",
        Expires:  time.Now().Add(config.SessionDuration),
    })

    return newSessionID, nil
}

// GetSessionID safely extracts session ID from Fiber Locals
func GetSessionID(c *fiber.Ctx) (string, error) {
    val := c.Locals("session_id")
    if val == nil {
        return "", fiber.ErrUnauthorized
    }
    sessionID, ok := val.(string)
    if !ok || sessionID == "" {
        return "", fiber.ErrUnauthorized
    }
    return sessionID, nil
}

// MustGetSessionID returns the session ID or panics if not found
func MustGetSessionID(c *fiber.Ctx) string {
    sessionID, err := GetSessionID(c)
    if err != nil {
        panic("missing session ID")
    }
    return sessionID
}
