package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

type Handler struct {
	URL          string
	Username     string
	Password     string
	Browser      *rod.Browser
	Page         *rod.Page
	QuestionBank []Questions
}

type Questions struct {
	QuestionNo int
	Question   string
	Options    string
}

func Sleep(timeDuration time.Duration) {
	time.Sleep(time.Second * timeDuration)
}

func NewHandler(username, password string) *Handler {
	browser := rod.New().
		ControlURL("ws://127.0.0.1:9222/devtools/browser/44f69c1f-6927-4e97-855b-fda47d55f9ec").
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

	page.MustElement("#pageid-1 > header > div > div > div.col-lg-4.col-md-4.col-sm-6.col-xs-6.header-right > ul > li:nth-child(1) > a").
		MustWaitVisible().
		MustClick()

	page.MustElement("#newcommonloginpmodal_popup > div > div > div > div > div > div:nth-child(2) > div > form > div.row > div:nth-child(1) > div > input").
		MustWaitVisible().
		MustInput(h.Username)

	page.MustElement("#newcommonloginpmodal_popup > div > div > div > div > div > div:nth-child(2) > div > form > div.row > div:nth-child(2) > div > input").
		MustInput(h.Password)

	page.MustElement("#newcommonloginpmodal_popup > div > div > div > div > div > div:nth-child(2) > div > form > div.row > div.submit-container > button").
		MustClick()

	page.MustElement("#accordionSidebar > li.nav-item.active > a > i").MustWaitVisible()
	return nil
}

func (h *Handler) NavigateToTest() error {
	page := h.Page
	page.MustNavigate("https://student.hitbullseye.com/testzone")
	page.MustWaitLoad()

	page.MustElement("#content > div > app-tests-menus > div > div.home-block.home-block1 > div > div:nth-child(1) > div > a").
		MustWaitVisible().
		MustClick()

	openTests := "#content > div > app-tests-menus > div > div.home-block.home-block1 > div > div:nth-child(6) > div > a"
	page.MustElement(openTests).MustWaitVisible().MustClick()
	page.MustElement("#content > div > app-tests-list-secure > div > div.d-sm-flex.align-items-center.justify-content-between.mb-3 > h1").
		MustWaitVisible()

	return nil
}

func (h *Handler) GiveTest() error {
	page := h.Page
	page.MustNavigate("https://student.hitbullseye.com/testzone/tests-list/Open%20Tests")

	Sleep(2)
	page.MustElement(".mk-start-btn-n-design").MustWaitVisible()

	var openTests rod.Elements
	resumeTestBtns, _ := page.Elements(".mk-resume-btn-n-design")
	openTests = append(openTests, resumeTestBtns...)

	startTestBtns, _ := page.Elements(".mk-start-btn-n-design")
	openTests = append(openTests, startTestBtns...)

	fmt.Printf("Found %d open tests\n", len(openTests))

	for i, btn := range openTests {
		fmt.Printf("Starting test %d\n", i+1)
		btn.MustClick()
		time.Sleep(time.Second * 5)
		h.startTest()
		fmt.Println("Doing next test")
		time.Sleep(time.Second * 5)
	}
	return nil
}

func (h *Handler) startTest() {
	page := h.Page
	page.MustNavigate("https://onlinetest.hitbullseye.com/online_load/directionpage.php")

	// Handle initial buttons
	deviceAndBrowserCheckNextBtn := "#nexinstructon"
	if nextBtnElement, err := page.Timeout(time.Second * 5).Element(deviceAndBrowserCheckNextBtn); err == nil {
		nextBtnElement.MustClick()
		fmt.Println("Device check button clicked")
	}

	// Wait for instructions page
	var nextBtn *rod.Element
	var err error
	for {
		nextBtn, err = page.Timeout(time.Second * 10).Element("#instPaginationa")
		if err == nil {
			break
		}

		if warningBtn, err := page.Timeout(time.Second * 5).Element("#mkOverlay > div > button"); err == nil {
			fmt.Println("Multiple session detected, closing")
			warningBtn.MustClick()
		}

		fmt.Println("Waiting for next button...")
		time.Sleep(time.Second * 2)
	}

	// Get test info
	testName, _ := page.MustElement("#firstpage > table > tbody > tr:nth-child(1) > td > table > tbody > tr > td:nth-child(2)").
		Text()
	totalQuestionsStr, _ := page.MustElement("#firstpage > table > tbody > tr:nth-child(1) > td > table > tbody > tr > td:nth-child(4)").
		Text()
	totalQuestionsInt, _ := strconv.Atoi(totalQuestionsStr)

	fmt.Printf("Test: %s | Questions: %d\n", testName, totalQuestionsInt)

	// Start test
	nextBtn.MustClick()
	page.MustElement("#disclaimer").MustWaitVisible().MustClick()
	page.MustElement("#secondPagep2 > div > input[type=submit]").MustClick()

	Sleep(3)

	// Collect questions
	h.collectQuestions(totalQuestionsInt)

	// Get AI answers and apply them
	h.applyAnswers()
}

func (h *Handler) collectQuestions(totalQuestions int) {
	fmt.Println("=== COLLECTING QUESTIONS ===")

	for i := 0; i < totalQuestions; i++ {
		questionNo := i + 1

		// Get question text
		questionText := h.getQuestionText(questionNo)

		// Get options
		options := h.getQuestionOptions(questionNo)

		// Better printing
		// fmt.Printf("\n--- Question %d ---\n", questionNo)
		// fmt.Printf("Q: %s\n", questionText)
		// fmt.Printf("Options: %v\n", options)
		fmt.Printf("Q: %d done\n", questionNo)

		// Store question
		h.QuestionBank = append(h.QuestionBank, Questions{
			QuestionNo: questionNo,
			Question:   questionText,
			Options:    strings.Join(options, " | "),
		})

		// Click next
		h.clickNext()
		time.Sleep(time.Millisecond * 100)
	}
}

func (h *Handler) applyAnswers() {
	fmt.Println("\n=== GETTING AI ANSWERS ===")
	answers := h.GetAnswers()

	fmt.Println("\n=== APPLYING ANSWERS ===")

	currentQ := 1

	for currentQ <= len(h.QuestionBank) {
		optNo, exists := answers[currentQ]
		if !exists || optNo < 1 || optNo > 4 {
			optNo = 2 // default
		}

		// Try clicking the radio button for current question
		clicked := h.clickOption(currentQ, optNo)

		if clicked {
			fmt.Printf("Q%d -> Option %d ✓\n", currentQ, optNo)
		} else {
			fmt.Printf("Q%d -> Failed to click option %d ✗\n", currentQ, optNo)
		}

		// Move to next question
		if currentQ < len(h.QuestionBank) {
			h.clickNext()
			time.Sleep(time.Millisecond * 200)
		}

		currentQ++
	}

	// Submit the test
	h.submitTest()
}

func (h *Handler) submitTest() {
	fmt.Println("\n=== SUBMITTING TEST ===")

	// Click submit button
	if submitBtn, err := h.Page.Timeout(time.Second * 5).Element("#activator"); err == nil {
		submitBtn.MustClick()
		fmt.Println("Submit button clicked ✓")

		// Wait for confirmation popup and select "Completed Test"
		time.Sleep(time.Second * 1)

		if completedRadio, err := h.Page.Timeout(time.Second * 5).Element("#box > div:nth-child(2) > div:nth-child(2) > div:nth-child(1) > input"); err == nil {
			completedRadio.MustClick()
			fmt.Println("'Completed Test' selected ✓")

			// Click OK to confirm
			if okBtn, err := h.Page.Timeout(time.Second * 3).Element("#close_confirmed"); err == nil {
				okBtn.MustClick()
				fmt.Println("Test submitted successfully ✓")
			} else {
				fmt.Println("OK button not found ✗")
			}
		} else {
			fmt.Println("Completed test radio button not found ✗")
		}
	} else {
		fmt.Println("Submit button not found ✗")
	}
}

func (h *Handler) getQuestionText(questionNo int) string {
	selectors := []string{
		fmt.Sprintf("#qno_%d > table > tbody > tr:nth-child(1) > td", questionNo),
		fmt.Sprintf("#qno_%d", questionNo),
	}

	for _, selector := range selectors {
		if element, err := h.Page.Timeout(time.Second * 2).Element(selector); err == nil {
			if text, _ := element.Text(); strings.TrimSpace(text) != "" {
				return strings.TrimSpace(text)
			}
		}
	}

	return fmt.Sprintf("Question %d text not found", questionNo)
}

func (h *Handler) getQuestionOptions(questionNo int) []string {
	selectors := []string{
		fmt.Sprintf("#answer_area_%d", questionNo),
		fmt.Sprintf(".answer_area_%d", questionNo),
	}

	for _, selector := range selectors {
		if optionsList, err := h.Page.Timeout(time.Second * 2).Element(selector); err == nil {
			if tables := optionsList.MustElements("table"); len(tables) > 0 {
				var options []string
				for _, table := range tables {
					var optionText string
					if labelElement, err := table.Element("label"); err == nil {
						optionText, _ = labelElement.Text()
					} else {
						optionText, _ = table.Text()
					}
					if text := strings.TrimSpace(optionText); text != "" {
						options = append(options, text)
					}
				}
				if len(options) > 0 {
					return options
				}
			}
		}
	}

	return []string{"A", "B", "C", "D"} // fallback
}

func (h *Handler) clickNext() {
	selectors := []string{
		"#main_div > div.tableWidthPercent > div.onlineTestLeftDiv.mk-onlineTestLeftDiv-new > div.qnav > span.saveNextButton > a",
		".saveNextButton a",
		"#saveNext",
		"input[value='Save & Next']",
	}

	for _, selector := range selectors {
		if element, err := h.Page.Timeout(time.Second * 2).Element(selector); err == nil {
			element.MustClick()
			return
		}
	}
}

func (h *Handler) clickOption(questionNo, optionNo int) bool {
	// Based on your HTML, the structure is answer_area_X > table:nth-child(Y) > tbody > tr > td:nth-child(1) > input
	selectors := []string{
		// Direct radio input targeting
		fmt.Sprintf(
			"#answer_area_%d > table:nth-child(%d) > tbody > tr > td:nth-child(1) > input",
			questionNo,
			optionNo,
		),
		fmt.Sprintf(
			"#answer_area_%d > table:nth-child(%d) input[type='radio']",
			questionNo,
			optionNo,
		),
		fmt.Sprintf(
			"#answer_area_%d input[value='%d']",
			questionNo,
			optionNo-1,
		), // sometimes 0-indexed
		fmt.Sprintf("#answer_area_%d input[value='%d']", questionNo, optionNo),
		// Generic radio button in answer area
		fmt.Sprintf("#answer_area_%d input[type='radio']:nth-of-type(%d)", questionNo, optionNo),
	}

	for _, selector := range selectors {
		if element, err := h.Page.Timeout(time.Millisecond * 500).Element(selector); err == nil {
			element.MustClick()
			return true
		}
	}

	return false
}
