package view

import (
	"fmt"
	"github.com/flosch/pongo2/v6"
	"io"
	"log"
	"os"
)

type Data map[string]any

var (
	Path = os.Getenv("VIEW_PATH")
	Ext  = os.Getenv("VIEW_EXT")
)

// LoadView load view content from template.
func LoadView(template string, data Data) string {
	file, err := pongo2.FromFile(fmt.Sprintf("%s/%s.%s", Path, template, Ext))
	if err != nil {
		log.Fatal(err)
	}
	ctx := pongo2.Context(data)
	res, err := file.Execute(ctx)

	if err != nil {
		log.Fatal(err)
	}

	return res
}

// LoadViewWriter load view content into writer.
func LoadViewWriter(template string, data Data, writer io.Writer) error {
	var templateFile = pongo2.Must(pongo2.FromFile(fmt.Sprintf("%s/%s.%s", Path, template, Ext)))
	ctx := pongo2.Context(data)

	return templateFile.ExecuteWriter(ctx, writer)
}
