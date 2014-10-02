package logsend

import (
	"regexp"
)

type Group struct {
	Mask  *regexp.Regexp
	Rules []*Rule
}

type Rule struct {
	regexp  *regexp.Regexp
	senders []Sender
}

func (rule *Rule) Match(line *string) interface{} {
	matches := rule.regexp.FindStringSubmatch(*line)

	if len(matches) == 0 {
		return nil
	}

	if len(matches) <= 1 {
		return true
	}

	// TODO: cache subexnames
	out := make(map[string]interface{})
	for i, value := range matches[1:] {
		key, val, err := prepareValue(rule.regexp.SubexpNames()[i+1], value)
		if err != nil {
			Conf.Logger.Printf("can't prepareValue with %+v and %+v have err %+v", rule.regexp.SubexpNames()[i+1], value, err)
			return nil
		}
		out[key] = val
	}
	if len(out) > 0 {
		return out
	}
	return nil
}

func (rule *Rule) send(data interface{}) {
	for _, sender := range rule.senders {
		sender.Send(data)
	}
}
