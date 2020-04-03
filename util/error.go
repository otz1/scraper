package util

import (
	"fmt"
	"github.com/otz1/scraper/entity"
)

func InvalidSiteCodeErr(siteCode entity.SiteCode) error {
	return fmt.Errorf("invalid site code '%s'", siteCode)
}

func InvalidSiteCodeHeaderErr(siteCode string) error {
	return fmt.Errorf("invalid site code '%s'", siteCode)
}