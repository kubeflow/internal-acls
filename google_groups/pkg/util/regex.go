package util

import (
	"regexp"
)

// StringLister is an interface to return a list of strings
type StringLister interface {
	List() ([]string, error)
}


// ArrayLister implements the Lister interface for an array.
type ArrayLister struct {
	Items []string
}

// List lists the items in the array
func (l *ArrayLister) List() ([]string, error) {
	return l.Items, nil
}

// ReMatch represents a regex match
type ReMatch struct {
	// Value is the matched value
	Value string

	// Groups is a map of the named groups if any captured from the regex
	Groups map[string]string
}

func FilterByRe(l StringLister, p *regexp.Regexp) ([]ReMatch, error){
	matches := []ReMatch{}

	items, err := l.List()

	if err != nil {
		return matches, err
	}

	for _, i := range items {
		m := p.FindStringSubmatch(i)

		if m == nil {
			continue
		}

		newMatch := ReMatch{
			Value:  m[0],
			Groups: map[string]string{},
		}
		for i, k := range p.SubexpNames() {
			// 0'th position corresponds to the whole string
			if i == 0 {
				continue
			}
			newMatch.Groups[k] = m[i]
		}

		matches = append(matches, newMatch)
	}

	return matches, nil
}