// Package emailtemplate - пакет для подготовки писем по html-шаблонам.
package emailtemplate

import (
	"bytes"
	"html/template"
)

// Content - значения, которые подставляются в html-шаблоны при формировании писем.
type Content struct {
	Name             string
	Password         string
	VerificationCode []string
	LinkLogin        string
	LinkHelp         string
	InvitationLink   string
}

// HandleMailTemplate - подготовка письма по html-шаблону.
func HandleMailTemplate(mcase MailTemplate, content Content) (string, error) {
	buf := &bytes.Buffer{}

	files := templates[mcase]

	ts, err := template.ParseFiles(files...)
	if err != nil {
		return "", err
	}

	err = ts.ExecuteTemplate(buf, "base", content)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
