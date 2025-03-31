// customdetections.go
package detection

import (
	"regexp"
)

func (d *Detector) ReloadCustomPatterns() {
	d.customPatterns = nil
	for _, cd := range config.customDetections {
		if cd.Type == "regex" {
			re, err := regexp.Compile(cd.Pattern)
			if err == nil {
				d.customPatterns = append(d.customPatterns, CustomPattern{
					Pattern:     re,
					EventType:   EventType(cd.EventType),
					MessageTmpl: cd.Message,
					IsRegex:     true,
				})
			}
		} else {
			d.customPatterns = append(d.customPatterns, CustomPattern{
				Keyword:     cd.Pattern,
				EventType:   EventType(cd.EventType),
				MessageTmpl: cd.Message,
				IsRegex:     false,
			})
		}
	}
}
