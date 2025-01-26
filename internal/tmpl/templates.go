package tmpl

import (
	"bytes"
	_ "embed"
	"text/template"
)

type MessageParams struct {
	Type     string
	Subject  string
	TimeAt   string
	Teacher  string
	Location string
	Class    string
	Link     string
}

//go:embed message.tmpl
var msgtmpl string

func Message(in *MessageParams) string {
	bb := bytes.NewBuffer([]byte{})

	template.Must(
		template.New("message").
			Parse(msgtmpl),
	).Execute(bb, in)

	return string(bb.Bytes())
}
