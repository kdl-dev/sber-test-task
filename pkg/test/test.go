package test

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kdl-dev/sber-test-task/pkg/global"
	"github.com/kdl-dev/sber-test-task/pkg/parser"
	"github.com/kdl-dev/sber-test-task/pkg/web"
)

type Test struct {
	sid             *http.Cookie
	host            string
	questionURLPath string
}

func (t *Test) SID() *http.Cookie {
	return t.sid
}

// return Test object, test start URL, error
func NewTest(addr string) (*Test, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return nil, err
	}

	var newTest Test

	newTest.host = addr
	newTest.sid = resp.Cookies()[0]

	newTest.questionURLPath, err = parser.GetPathToSartTest(resp.Body)
	if err != nil {
		return nil, err
	}

	return &newTest, nil
}

func (t *Test) SolveTest() (string, error) {
	resp, err := web.SendRequest("Get", t.host+t.questionURLPath, &web.HttpHeaders{
		Cookie: t.sid,
	}, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body []byte

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var QuestionIndex int

	for {
		answers, err := getRightAnswers(body)
		if err != nil {
			return "", err
		}

		if answers == nil && err == nil {
			return parser.ParseSuccessMsg(body)
		}

		QuestionIndex++
		log.Printf("Question #%d SUCCESS\n", QuestionIndex)
		global.PrintVerboseInfo(global.CLI_Border)

		resp, err = web.SendRequest("POST", t.host+t.questionURLPath, &web.HttpHeaders{
			Cookie: t.sid, ContentType: web.ContentType}, answers)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		t.questionURLPath = resp.Request.URL.Path

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
	}
}

func getRightAnswers(body []byte) (io.Reader, error) {
	userInputs, err := parser.ParseUserInputs(body)
	if err != nil {
		return nil, err
	}

	httpBody := getLongestValueFromUserInputs(*userInputs)
	if httpBody == nil {
		return nil, nil
	}

	data := url.Values{}
	for i := 0; i < len(httpBody); i++ {
		data.Set(httpBody[i].Key, httpBody[i].Value)
	}

	return strings.NewReader(data.Encode()), nil
}

func getLongestValueFromUserInputs(userInputs parser.UserInputs) []web.HttpBody {
	var httpBody []web.HttpBody

	texts := userInputs.TextInput
	if texts != nil {
		global.PrintVerboseInfo("Text input type:\n")
		httpBody = append(httpBody, getLongestValueFromUserInput(texts)...)
	}

	selects := userInputs.SelectInput
	if selects != nil {
		global.PrintVerboseInfo("Select input type:\n")
		httpBody = append(httpBody, getLongestValueFromUserInput(selects)...)
	}

	radios := userInputs.RadioInput
	if radios != nil {
		global.PrintVerboseInfo("Radios input type:\n")
		httpBody = append(httpBody, getLongestValueFromUserInput(radios)...)
	}

	// If the parser did not find any input element
	if texts == nil && selects == nil && radios == nil {
		return nil
	}

	return httpBody
}

func getLongestValueFromUserInput(userInput []parser.UserInput) []web.HttpBody {
	var httpBody []web.HttpBody
	var longer *parser.Option
	var value string

	for _, input := range userInput {
		printUserInput(&input)

		// For selects and radios
		if input.Options != nil {
			longer = calcLongestValue(input.Options)
			value = longer.Value
		} else {
			value = os.Getenv("TEXT_WORD") // For text field
		}

		httpBody = append(httpBody, web.HttpBody{
			Key:   input.Name,
			Value: value,
		})
	}

	return httpBody
}

func calcLongestValue(arr []parser.Option) *parser.Option {
	if arr == nil {
		return nil
	}

	longer := arr[0]

	for i := 1; i < len(arr); i++ {
		if len(arr[i].Content) > len(longer.Content) {
			longer = arr[i]
		}
	}

	return &longer
}

func printUserInput(userInput *parser.UserInput) {
	if userInput == nil {
		return
	}

	global.PrintVerboseInfo("Name: %v\n", userInput.Name)

	if userInput.Options == nil {
		global.PrintVerboseInfo("\n")
		return
	}

	for j, opt := range userInput.Options {
		global.PrintVerboseInfo("%d. Value: %v\n", j+1, opt.Value)
	}
	global.PrintVerboseInfo("\n")
}
