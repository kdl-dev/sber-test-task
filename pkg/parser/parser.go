package parser

import (
	"bytes"
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
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

	path, isExists := doc.Find("body").Find("a").Attr("href")
	if !isExists {
		return "", fmt.Errorf("href attribute for tag a not found")
	}

	return path, nil
}

func ParseSelects(body io.Reader) ([]UserInput, error) {
	var selectInputs []UserInput

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	doc.Find("select").Each(func(i int, s *goquery.Selection) {
		var opts []Option
		s.Find("option").Each(func(j int, s *goquery.Selection) {
			opts = append(opts, Option{})
			opts[j].Value, _ = s.Attr("value")
			opts[j].Content = s.Text()
		})
		selectInputs = append(selectInputs, UserInput{})
		selectInputs[i].Name, _ = s.Attr("name")

		selectInputs[i].Options = append(selectInputs[i].Options, opts...)
	})

	return selectInputs, nil
}

func ParseRadios(body io.Reader) ([]UserInput, error) {
	var selectInputs []UserInput

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	index := -1
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		inputType, _ := s.Attr("type")
		if inputType == "radio" {
			inputRadioName, _ := s.Attr("name")

			if index == -1 || inputRadioName != selectInputs[index].Name {
				selectInputs = append(selectInputs, UserInput{Name: inputRadioName})
				index++
			}

			inputRadioValue, _ := s.Attr("value")
			selectInputs[index].Options = append(selectInputs[index].Options, Option{
				Value:   inputRadioValue,
				Content: inputRadioValue,
			})
		}
	})

	return selectInputs, nil
}

func ParseTextField(body io.Reader) ([]UserInput, error) {
	var selectInputs []UserInput

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		if inputType, _ := s.Attr("type"); inputType == "text" {
			name, _ := s.Attr("name")

			selectInputs = append(selectInputs, UserInput{
				Name: name,
			})
		}
	})

	return selectInputs, nil
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

func ParseSuccessMsg(body []byte) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	return doc.Find("body").Find("h1").Text(), nil

}
