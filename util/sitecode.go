package util

import (
	"github.com/getsentry/sentry-go"
	"github.com/otz1/scraper/entity"
	"log"
)

var siteCodeMap = map[string]entity.SiteCode {
	"OTZIT_UK": entity.OTZIT_UK,
	"OTZIT_US": entity.OTZIT_US,
	"OTZIT_FR": entity.OTZIT_FR,
	"OTZIT_IT": entity.OTZIT_IT,
	"OTZIT_ES": entity.OTZIT_ES,
}

func GetSiteCode(siteCodeHeader string) (entity.SiteCode) {
	log.Println("Parsing siteCode", siteCodeHeader)

	siteCode, ok := siteCodeMap[siteCodeHeader]
	if !ok {
		err := InvalidSiteCodeHeaderErr(siteCodeHeader)
		sentry.CaptureException(err)
		panic(err)
	}
	return siteCode
}