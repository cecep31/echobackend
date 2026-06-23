package service

import (
	"context"
	"sort"
	"time"

	"echobackend/internal/dto"
	"echobackend/pkg/market"
)

const (
	// corporateActionCacheTTL is how long the calendar result is cached in Redis.
	// Set to 6 hours so that data is reasonably fresh without hammering the API.
	corporateActionCacheTTL = 6 * time.Hour

	// maxCalendarRangeMonths caps the from–to range accepted by the handler.
	maxCalendarRangeMonths = 6

	calendarDateFormat = "2006-01-02"
)

// corporateActionCache is the minimal cache interface required by this service.
type corporateActionCache interface {
	BuildKey(parts ...string) string
	GetJSON(ctx context.Context, key string, dest any) (bool, error)
	SetJSONWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error
}

// CorporateActionService fetches dividend and RUPS events for the user's
// stock holdings and caches the results in Redis.
type CorporateActionService interface {
	GetCalendar(ctx context.Context, userID string, from, to time.Time) (*dto.CorporateActionCalendarResponse, error)
}

type corporateActionService struct {
	idxClient    market.CorporateActionClient // RapidAPI IDX (IDX stocks)
	cache        corporateActionCache
}

// NewCorporateActionService creates the service.
func NewCorporateActionService(
	idxClient market.CorporateActionClient,
	cache corporateActionCache,
) CorporateActionService {
	return &corporateActionService{
		idxClient:   idxClient,
		cache:       cache,
	}
}

// GetCalendar returns corporate actions (dividend + RUPS) within the [from, to] date range.
//
// Results are cached globally for corporateActionCacheTTL. When the cache
// is cold, external API calls are made to IDX. Individual API errors are swallowed
// (fail-open) so that a partial result is always returned.
func (s *corporateActionService) GetCalendar(ctx context.Context, userID string, from, to time.Time) (*dto.CorporateActionCalendarResponse, error) {
	// Normalise to date-only precision
	from = truncateToDay(from)
	to = truncateToDay(to)

	// Cap range to maxCalendarRangeMonths
	maxTo := from.AddDate(0, maxCalendarRangeMonths, 0)
	if to.After(maxTo) {
		to = maxTo
	}

	fromStr := from.Format(calendarDateFormat)
	toStr := to.Format(calendarDateFormat)

	// Try cache first
	cacheKey := s.buildCacheKey(fromStr, toStr)
	if s.cache != nil {
		var cached dto.CorporateActionCalendarResponse
		if ok, _ := s.cache.GetJSON(ctx, cacheKey, &cached); ok {
			cached.Cached = true
			return &cached, nil
		}
	}

	var actions []market.CorporateAction

	// Fetch IDX corporate actions (dividend + RUPS)
	if s.idxClient != nil {
		idxActions, _ := s.idxClient.GetCorporateActions(ctx, from, to)
		actions = append(actions, idxActions...)
	}

	// Sort by date ascending
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].Date.Before(actions[j].Date)
	})

	// Convert to DTOs
	responses := make([]dto.CorporateActionResponse, 0, len(actions))
	for _, a := range actions {
		r := dto.CorporateActionResponse{
			Symbol:   a.Symbol,
			Name:     a.Name,
			Type:     string(a.Type),
			Date:     a.Date.Format(calendarDateFormat),
			Currency: a.Currency,
			Note:     a.Note,
			Market:   a.Market,
		}
		if a.PayDate != nil {
			s := a.PayDate.Format(calendarDateFormat)
			r.PayDate = &s
		}
		if a.Amount != nil {
			r.Amount = a.Amount
		}
		responses = append(responses, r)
	}

	result := &dto.CorporateActionCalendarResponse{
		From:    fromStr,
		To:      toStr,
		Total:   len(responses),
		Cached:  false,
		Actions: responses,
	}

	// Persist to cache (best-effort, ignore errors)
	if s.cache != nil {
		_ = s.cache.SetJSONWithTTL(ctx, cacheKey, result, corporateActionCacheTTL)
	}

	return result, nil
}

func (s *corporateActionService) buildCacheKey(from, to string) string {
	if s.cache == nil {
		return ""
	}
	return s.cache.BuildKey("corporate-actions-v2", from, to)
}



// truncateToDay zeroes out hours/minutes/seconds/nanoseconds in UTC.
func truncateToDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
