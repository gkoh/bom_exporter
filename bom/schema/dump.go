package schema

import (
	log "github.com/sirupsen/logrus"
)

// DumpElements iterates and logs the specified Element array.
func DumpElements(elements []Element) {
	for _, e := range elements {
		log.Debugf("%s (%s): %v\n", e.Type, e.Unit, e.Value)
	}
}

// DumpTexts iterates and logs the specified Text array.
func DumpTexts(texts []Text) {
	for _, t := range texts {
		log.Debugf("%s: %v\n", t.Type, t.Value)
	}
}

// DumpPeriod iterates and logs the specified ForecastPeriod.
func DumpPeriod(period ForecastPeriod) {
	log.Debugf("%v %v", period.StartTimeUTC, period.EndTimeUTC)
	DumpElements(period.Elements)
	DumpTexts(period.Texts)
}
