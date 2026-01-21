package resolver

import (
	"context"
	"fmt"

	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// ResolveResult represents the result of keyword resolution
type ResolveResult struct {
	// Keyword is the original keyword
	Keyword string

	// MatchType is the type of match found
	MatchType MatchType

	// Environment is the resolved environment (if exactly one match)
	Environment *types.Environment

	// Matches contains all matching environments (if multiple matches)
	Matches []Match

	// ExactMatch indicates if this was an exact match (CVE or path)
	ExactMatch bool
}

// Resolver defines the interface for keyword resolution
type Resolver interface {
	// Resolve resolves a keyword to one or more environments
	Resolve(ctx context.Context, keyword string) (*ResolveResult, error)

	// IsCVEFormat checks if keyword is in CVE format
	IsCVEFormat(keyword string) bool

	// IsPathFormat checks if keyword looks like a path
	IsPathFormat(keyword string) bool
}

// EnvironmentResolver implements the Resolver interface
type EnvironmentResolver struct {
	configMgr config.Manager
	matcher   *Matcher
}

// NewEnvironmentResolver creates a new EnvironmentResolver
func NewEnvironmentResolver(configMgr config.Manager) *EnvironmentResolver {
	return &EnvironmentResolver{
		configMgr: configMgr,
		matcher:   NewMatcher(),
	}
}

// Resolve resolves a keyword to one or more environments
func (r *EnvironmentResolver) Resolve(ctx context.Context, keyword string) (*ResolveResult, error) {
	// Load environment list
	envList, err := r.configMgr.LoadEnvironments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load environments: %w", err)
	}

	result := &ResolveResult{
		Keyword: keyword,
	}

	// First try exact match (CVE or path)
	if exactEnv := r.matcher.FindExactMatch(keyword, envList.Environment); exactEnv != nil {
		result.Environment = exactEnv
		result.ExactMatch = true
		if r.matcher.IsCVEFormat(keyword) {
			result.MatchType = MatchTypeExactCVE
		} else {
			result.MatchType = MatchTypeExactPath
		}
		return result, nil
	}

	// Find all matches
	matches := r.matcher.FindMatches(keyword, envList.Environment)

	if len(matches) == 0 {
		result.MatchType = MatchTypeNone
		return result, nil
	}

	if len(matches) == 1 {
		result.Environment = &matches[0].Environment
		result.MatchType = matches[0].Type
		result.Matches = matches
		return result, nil
	}

	// Multiple matches
	result.MatchType = matches[0].Type
	result.Matches = matches
	return result, nil
}

// IsCVEFormat checks if keyword is in CVE format
func (r *EnvironmentResolver) IsCVEFormat(keyword string) bool {
	return r.matcher.IsCVEFormat(keyword)
}

// IsPathFormat checks if keyword looks like a path
func (r *EnvironmentResolver) IsPathFormat(keyword string) bool {
	return r.matcher.IsPathFormat(keyword)
}

// HasMultipleMatches checks if the result has multiple matches
func (r *ResolveResult) HasMultipleMatches() bool {
	return len(r.Matches) > 1
}

// HasNoMatches checks if the result has no matches
func (r *ResolveResult) HasNoMatches() bool {
	return r.MatchType == MatchTypeNone
}

// HasSingleMatch checks if the result has exactly one match
func (r *ResolveResult) HasSingleMatch() bool {
	return r.Environment != nil
}

// GetMatchedEnvironments returns all matched environments
func (r *ResolveResult) GetMatchedEnvironments() []types.Environment {
	// If exact match and no matches slice, return the single environment
	if r.Environment != nil && len(r.Matches) == 0 {
		return []types.Environment{*r.Environment}
	}

	envs := make([]types.Environment, len(r.Matches))
	for i, m := range r.Matches {
		envs[i] = m.Environment
	}
	return envs
}
