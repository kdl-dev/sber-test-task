package parser

import (
	"bytes"
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

const (
	TAG_BODY   = "body"
	TAG_A      = "a"
	TAG_SELECT = "select"
	TAG_OPTION = "option"
	TAG_INPUT  = "input"
	TAG_H1     = "h1"

	ATTR_HREF  = "href"
	ATTR_VALUE = "value"
	ATTR_NAME  = "name"
	ATTR_TYPE  = "type"

	INPUT_TYPE_RADIO = "radio"
	INPUT_TYPE_TEXT  = "text"
)

type UserInput struct {
	Name    string
	Options []Option
}

type Option struct {
	Value   string
	Content string
}

type UserInputs struct {
	TextInput   []UserInput
	SelectInput []UserInput
	RadioInput  []UserInput
}

func GetPathToSartTest(body io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err
	}

	urlPath, isExists := doc.Find(TAG_BODY).Find(TAG_A).Attr(ATTR_HREF)
	if !isExists {
		return "", fmt.Errorf("href attribute for tag 'a' not found")
	}

	return urlPath, nil
}

func ParseSelects(body io.Reader) ([]UserInput, error) {
	var selectInputs []UserInput

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	doc.Find(TAG_SELECT).Each(func(i int, s *goquery.Selection) {
		var opts []Option
		s.Find(TAG_OPTION).Each(func(j int, s *goquery.Selection) {
			opts = append(opts, Option{})
			opts[j].Value, _ = s.Attr(ATTR_VALUE)
			opts[j].Content = s.Text()
		})

		// Add info about select
		selectInputs = append(selectInputs, UserInput{})
		selectInputs[i].Name, _ = s.Attr(ATTR_NAME)
		selectInputs[i].Options = append(selectInputs[i].Options, opts...)
	})

	return selectInputs, nil
}

func ParseRadios(body io.Reader) ([]UserInput, error) {
	var radioInputs []UserInput

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	index := -1
	doc.Find(TAG_INPUT).Each(func(i int, s *goquery.Selection) {

		if inputType, _ := s.Attr(ATTR_TYPE); inputType != INPUT_TYPE_RADIO {
			return
		}

		inputRadioName, _ := s.Attr(ATTR_NAME)

		// If this is the first counter element "radio" (index == -1)
		// or the current parsed name is not the name of the previous added
		// element "radio", then need to add a new element "radio" to the slice.
		if index == -1 || inputRadioName != radioInputs[index].Name {
			radioInputs = append(radioInputs, UserInput{Name: inputRadioName})
			index++
		}

		inputRadioValue, _ := s.Attr(ATTR_VALUE)
		radioInputs[index].Options = append(radioInputs[index].Options, Option{
			Value:   inputRadioValue,
			Content: inputRadioValue,
		})
	})

	return radioInputs, nil
}

func ParseTextField(body io.Reader) ([]UserInput, error) {
	var textInputs []UserInput

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	doc.Find(TAG_INPUT).Each(func(i int, s *goquery.Selection) {
		if inputType, _ := s.Attr(ATTR_TYPE); inputType != INPUT_TYPE_TEXT {
			return
		}

		name, _ := s.Attr(ATTR_NAME)

		textInputs = append(textInputs, UserInput{
			Name: name,
		})
	})

	return textInputs, nil
}

func ParseUserInputs(body []byte) (*UserInputs, error) {
	textFields, err := ParseTextField(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	selects, err := ParseSelects(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	radios, err := ParseRadios(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return &UserInputs{
		TextInput:   textFields,
		SelectInput: selects,
		RadioInput:  radios,
	}, nil
}

// Parse last html (Test result)
func ParseSuccessMsg(body []byte) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	return doc.Find(TAG_BODY).Find(TAG_H1).Text(), nil
}
