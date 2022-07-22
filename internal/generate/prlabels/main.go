//go:build generate
// +build generate

package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-provider-aws/names"
)

const (
	filename      = `../../../.github/labeler-pr-triage.yml`
	namesDataFile = "../../../names/names_data.csv"
)

type ServiceDatum struct {
	ProviderPackage string
	ActualPackage   string
	FilePrefix      string
	DocPrefixes     []string
}

type TemplateData struct {
	Services []ServiceDatum
}

func main() {
	fmt.Printf("Generating %s\n", strings.TrimPrefix(filename, "../../../"))

	f, err := os.Open(namesDataFile)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)

	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	td := TemplateData{}

	for i, l := range data {
		if i < 1 { // no header
			continue
		}

		if l[names.ColExclude] != "" && l[names.ColAllowedSubcategory] == "" {
			continue
		}

		if l[names.ColProviderPackageActual] == "" && l[names.ColProviderPackageCorrect] == "" {
			continue
		}

		p := l[names.ColProviderPackageCorrect]

		if l[names.ColProviderPackageActual] != "" {
			p = l[names.ColProviderPackageActual]
		}

		ap := p

		if l[names.ColSplitPackageRealPackage] != "" {
			ap = l[names.ColSplitPackageRealPackage]
		}

		s := ServiceDatum{
			ProviderPackage: p,
			ActualPackage:   ap,
			FilePrefix:      l[names.ColFilePrefix],
			DocPrefixes:     strings.Split(l[names.ColDocPrefix], ";"),
		}

		td.Services = append(td.Services, s)
	}

	sort.SliceStable(td.Services, func(i, j int) bool {
		return td.Services[i].ProviderPackage < td.Services[j].ProviderPackage
	})

	writeTemplate(tmpl, "prlabeler", td)
}

func writeTemplate(body string, templateName string, td TemplateData) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file (%s): %s", filename, err)
	}

	tplate, err := template.New(templateName).Parse(body)
	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tplate.Execute(&buffer, td)
	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	if _, err := f.Write(buffer.Bytes()); err != nil {
		f.Close()
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}

	if err := f.Close(); err != nil {
		log.Fatalf("error closing file (%s): %s", filename, err)
	}
}

var tmpl = `# YAML generated by internal/generate/prlabels/main.go; DO NOT EDIT.
client-connections:
  - 'internal/conns/**/*'
create:
  - 'internal/create/**/*'
dependencies:
  - '.github/dependabot.yml'
documentation:
  - '**/*.md'
  - 'docs/**/*'
  - 'website/**/*'
examples:
  - 'examples/**/*'
flex:
  - 'internal/flex/**/*'
generators:
  - 'internal/**/*_gen.go'
  - 'internal/**/*_gen_test.go'
  - 'internal/**/generate.go'
  - 'internal/generate/**/*'
github_actions:
  - '.github/*.yml'
  - '.github/workflows/*.yml'
linter:
  - '.github/workflows/acctest-terraform-lint.yml'
  - '.github/workflows/terraform_provider.yml'
  - '.github/workflows/website.yml'
  - '.golangci.yml'
  - '.markdownlinkcheck.json'
  - '.markdownlint.yml'
  - '.semgrep.yml'
  - '.tflint.hcl'
  - 'staticcheck.conf'
  - 'tools/providerlint/**/*'
pre-service-packages:
  - '**/data_source_aws_*'
  - '**/resource_aws_*'
  - 'aws/**/*'
  - 'awsproviderlint/**/*'
provider:
  - '.gitignore'
  - '.go-version'
  - '*.md'
  - 'docs/contributing/**/*'
  - 'internal/provider/**/*'
  - 'main.go'
  - 'website/docs/index.html.markdown'
repository:
  - '.github/**/*'
  - 'GNUmakefile'
  - 'infrastructure/**/*'
skaff:
  - 'skaff/**/*'
sweeper:
  - 'internal/sweep/**/*'
  - 'internal/service/**/sweep.go'
tags:
  - 'internal/**/tag_gen.go'
  - 'internal/**/tag_gen_test.go'
  - 'internal/**/tag_test.go'
  - 'internal/**/tags_gen.go'
  - 'internal/tags/**/*'
tests:
  - '**/*_test.go'
  - 'internal/**/test-fixtures/**/*'
  - 'internal/**/testdata/**/*'
  - 'internal/acctest/**/*'
verify:
  - 'internal/verify/**/*'
{{- range .Services }}
service/{{ .ProviderPackage }}:
  - 'internal/service/{{ .ActualPackage }}/**/{{ .FilePrefix }}*'
  {{- range .DocPrefixes }}
  - 'website/**/{{ . }}*'
  {{- end }}
{{- end }}
`
