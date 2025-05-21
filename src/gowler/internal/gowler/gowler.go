package gowler

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/temoto/robotstxt"
)

type Gowler struct {
	Site string
	Domain string
	SiteUrls []string
	OtherUrls []string
	Telephones []string
	Emails []string
	UserAgents string
	Timeout time.Duration

	wg *sync.WaitGroup
	sem chan struct{}
	mu   sync.Mutex
	robotsData robotstxt.RobotsData
}

// NewGowler creates a new Gowler struct
// site is the site to scrape
func NewGowler(site string) *Gowler {
	u, err := url.Parse(site)
	if err != nil {
		log.Fatal(err)
	}
	domain := u.Hostname()
	return &Gowler{
		Site: site,
		Domain: domain,
		SiteUrls: []string{site},
		UserAgents: "Gowler",
		Timeout: 60 * time.Second,
		wg: &sync.WaitGroup{},
		sem: make(chan struct{}, 10),
	}
}

// Crawl starts the crawling process
func (entity *Gowler) Crawl() {
	_, err := entity.getRobotData()
	if err != nil {
		log.Println(ColorRed, err, ColorReset)
	}
	entity.wg.Add(1)
	go entity.gowler(entity.Site)
	entity.wg.Wait()
}

// gowler is the recursive function that scrapes the site
// urlToScrape is the URL to scrape
func (entity *Gowler) gowler(urlToScrape string) {
	defer entity.wg.Done()
	entity.sem <- struct{}{}
	defer func() {
		<-entity.sem
	}()

	log.Println(ColorGreen, "Scraping site:", urlToScrape, ColorReset)
	req, err := http.NewRequest("GET", urlToScrape, nil)
	if err != nil {
		log.Println(ColorRed, err, ColorReset)
		return
	}
	req.Header.Set("User-Agent", entity.UserAgents)
	
	client := &http.Client{
		Timeout: entity.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(ColorRed, err, ColorReset)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Println(ColorRed, err, ColorReset)
			return
		}

		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			link, exists := s.Attr("href")
			if exists {
				if IsTelephone(link) && !siteInSlice(link, entity.Telephones) {
					telephone := link
					if link[:4] == "tel:" {
						telephone = link[4:]
					}
					entity.Telephones = append(entity.Telephones, telephone)
				} else if IsEmail(link) && !siteInSlice(link, entity.Emails) {
					email := link
					if link[:7] == "mailto:" {
						email = link[7:]
					}
					entity.Emails = append(entity.Emails, email)
				} else {
					u, err := url.Parse(link)
					if err != nil {
						log.Println(ColorRed, err, ColorReset)
						return
					}
					if !u.IsAbs() {
						u = resp.Request.URL.ResolveReference(u)
						log.Println(ColorYellow, "Relative URL resolved to:", u.String(), ColorReset)
					}
					if !entity.robotsData.TestAgent(entity.UserAgents, u.Path) {
						log.Println(ColorYellow, "URL disallowed by robots.txt:", u.String(), ColorReset)
						return
					}
					domain := u.Hostname()
					if domain == entity.Domain && !siteInSlice(link, entity.SiteUrls) {
						entity.addSiteUrl(link)
						entity.wg.Add(1)
						go entity.gowler(link)
						log.Println(ColorBlue, "Active Goroutines:", len(entity.sem), ColorReset)
					} else if domain != entity.Domain && !siteInSlice(link, entity.OtherUrls) {
						entity.addOtherUrl(link)
					}
				}
			}
		})

		doc.Find("body").Each(func(i int, s *goquery.Selection) {
			html, err := s.Html()
			if err != nil {
				log.Println(ColorRed, err, ColorReset)
				return
			}
			// Email regex
			emailRegex := `[\w._%+-]+@[\w.-]+\.[a-zA-Z]{2,}`
			// Phone regex (basic support for various formats)
			phoneRegex := `(\+?\d{1,3}[-.\s]?)?\(?\d{2,4}\)?[-.\s]?\d{3,4}[-.\s]?\d{4}`

			// Compile regex
			emailRe := regexp.MustCompile(emailRegex)
			phoneRe := regexp.MustCompile(phoneRegex)

			// Find matches
			emails := emailRe.FindAllString(html, -1)
			phones := phoneRe.FindAllString(html, -1)

			for _, email := range emails {
				if !siteInSlice(email, entity.Emails) {
					entity.Emails = append(entity.Emails, email)
				}
			}
			for _, phone := range phones {
				if !siteInSlice(phone, entity.Telephones) {
					entity.Telephones = append(entity.Telephones, phone)
				}
			}
		})
	} else {
		log.Printf(ColorYellow + "Site %s returned status code %d" +  ColorReset, urlToScrape, resp.StatusCode)
	}
}

// siteInSlice checks if a string is in a slice
// returns true if the string is in the slice
func siteInSlice(site string, sites []string) bool {
	for _, s := range sites {
		if s == site {
			return true
		}
	}
	return false
}

// addSiteUrl adds a site URL to the Gowler struct
func (entity *Gowler) addSiteUrl(url string) {
	entity.mu.Lock()
	defer entity.mu.Unlock()
	entity.SiteUrls = append(entity.SiteUrls, url)
}

// addOtherUrl adds a non-site URL to the Gowler struct
func (entity *Gowler) addOtherUrl(url string) {
	entity.mu.Lock()
	defer entity.mu.Unlock()
	entity.OtherUrls = append(entity.OtherUrls, url)
}

// getRobotData gets the robots.txt data
func (entity *Gowler) getRobotData() (robotstxt.RobotsData, error) {
	robotsUrl := entity.Site + "/robots.txt"
	req, err := http.NewRequest("GET", robotsUrl, nil)
	if err != nil {
		return robotstxt.RobotsData{}, err
	}
	req.Header.Set("User-Agent", entity.UserAgents)
	
	client := &http.Client{
		Timeout: entity.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return robotstxt.RobotsData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		robots, err := robotstxt.FromResponse(resp)
		if err != nil {
			return robotstxt.RobotsData{}, err
		}
		return *robots, nil
	}

	return robotstxt.RobotsData{}, nil
}

func IsTelephone(link string) bool {
	// Check if the link starts with "tel:"
	return len(link) > 4 && link[:4] == "tel:"
}

func IsEmail(link string) bool {
	// Check if the link starts with "mailto:"
	return len(link) > 7 && link[:7] == "mailto:"
}