package creds

import "github.com/pf-qiu/concourse/v6/vars"

type String struct {
	variablesResolver vars.Variables
	rawCredString     string
}

func NewString(variables vars.Variables, credString string) String {
	return String{
		variablesResolver: variables,
		rawCredString:     credString,
	}
}

func (s String) Evaluate() (string, error) {
	var credsString string

	err := evaluate(s.variablesResolver, s.rawCredString, &credsString)
	if err != nil {
		return s.rawCredString, err
	}

	return credsString, nil
}
