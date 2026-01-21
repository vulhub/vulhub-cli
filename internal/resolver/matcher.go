package resolver

import (
	"regexp"
	"sort"
	"strings"

	"github.com/vulhub/vulhub-cli/pkg/types"
)

// CVE pattern: CVE-YYYY-XXXXX (4 digit year, 4-7 digit number)
var cvePattern = regexp.MustCompile(`^CVE-\d{4}-\d{4,7}$`)

// MatchType represents the type of match
type MatchType int

const (
	// MatchTypeNone indicates no match
	MatchTypeNone MatchType = iota
	// MatchTypeExactCVE indicates exact CVE match
	MatchTypeExactCVE
	// MatchTypeExactPath indicates exact path match
	MatchTypeExactPath
	// MatchTypeAppName indicates application name match
	MatchTypeAppName
	// MatchTypePartial indicates partial match (name, description, tags)
	MatchTypePartial
)

// Match represents a single match result
type Match struct {
	Environment types.Environment
	Type        MatchType
	Score       int
}

// Matcher provides environment matching functionality
type Matcher struct{}

// NewMatcher creates a new Matcher
func NewMatcher() *Matcher {
	return &Matcher{}
}

// IsCVEFormat checks if the keyword is in CVE format
func (m *Matcher) IsCVEFormat(keyword string) bool {
	return cvePattern.MatchString(strings.ToUpper(keyword))
}

// IsPathFormat checks if the keyword looks like a path
func (m *Matcher) IsPathFormat(keyword string) bool {
	return strings.Contains(keyword, "/")
}

// FindMatches finds all environments matching the keyword
func (m *Matcher) FindMatches(keyword string, envs []types.Environment) []Match {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil
	}

	var matches []Match

	// Normalize keyword for comparison
	keywordUpper := strings.ToUpper(keyword)
	keywordLower := strings.ToLower(keyword)

	for _, env := range envs {
		match := m.matchEnvironment(keyword, keywordUpper, keywordLower, env)
		if match.Type != MatchTypeNone {
			matches = append(matches, match)
		}
	}

	// Sort by match type (priority) and score
	sort.Slice(matches, func(i, j int) bool {
		// Lower type is better: ExactCVE=1, ExactPath=2, AppName=3, Partial=4
		if matches[i].Type != matches[j].Type {
			return matches[i].Type < matches[j].Type
		}
		// Higher score is better
		return matches[i].Score > matches[j].Score
	})

	return matches
}

// matchEnvironment checks if an environment matches the keyword
func (m *Matcher) matchEnvironment(keyword, keywordUpper, keywordLower string, env types.Environment) Match {
	// Priority 1: Exact CVE match
	if m.IsCVEFormat(keyword) {
		for _, cve := range env.CVE {
			if strings.EqualFold(cve, keyword) {
				return Match{
					Environment: env,
					Type:        MatchTypeExactCVE,
					Score:       100,
				}
			}
		}
	}

	// Priority 2: Exact path match
	if strings.EqualFold(env.Path, keyword) {
		return Match{
			Environment: env,
			Type:        MatchTypeExactPath,
			Score:       100,
		}
	}

	// Priority 3: Application name match (case-insensitive)
	if strings.EqualFold(env.App, keyword) {
		return Match{
			Environment: env,
			Type:        MatchTypeAppName,
			Score:       80,
		}
	}

	// Priority 4: Partial matches
	score := m.calculatePartialScore(keywordLower, env)
	if score > 0 {
		return Match{
			Environment: env,
			Type:        MatchTypePartial,
			Score:       score,
		}
	}

	return Match{Type: MatchTypeNone}
}

// calculatePartialScore calculates a score for partial matches
func (m *Matcher) calculatePartialScore(keywordLower string, env types.Environment) int {
	score := 0

	// Check if keyword is in path
	if strings.Contains(strings.ToLower(env.Path), keywordLower) {
		score += 50
	}

	// Check if keyword is in name
	if strings.Contains(strings.ToLower(env.Name), keywordLower) {
		score += 40
	}

	// Check if keyword is in app name
	if strings.Contains(strings.ToLower(env.App), keywordLower) {
		score += 35
	}

	// Check if keyword is in any CVE
	for _, cve := range env.CVE {
		if strings.Contains(strings.ToLower(cve), keywordLower) {
			score += 45
			break
		}
	}

	// Check if keyword is in tags
	for _, tag := range env.Tags {
		if strings.Contains(strings.ToLower(tag), keywordLower) {
			score += 30
			break
		}
	}

	return score
}

// FindExactMatch finds an exact match for the keyword (CVE or path)
func (m *Matcher) FindExactMatch(keyword string, envs []types.Environment) *types.Environment {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil
	}

	keywordUpper := strings.ToUpper(keyword)

	// Check for exact CVE match
	if m.IsCVEFormat(keyword) {
		for i := range envs {
			for _, cve := range envs[i].CVE {
				if strings.EqualFold(cve, keywordUpper) {
					return &envs[i]
				}
			}
		}
	}

	// Check for exact path match
	for i := range envs {
		if strings.EqualFold(envs[i].Path, keyword) {
			return &envs[i]
		}
	}

	return nil
}
