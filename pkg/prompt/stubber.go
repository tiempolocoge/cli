package prompt

import (
	"fmt"
	"reflect"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
)

type askStubber struct {
	Asks  [][]*survey.Question
	Count int
	Stubs [][]*QuestionStub
}

func InitAskStubber() (*askStubber, func()) {
	origSurveyAsk := SurveyAsk
	as := askStubber{}
	SurveyAsk = func(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
		as.Asks = append(as.Asks, qs)
		count := as.Count
		as.Count += 1
		if count >= len(as.Stubs) {
			panic(fmt.Sprintf("more asks than stubs. most recent call: %v", qs))
		}

		// actually set response
		stubbedQuestions := as.Stubs[count]
		for i, sq := range stubbedQuestions {
			q := qs[i]
			if q.Name != sq.Name {
				panic(fmt.Sprintf("stubbed question mismatch: %s != %s", q.Name, sq.Name))
			}
			if sq.Default {
				defaultValue := reflect.ValueOf(q.Prompt).Elem().FieldByName("Default")
				_ = core.WriteAnswer(response, q.Name, defaultValue)
			} else {
				_ = core.WriteAnswer(response, q.Name, sq.Value)
			}
		}

		return nil
	}
	teardown := func() {
		SurveyAsk = origSurveyAsk
	}
	return &as, teardown
}

type QuestionStub struct {
	Name    string
	Value   interface{}
	Default bool
}

func (as *askStubber) Stub(stubbedQuestions []*QuestionStub) {
	// A call to .Ask takes a list of questions; a stub is then a list of questions in the same order.
	as.Stubs = append(as.Stubs, stubbedQuestions)
}
