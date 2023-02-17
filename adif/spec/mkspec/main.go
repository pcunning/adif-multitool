// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"embed"
	"encoding/xml"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/stoewer/go-strcase"
)

type recordValue struct {
	Name string `xml:"name,attr"`
	Val  string `xml:",chardata"`
}

type header struct {
	Values []string `xml:"value"`
}
type record struct {
	Values []recordValue `xml:"value"`
}

func (r record) Value(name string) string {
	for _, v := range r.Values {
		if v.Name == name {
			return v.Val
		}
	}
	return ""
}

func (r record) CamelCaseValue(name string) string {
	return strcase.UpperCamelCase(r.Value(name))
}

func (r record) ValueMap() map[string]string {
	res := make(map[string]string)
	for _, v := range r.Values {
		res[strcase.UpperCamelCase(v.Name)] = v.Val
	}
	return res
}

type specification struct {
	Header  header   `xml:"header"`
	Records []record `xml:"record"`
}

type field record

func (f field) Value(name string) string { return record(f).Value(name) }

func (f field) Identifier() string {
	return upperCamelIdent(f.Value("Field Name")) + "Field"
}

func (f field) EnumName() string {
	e := f.Value("Enumeration")
	if e == "" {
		return ""
	}
	if e == "Sponsored_Award" {
		// The Award_Sponsor enum has prefixes like ADIF_ and the SponsoredAwardList
		// fields have values like ADIF_CENTURY_BASIC
		return "Award_Sponsor"
	}
	if e == "Credit Award (import-only)" {
		// CREDIT_GRANTED and CREDIT_SUBMITTED used to take an AwardList and now prefer a CreditList
		return "Credit"
	}
	if i := strings.IndexRune(e, '['); i > 0 {
		return e[:i]
	}
	return e
}

func (f field) EnumScope() string {
	e := f.Value("Enumeration")
	// e.g. "Submode[MODE]" or "Primary_Administrative_Subdivision[MY_DXCC]" for MY_STATE
	i := strings.IndexRune(e, '[')
	if i < 0 {
		return ""
	}
	return strings.TrimSuffix(e[i+1:], "]")
}

func (f field) FieldName() string    { return f.Value("Field Name") }
func (f field) DataType() string     { return f.Value("Data Type") }
func (f field) Description() string  { return f.Value("Description") }
func (f field) Comments() string     { return f.Value("Comments") }
func (f field) MinimumValue() string { return f.Value("Minimum Value") }
func (f field) MaximumValue() string { return f.Value("Maximum Value") }
func (f field) ImportOnly() bool     { return f.Value("Import-only") == "true" }
func (f field) HeaderField() bool    { return f.Value("Header Field") == "true" }

type fieldSpec struct {
	Header header  `xml:"header"`
	Fields []field `xml:"record"`
}

type enumSpec struct {
	Name string `xml:"name,attr"`
	specification
}

func (e *enumSpec) ValueField() string {
	return e.Header.Values[1]
}

func (e *enumSpec) TypeIdentifier() string {
	return strcase.UpperCamelCase(e.Name) + "Enum"
}

func (e *enumSpec) ValueIdentifier(r record) string {
	// it seems the value used in data files is always the second field, first is the enum name, e.g.
	// "Enumeration Name", "Mode", "Submodes", "Description", "Import-only", "Comments"
	name := fixIdentifierPat.ReplaceAllString(r.Values[1].Val, "_")
	if r.Values[0].Val == "Country" {
		// COMOROS and PALESTINE were deleted and added again with new entity codes
		if r.Value("Deleted") == "true" {
			name += "_Deleted"
		}
	} else if r.Values[0].Val == "Primary_Administrative_Subdivision" {
		// Many countries use the same short abbreviations for states/regions,
		// so disambiguate with the DXCC code
		for _, v := range r.Values {
			if v.Name == "DXCC Entity Code" {
				// special case for reused Russian/Austrian subregions
				if (v.Val == "15" || v.Val == "206") && strings.Contains(r.Value("Comments"), "for contacts made before") {
					name = fmt.Sprintf("%s_%s_old", name, v.Val)
				} else {
					name = fmt.Sprintf("%s_%s", name, v.Val)
				}
				break
			}
		}
	}
	return strcase.UpperCamelCase(e.Name) + name
}

type enumerationList struct {
	Enums []enumSpec `xml:"enumeration"`
}

// <adif version="3.1.4" status="Released" created="2022-12-06T22:03:52Z">
type adifSpec struct {
	DataTypes    specification   `xml:"dataTypes"`
	Fields       fieldSpec       `xml:"fields"`
	Enumerations enumerationList `xml:"enumerations"`
	Version      string          `xml:"version,attr"`
	Status       string          `xml:"status,attr"`
	Created      string          `xml:"created,attr"`
	Source       string
	SpecUrl      string
}

func main() {
	var content []byte
	var filename string
	if len(os.Args) > 1 {
		filename = os.Args[1]
		if strings.HasPrefix(filename, "https:") {
			name, err := fetch(filename)
			if err != nil {
				log.Fatalf("Could not downoad %s: %v", filename, err)
			}
			c, err := xmlFromZip(name)
			if err != nil {
				log.Fatalf("Could not read all.xml from %s: %v", name, err)
			}
			content = c
		} else {
			c, err := os.ReadFile(filename)
			if err != nil {
				log.Fatalf("Could not read %s: %v", filename, err)
			}
			content = c
		}
	} else {
		filename = os.Stdin.Name()
		c, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Could not read %s: %v", filename, err)
		}
		content = c
	}

	spec := adifSpec{Source: filename}
	if err := xml.Unmarshal(content, &spec); err != nil {
		log.Fatalf("XML decoding error in %s: %v", filename, err)
	}
	shortVer := strings.ReplaceAll(spec.Version, ".", "")
	spec.SpecUrl = fmt.Sprintf("https://adif.org/%s/ADIF_%s.htm", shortVer, shortVer)
	log.Printf("Parsed ADIF version %s from %s", spec.Version, filename)
	log.Printf("Version %s has %d data types, %d fields, %d enums", spec.Version,
		len(spec.DataTypes.Records), len(spec.Fields.Fields), len(spec.Enumerations.Enums))
	addCountryEnum(&spec.Enumerations)
	files := []string{"version", "data_types", "fields", "enumerations"}
	for _, f := range files {
		f = f + ".go"
		if err := generateFile(f, fmt.Sprintf("templates/%s.tmpl", f), &spec); err != nil {
			log.Fatalf("could not generate %s: %v", f, err)
		}
	}
}

func generateFile(filename, tmplPath string, spec *adifSpec) error {
	name := path.Base(tmplPath)
	log.Printf("generating %s from %s", filename, tmplPath)
	tmpl := template.New(name).Funcs(templateFuncs)
	tmpl, err := tmpl.ParseFS(templates, tmplPath)
	if err != nil {
		return err
	}
	var out bytes.Buffer
	if err = tmpl.Execute(&out, spec); err != nil {
		return fmt.Errorf("generating %s: %v", filename, err)
	}
	src, err := format.Source(out.Bytes())
	if err != nil {
		return fmt.Errorf("formatting %s: %v", filename, err)
	}
	return os.WriteFile(filename, src, fileMode)
}

func addCountryEnum(el *enumerationList) {
	for _, e := range el.Enums {
		if e.Name == "DXCC_Entity_Code" {
			country := enumSpec{Name: "Country"}
			country.Header = header{Values: []string{"Enumeration Name", "Entity Name", "Entity Code", "Deleted"}}
			country.Records = make([]record, len(e.Records))
			for i, r := range e.Records {
				country.Records[i] = record{Values: []recordValue{
					{Name: "Enumeration Name", Val: "Country"},
					{Name: "Entity Name", Val: r.Value("Entity Name")},
					{Name: "Entity Code", Val: r.Value("Entity Code")},
					{Name: "Deleted", Val: r.Value("Deleted")},
				}}
			}
			el.Enums = append(el.Enums, country)
			return
		}
	}
	log.Print("Could not make Country enum, DXCC_Entity_Code enum missing")
}

func upperCamelIdent(name string) string {
	return strcase.UpperCamelCase(fixIdentifierPat.ReplaceAllString(name, "_"))
}

var (
	fixIdentifierPat = regexp.MustCompile(`\W+`)
	//go:embed templates
	templates     embed.FS
	fileFlags                 = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	fileMode      os.FileMode = 0644
	templateFuncs             = template.FuncMap{
		"upperCamel": upperCamelIdent,
	}
)
