package app_profile

import (
	"os"
	"strings"
)

// GetProfileByScopeSuffix retrieves the app_profile by coincidence in the scope suffix
// this is an adaptation of the ScopeUtils.java
func GetProfileByScope() string {
	tokens := strings.Split(GetScopeValue(), "-")
	return tokens[len(tokens)-1]
}

// IsLocalProfile retrieves information about if the app_profile is local or not
// this is an adaptation of the ScopeUtils.java
func IsLocalProfile() bool {
	return Local == GetScopeValue()
}

// IsTestProfile retrieves information about if the app_profile is test or not
// this is an adaptation of the ScopeUtils.java
func IsTestProfile() bool {
	return strings.HasSuffix(GetScopeValue(), Test)
}

// IsProdProfile retrieves information about if the app_profile is prod or not
// this is an adaptation of the ScopeUtils.java
func IsProdProfile() bool {
	return strings.HasSuffix(GetScopeValue(), Prod)
}

// IsStageProfile retrieves information about if the app_profile is stage or not
// this is an adaptation of the ScopeUtils.java
func IsStageProfile() bool {
	return strings.HasSuffix(GetScopeValue(), Stage)
}

func GetScopeValue() string {
	scope := os.Getenv("SCOPE")
	if scope != "" {
		return scope
	}
	return Local
}
