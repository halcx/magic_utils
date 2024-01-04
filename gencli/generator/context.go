package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

type objectGenCtx struct {
	Collector  *markers.Collector
	Checker    *loader.TypeChecker
	HeaderText string
	DebugMode  bool
}

func newObjectGenCtx(collector *markers.Collector, checker *loader.TypeChecker, headerText string, debugMode bool) *objectGenCtx {
	return &objectGenCtx{Collector: collector, Checker: checker, HeaderText: headerText, DebugMode: debugMode}
}

// generateForPackage generates IOCGolang init and runtime.Object implementations for
// types in the given package, writing the formatted result to given writer.
// May return nil if source could not be generated.
func (ctx *objectGenCtx) generateForPackage(genCtx *genall.GenerationContext, root *loader.Package) []byte {
	ctx.Checker.Check(root)

	root.NeedTypesInfo()

	imports := newImportsList(root)

	// avoid confusing aliases by "reserving" the root package's name as an alias
	imports.byAlias[root.Name] = ""

	infos := make([]*markers.TypeInfo, 0)
	if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
		infos = append(infos, info)
	}); err != nil {
		root.AddError(err)
		return nil
	}
	outContent := new(bytes.Buffer)

	copyCtx := newCopyMethodMaker(root, imports, outContent, ctx.DebugMode)

	needGen := false
	for _, info := range infos {
		if len(info.Markers["ioc:bean"]) != 0 {
			needGen = true
			if ctx.DebugMode {
				fmt.Printf("==========\n[Gen Pkg %s] Found struct that needs to gen code\n", root.PkgPath)
			}
			break
		}
	}
	if !needGen {
		if ctx.DebugMode {
			fmt.Printf("==========\n[Skip Pkg %s] Not found struct under the package that needs to gen code\n", root.PkgPath)
		}
		return nil
	}

	copyCtx.generateMethodsFor(genCtx, root, imports, infos)

	outBytes := outContent.Bytes()

	outContent = new(bytes.Buffer)
	writeHeader(root, outContent, root.Name, imports, ctx.HeaderText)
	writeMethods(root, outContent, outBytes)

	outBytes = outContent.Bytes()
	formattedBytes, err := format.Source(outBytes)
	if err != nil {
		root.AddError(err)
		// we still write the invalid source to disk to figure out what went wrong
	} else {
		outBytes = formattedBytes
	}

	return outBytes
}
