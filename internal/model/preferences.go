package model

import "strings"

type PreferenceSet struct {
	Values []string
}

func NewPreferenceSet(values ...string) PreferenceSet {
	set := PreferenceSet{Values: make([]string, 0, len(values))}
	for _, value := range values {
		if PreferenceAllowed(value) && !set.Contains(value) {
			set.Values = append(set.Values, value)
		}
	}
	return set
}

func (p PreferenceSet) Contains(value string) bool {
	for _, candidate := range p.Values {
		if candidate == value {
			return true
		}
	}
	return false
}

func (p PreferenceSet) Merge(other PreferenceSet) PreferenceSet {
	merged := NewPreferenceSet(p.Values...)
	for _, value := range other.Values {
		if PreferenceAllowed(value) && !merged.Contains(value) {
			merged.Values = append(merged.Values, value)
		}
	}
	return merged
}

func (p PreferenceSet) Missing(required []string) []string {
	missing := make([]string, 0)
	for _, value := range required {
		if !p.Contains(value) {
			missing = append(missing, value)
		}
	}
	return missing
}

func (p PreferenceSet) String() string { return strings.Join(p.Values, ", ") }

func ComparePreferences(left, right PreferenceSet) (shared, onlyLeft, onlyRight []string) {
	for _, value := range left.Values {
		if right.Contains(value) {
			shared = append(shared, value)
		} else {
			onlyLeft = append(onlyLeft, value)
		}
	}
	for _, value := range right.Values {
		if !left.Contains(value) {
			onlyRight = append(onlyRight, value)
		}
	}
	return shared, onlyLeft, onlyRight
}
