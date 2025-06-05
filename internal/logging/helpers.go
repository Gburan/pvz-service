package logging

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type logCtx struct {
	RequestID       uuid.UUID
	Status          int
	RequestStart    string
	RequestDuration string
	Method          string
	Path            string

	PVZID       uuid.UUID
	ProductType string
	City        string
	Role        string
	Email       string
	StartDate   time.Time
	EndDate     time.Time
	Page        int
	Limit       int
}

func WithLogLimit(ctx context.Context, limit int) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Limit = limit
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Limit: limit})
}

func WithLogPage(ctx context.Context, page int) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Page = page
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Page: page})
}

func WithLogEndDate(ctx context.Context, date time.Time) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.EndDate = date
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{EndDate: date})
}

func WithLogStartDate(ctx context.Context, date time.Time) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.StartDate = date
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{StartDate: date})
}

func WithLogEmail(ctx context.Context, email string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Email = email
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Email: email})
}

func WithLogRole(ctx context.Context, role string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Role = role
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Role: role})
}

func WithLogCity(ctx context.Context, city string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.City = city
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{City: city})
}

func WithLogProductType(ctx context.Context, prodType string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.ProductType = prodType
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{ProductType: prodType})
}

func WithLogPVZID(ctx context.Context, pvzID uuid.UUID) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.PVZID = pvzID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{PVZID: pvzID})
}

func WithLogRequestID(ctx context.Context, requestID uuid.UUID) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.RequestID = requestID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{RequestID: requestID})
}

func WithLogRequestPath(ctx context.Context, path string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Path = path
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Path: path})
}

func WithLogRequestMethod(ctx context.Context, method string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Method = method
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Method: method})
}

func WithLogRequestStatus(ctx context.Context, status int) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Status = status
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Status: status})
}

func WithLogRequestDuration(ctx context.Context, duration string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.RequestDuration = duration
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{RequestDuration: duration})
}
