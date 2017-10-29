package json

import (
	ejson "encoding/json"

	"github.com/golang/glog"
	"github.com/koki/templates/generic"
	"github.com/kr/pretty"
)

// Fill a template with values from params.
func Fill(template, params []byte) ([]byte, error) {
	var err error
	t := map[string]interface{}{}
	err = ejson.Unmarshal(template, &t)
	if err != nil {
		glog.Errorf(
			"Template is not a JSON object:\n%v\n%s",
			err,
			params)
		return nil, err
	}

	p := map[string]interface{}{}
	err = ejson.Unmarshal(params, &p)
	if err != nil {
		glog.Errorf(
			"Params are not a JSON object:\n%s",
			params)
		return nil, err
	}

	filled := generic.Fill(t, p)
	var output []byte
	output, err = ejson.Marshal(filled)
	if err != nil {
		glog.Error(pretty.Sprintf(
			"Couldn't serialize result:\n%v\n%# v",
			err,
			filled))
		return nil, err
	}

	return output, nil
}
