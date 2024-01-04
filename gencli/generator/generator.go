package generator

import (
	"fmt"
	"go/ast"

	"github.com/halcx/magic_utils/gencli/marker"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

var DebugMode = false

type Generator struct {
	HeaderFile string `marker:",optional"`
	Year       string `marker:",optional"`
}

func (Generator) CheckFilter() loader.NodeFilter {
	return func(node ast.Node) bool {
		// ignore interfaces
		_, isIface := node.(*ast.InterfaceType)
		return !isIface
	}
}

// RegisterMarkers is called in main, register all markers
func (Generator) RegisterMarkers(into *markers.Registry) error {
	defs := []*markers.Definition{
		marker.NewEnableBeanNameDefinition().GetMarkerDefinition(),
		marker.NewEnableBeanTypeDefinition().GetMarkerDefinition(),
		marker.NewEnableBeanDefinition().GetMarkerDefinition(),
	}

	return markers.RegisterAll(into, defs...)
}

func (d Generator) Generate(ctx *genall.GenerationContext) error {
	var headerText string

	if d.HeaderFile != "" {
		headerBytes, err := ctx.ReadFile(d.HeaderFile)
		if err != nil {
			return err
		}
		headerText = string(headerBytes)
	}

	objGenCtx := newObjectGenCtx(ctx.Collector, ctx.Checker, headerText, DebugMode)

	for _, root := range ctx.Roots {
		// 2. generate codes under current pkg
		outContents := objGenCtx.generateForPackage(ctx, root)
		fmt.Println(outContents)
		if outContents == nil {
			continue
		}
		// 3. write codes to file
		writeOut(ctx, nil, root, outContents)
		fmt.Println(outContents)
	}
	return nil
}
