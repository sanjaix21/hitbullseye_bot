package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

type Handler struct {
	URL      string
	Username string
	Password string
	Browser  *rod.Browser
	Page     *rod.Page
}

func NewHandler(username, password string) *Handler {
	// browser := rod.New().MustConnect()
	browser := rod.New().
		ControlURL("ws://127.0.0.1:9222/devtools/browser/d31f2563-f4d3-4cc4-97c2-23efed03d77d").
		MustConnect()
	return &Handler{
		URL:      "https://www.hitbullseye.com/",
		Username: username,
		Password: password,
		Browser:  browser,
		Page:     browser.MustPage(),
	}
}

func (h *Handler) Login() error {
	page := h.Page
	page.MustNavigate(h.URL).MustWaitLoad()
	// page.MustWindowFullscreen()

	// Clicking login button
	page.MustElement("#pageid-1 > header > div > div > div.col-lg-4.col-md-4.col-sm-6.col-xs-6.header-right > ul > li:nth-child(1) > a").
		MustWaitVisible().
		MustClick()

	// Entering username/id
	page.MustElement("#newcommonloginpmodal_popup > div > div > div > div > div > div:nth-child(2) > div > form > div.row > div:nth-child(1) > div > input").
		MustWaitVisible().
		MustInput(h.Username)

	// Entering password
	page.MustElement("#newcommonloginpmodal_popup > div > div > div > div > div > div:nth-child(2) > div > form > div.row > div:nth-child(2) > div > input").
		MustInput(h.Password)

	// Clicking on submit button to login into the site
	page.MustElement("#newcommonloginpmodal_popup > div > div > div > div > div > div:nth-child(2) > div > form > div.row > div.submit-container > button").
		MustClick()

	// Making sure the site is loaded by checking if dashboard is visible
	page.MustElement("#accordionSidebar > li.nav-item.active > a > i").MustWaitVisible()
	// for debug/development
	return nil
}

func (h *Handler) NavigateToTest() error {
	// clicking on tests tab
	page := h.Page
	page.MustNavigate("https://student.hitbullseye.com/testzone")
	page.MustWaitLoad()

	// Opening the available tests
	page.MustElement("#content > div > app-tests-menus > div > div.home-block.home-block1 > div > div:nth-child(1) > div > a").
		MustWaitVisible().
		MustClick()

	// Navigating to Open Tests
	openTests := "#content > div > app-tests-menus > div > div.home-block.home-block1 > div > div:nth-child(6) > div > a"

	page.MustElement(openTests).MustWaitVisible().MustClick()
	page.MustElement("#content > div > app-tests-list-secure > div > div.d-sm-flex.align-items-center.justify-content-between.mb-3 > h1").
		MustWaitVisible()

	return nil
}

func (h *Handler) GiveTest() error {
	page := h.Page
	// Making sure all the tests load before couting them
	page.MustElementR("td.details1", ".*").MustWaitVisible()

	tableSelector := "#content > div > app-tests-list-secure > div > div.test-screen-table.test-screen-table1.mk-test-screen-table > div.container > div:nth-child(3) > div > div > table"

	table := page.MustElement(tableSelector)
	rows := table.MustElements("tr")
	fmt.Println("Founded", len(rows), "rows")

	for i, row := range rows {
		statusCol, err := row.Timeout(5*time.Second).ElementR(`td.details1`, `.*`)
		if err != nil {
			continue
		}

		statusText := statusCol.MustText()
		if strings.Contains(statusText, "Start Now") {
			fmt.Println("Found Open Test at", i+1)
			h.startTest(statusCol)
		}
	}
	return nil
}

func (h *Handler) startTest(startBtn *rod.Element) {
	// page := h.Page
	startBtn.MustClick()
}
