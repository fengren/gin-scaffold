package template

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/volatiletech/inflect"
)

type state struct {
	n int
}

func (s *state) Set(n int) int {
	s.n = n
	return n
}

func (s *state) Inc() int {
	s.n++
	return s.n
}

var s state
var (
	funcMap = template.FuncMap{
		"Pluralize":  inflect.Pluralize,
		"Underscore": inflect.Underscore,
		"ToUpper":    strings.ToUpper,
		"ToLower":    strings.ToLower,
		"set":        s.Set,
		"inc":        s.Inc,
		"is_tmp": func(fieldType string) bool {
			switch fieldType {
			case
				"int16",
				"int32",
				"int64":
				return true
			}
			return false
		},
		"ret": func(fieldType string) string {
			switch fieldType {
			case
				"int",
				"float64",
				"int16",
				"int32",
				"int64",
				"bool":
				return ", _"
			}
			return ""
		},
		"conv": func(origin string, fieldType string, misc string) string {
			if fieldType == "int" {
				return "strconv.Atoi(" + origin + ")"
			} else if fieldType == "int16" {
				return "strconv.ParseInt(" + origin + ", 10, 16)"
			} else if fieldType == "int32" {
				return "strconv.ParseInt(" + origin + ", 10, 32)"
			} else if fieldType == "int64" {
				// original time is time
				if misc == "time" {
					return "time.Parse(\"2006-01-02 15:04:05\", " + origin + ")"
				}
				return "strconv.ParseInt(" + origin + ", 10, 64)"
			} else if fieldType == "float64" {
				return "strconv.ParseFloat(" + origin + ", 64)"
			} else if fieldType == "bool" {
				return "strconv.ParseBool(" + origin + ")"
			}
			return origin
		},
	}
)

type Builder struct {
	TemplateName string
	TemplatePath string
}

func NewBuilder(templatePath string) *Builder {
	if !filepath.IsAbs(templatePath) {
		templatePath = TemplatePath(templatePath)
	}

	templateName := filepath.Base(templatePath)
	builder := &Builder{
		TemplateName: templateName,
		TemplatePath: templatePath,
	}

	return builder
}

func (builder *Builder) Template() *template.Template {
	contents := LoadTemplateFromFile(builder.TemplatePath)
	tmpl := template.Must(template.New(builder.TemplateName).Funcs(funcMap).Parse(contents))

	return tmpl
}

func (builder *Builder) Write(writer io.Writer, data interface{}) {
	tmpl := builder.Template()
	err := tmpl.Execute(writer, data)
	if err != nil {
		panic(err)
	}
}

func (builder *Builder) WriteToPath(outputPath string, data interface{}) {
	printAction("green+h:black", "create", outputPath)
	if _, err := os.Stat(outputPath); err == nil {
		printAction("red+h:black", "skip", outputPath)
		return
	}

	file, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	builder.Write(file, data)
}

func (builder *Builder) InsertAfterToPath(outputPath string, after string, data interface{}) {
	printAction("cyan+h:black", "insert", outputPath)

	newFilePath := outputPath + ".new"

	file, err := os.Open(outputPath)
	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create(newFilePath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outputFile)

	for scanner.Scan() {
		line := scanner.Text()

		writer.WriteString(line + "\n")
		if strings.HasPrefix(line, after) {
			builder.Write(writer, data)
		}
	}

	writer.Flush()
	outputFile.Close()
	file.Close()

	os.Rename(newFilePath, outputPath)
}
