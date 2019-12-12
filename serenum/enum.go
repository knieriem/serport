// Package serenum implements enumeration of a system's serial ports.
package serenum

import (
	"bytes"
	"log"
	"strings"
	"text/template"
)

type PortInfo struct {
	Desc   string
	Device string
	Driver string

	VendorID     string
	ProductID    string
	Manufacturer string
	SerialNumber string

	Enumerator string
}

// String returns just the device name, e.g. /dev/ttyUSB0 or COM1.
func (p *PortInfo) String() string {
	return p.Device
}

// Format creates a string from the information found in a PortInfo struct,
// making use of a template. If t is nil, the StdFormat will be used.
func (p *PortInfo) Format(t *template.Template) string {
	var b bytes.Buffer

	if t == nil {
		t = tpl
	}
	err := t.Execute(&b, p)
	if err != nil {
		log.Println(err)
		return "<error>"
	}
	return b.String()
}

const StdFormat = `
{{- if $d := .Desc}}
	{{- with .Enumerator}}
		{{- if not (contains $d .)}}{{.}}: {{end}}
	{{- end}}{{$d}}
	{{- if .Driver}}, {{end}}
{{- end}}

{{- with .Driver}}driver: {{.}}
{{- end}}

{{- if .VendorID}}, v/p: {{.VendorID}}{{with .ProductID}}:{{.}}{{end}}
{{- end}}

{{- with .SerialNumber}}, s/n: {{.}}
{{- end}}`

var tpl = template.Must(template.New("format").Funcs(template.FuncMap{
	"contains": func(s, sub string) bool {
		return strings.Contains(s, sub)
	},
}).Parse(StdFormat))

func (p *PortInfo) isPL2303() bool {
	if v := p.VendorID; v != "067b" && v != "067B" {
		return false
	}
	return p.ProductID == "2303"
}

func matchPL2303(p0, p1 *PortInfo) (isLess bool, match bool) {
	is1 := p1.isPL2303()
	if p0.isPL2303() {
		match = true
		if is1 {
			isLess = less(p0, p1)
		} else {
			isLess = true
		}
	} else if is1 {
		match = true
	}
	return
}

func matchDesc(p0, p1 *PortInfo, s string, toFront bool) (isLess bool, match bool) {
	c0 := strings.Contains(p0.Desc, s)
	c1 := strings.Contains(p1.Desc, s)
	if c0 {
		match = true
		if c1 {
			isLess = less(p0, p1)
		} else {
			isLess = toFront
		}
	} else if c1 {
		match = true
		isLess = !toFront
	}
	return
}

func less(p0, p1 *PortInfo) bool {
	return p0.Device < p1.Device
}

type portList []*PortInfo

func (list portList) Len() int {
	return len(list)
}

func (list portList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
