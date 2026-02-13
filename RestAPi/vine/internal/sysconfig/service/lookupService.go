package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/Deo-Mugabe/GOLANG/internal/sysconfig/domain"
	"github.com/rs/zerolog"
)

// LookupService manages system lookups with caching
type LookupService struct {
	lookupRepo domain.LookupRepository
	cache      map[string]*domain.SystemLookup
	cacheMu    sync.RWMutex
	logger     zerolog.Logger
}

// NewLookupService creates a new lookup service
func NewLookupService(
	lookupRepo domain.LookupRepository,
	logger zerolog.Logger,
) *LookupService {
	return &LookupService{
		lookupRepo: lookupRepo,
		cache:      make(map[string]*domain.SystemLookup),
		logger:     logger,
	}
}

// GetLookup retrieves a lookup value
func (s *LookupService) GetLookup(ctx context.Context, codeAgency, codeKey string) (*domain.SystemLookup, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("%s:%s", codeAgency, codeKey)

	s.cacheMu.RLock()
	if cached, ok := s.cache[cacheKey]; ok {
		s.cacheMu.RUnlock()
		s.logger.Debug().
			Str("code_agency", codeAgency).
			Str("code_key", codeKey).
			Msg("lookup retrieved from cache")
		return cached, nil
	}
	s.cacheMu.RUnlock()

	// Fetch from database
	lookup, err := s.lookupRepo.GetByKey(ctx, codeAgency, codeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get lookup: %w", err)
	}

	// Cache the result
	s.cacheMu.Lock()
	s.cache[cacheKey] = lookup
	s.cacheMu.Unlock()

	s.logger.Debug().
		Str("code_agency", codeAgency).
		Str("code_key", codeKey).
		Str("message", lookup.Message).
		Msg("lookup retrieved and cached")

	return lookup, nil
}

// GetAgencyLookups retrieves all lookups for an agency
func (s *LookupService) GetAgencyLookups(ctx context.Context, codeAgency string) ([]*domain.SystemLookup, error) {
	lookups, err := s.lookupRepo.GetByAgency(ctx, codeAgency)
	if err != nil {
		return nil, fmt.Errorf("failed to get agency lookups: %w", err)
	}

	// Cache all lookups
	s.cacheMu.Lock()
	for _, lookup := range lookups {
		cacheKey := fmt.Sprintf("%s:%s", lookup.CodeAgency, lookup.CodeKey)
		s.cache[cacheKey] = lookup
	}
	s.cacheMu.Unlock()

	s.logger.Debug().
		Str("code_agency", codeAgency).
		Int("count", len(lookups)).
		Msg("agency lookups retrieved and cached")

	return lookups, nil
}

// TranslateCode translates a code to its message
func (s *LookupService) TranslateCode(ctx context.Context, codeAgency, codeKey string) (string, error) {
	lookup, err := s.GetLookup(ctx, codeAgency, codeKey)
	if err != nil {
		// Return empty string if lookup not found (graceful degradation)
		s.logger.Warn().
			Err(err).
			Str("code_agency", codeAgency).
			Str("code_key", codeKey).
			Msg("lookup not found, returning empty string")
		return "", nil
	}

	return lookup.Message, nil
}

// InvalidateCache clears the lookup cache
func (s *LookupService) InvalidateCache(ctx context.Context) error {
	s.cacheMu.Lock()
	s.cache = make(map[string]*domain.SystemLookup)
	s.cacheMu.Unlock()

	s.logger.Info().Msg("lookup cache invalidated")
	return nil
}

// PreloadCache preloads all lookups into cache
func (s *LookupService) PreloadCache(ctx context.Context) error {
	s.logger.Info().Msg("preloading lookup cache")

	// Get all lookups (with reasonable limit)
	lookups, err := s.lookupRepo.List(ctx, 1000, 0)
	if err != nil {
		return fmt.Errorf("failed to preload cache: %w", err)
	}

	// Cache all lookups
	s.cacheMu.Lock()
	for _, lookup := range lookups {
		cacheKey := fmt.Sprintf("%s:%s", lookup.CodeAgency, lookup.CodeKey)
		s.cache[cacheKey] = lookup
	}
	s.cacheMu.Unlock()

	s.logger.Info().
		Int("count", len(lookups)).
		Msg("lookup cache preloaded")

	return nil
}

// CreateLookup creates a new lookup
func (s *LookupService) CreateLookup(ctx context.Context, lookup *domain.SystemLookup) error {
	if err := s.lookupRepo.Create(ctx, lookup); err != nil {
		return fmt.Errorf("failed to create lookup: %w", err)
	}

	// Add to cache
	cacheKey := fmt.Sprintf("%s:%s", lookup.CodeAgency, lookup.CodeKey)
	s.cacheMu.Lock()
	s.cache[cacheKey] = lookup
	s.cacheMu.Unlock()

	s.logger.Info().
		Str("code_agency", lookup.CodeAgency).
		Str("code_key", lookup.CodeKey).
		Msg("lookup created")

	return nil
}

// UpdateLookup updates a lookup
func (s *LookupService) UpdateLookup(ctx context.Context, lookup *domain.SystemLookup) error {
	if err := s.lookupRepo.Update(ctx, lookup); err != nil {
		return fmt.Errorf("failed to update lookup: %w", err)
	}

	// Update cache
	cacheKey := fmt.Sprintf("%s:%s", lookup.CodeAgency, lookup.CodeKey)
	s.cacheMu.Lock()
	s.cache[cacheKey] = lookup
	s.cacheMu.Unlock()

	s.logger.Info().
		Str("code_agency", lookup.CodeAgency).
		Str("code_key", lookup.CodeKey).
		Msg("lookup updated")

	return nil
}

// DeleteLookup deletes a lookup
func (s *LookupService) DeleteLookup(ctx context.Context, id int) error {
	// Get lookup first to find cache key
	lookup, err := s.lookupRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get lookup: %w", err)
	}

	if err := s.lookupRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete lookup: %w", err)
	}

	// Remove from cache
	cacheKey := fmt.Sprintf("%s:%s", lookup.CodeAgency, lookup.CodeKey)
	s.cacheMu.Lock()
	delete(s.cache, cacheKey)
	s.cacheMu.Unlock()

	s.logger.Info().
		Int("id", id).
		Str("code_agency", lookup.CodeAgency).
		Str("code_key", lookup.CodeKey).
		Msg("lookup deleted")

	return nil
}
