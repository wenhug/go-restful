package restful

import (
	"fmt"
	"regexp"
	"sync"
)

var (
	customVerbReg     = regexp.MustCompile(":([A-Za-z]+)$")
	customVerbCache   sync.Map // Cache for compiled custom verb regexes
	customVerbCacheEnabled = false // Enable/disable custom verb regex caching
)

// SetCustomVerbCacheEnabled enables or disables custom verb regex caching.
// When disabled (default), custom verb regex patterns will be compiled on every request.
// When enabled, compiled custom verb regex patterns are cached for better performance.
func SetCustomVerbCacheEnabled(enabled bool) {
	customVerbCacheEnabled = enabled
	if !enabled {
		// Clear existing cache when disabling
		customVerbCache = sync.Map{}
	}
}

func hasCustomVerb(routeToken string) bool {
	return customVerbReg.MatchString(routeToken)
}

func isMatchCustomVerb(routeToken string, pathToken string) bool {
	rs := customVerbReg.FindStringSubmatch(routeToken)
	if len(rs) < 2 {
		return false
	}

	customVerb := rs[1]
	regexPattern := fmt.Sprintf(":%s$", customVerb)

	// Check cache first (if enabled)
	if customVerbCacheEnabled {
		if cached, found := customVerbCache.Load(regexPattern); found {
			specificVerbReg := cached.(*regexp.Regexp)
			return specificVerbReg.MatchString(pathToken)
		}
	}

	// Compile the regex
	specificVerbReg := regexp.MustCompile(regexPattern)
	
	// Cache the regex (if enabled)
	if customVerbCacheEnabled {
		customVerbCache.Store(regexPattern, specificVerbReg)
	}
	
	return specificVerbReg.MatchString(pathToken)
}

func removeCustomVerb(str string) string {
	return customVerbReg.ReplaceAllString(str, "")
}
