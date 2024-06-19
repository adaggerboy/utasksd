package renderer

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"os"
	"sync"
)

type IPageRenderer interface {
	Render(data any) (renderedData []byte, err error)
	LoadSubTemplate(renderer IPageRenderer) (err error)
}

type TemplatePageRenderer struct {
	IPageRenderer
	templateEngine template.Template
}

func (r *TemplatePageRenderer) Render(data any) (result []byte, err error) {
	result = []byte{}
	buffer := bytes.NewBuffer(result)
	err = r.templateEngine.Execute(buffer, data)
	result = buffer.Bytes()
	return
}

func NewTemplatePageRenderer(name string, templateText string) (renderer IPageRenderer, err error) {

	funcMap := template.FuncMap{
		"jsArray": func(slice []string) template.JS {
			str := "["
			first := true
			for _, v := range slice {
				if first {
					str += fmt.Sprintf("'%s'", v)
				} else {
					str += fmt.Sprintf(", '%s'", v)
				}
			}
			str += "]"
			return template.JS(str)
		},
	}

	templ, err := template.New(name).Funcs(funcMap).Parse(templateText)
	if err != nil {
		return
	}
	renderer = &TemplatePageRenderer{
		templateEngine: *templ,
	}
	return
}

func NewTemplatePageRendererFromFile(name string, file string) (renderer IPageRenderer, err error) {
	fl, err := os.Open(file)
	if err != nil {
		return
	}
	data, err := io.ReadAll(fl)
	if err != nil {
		return
	}
	return NewTemplatePageRenderer(name, string(data))
}

var (
	poolMtx sync.RWMutex
	poolMap map[string]IPageRenderer = map[string]IPageRenderer{}
)

func LoadToPool(files map[string]string) error {
	poolMtx.Lock()
	defer poolMtx.Unlock()
	for k, v := range files {
		rnd, err := NewTemplatePageRendererFromFile(k, v)
		if err != nil {
			return err
		}
		poolMap[k] = rnd
	}
	return nil
}

func GetRenderer(key string) (renderer IPageRenderer, err error) {
	poolMtx.RLock()
	defer poolMtx.RUnlock()
	renderer, ok := poolMap[key]
	if !ok {
		err = fmt.Errorf("renderer key not found")
	}
	return
}
