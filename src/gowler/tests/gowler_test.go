package tests

import (
	"fmt"
	"testing"

	"gowler/internal/gowler"
)

func TestGetRobotData(t *testing.T) {
	g := gowler.NewGowler("https://www.scponline.it")

	g.Crawl()
	fmt.Println("Site URLs:", g.SiteUrls)
	fmt.Println("Other URLs:", g.OtherUrls)
	fmt.Println("Telephones:", g.Telephones)
	fmt.Println("Emails:", g.Emails)
	fmt.Println("User Agents:", g.UserAgents)
	fmt.Println("Timeout:", g.Timeout)
	fmt.Println("Domain:", g.Domain)
	fmt.Println("Site:", g.Site)

}