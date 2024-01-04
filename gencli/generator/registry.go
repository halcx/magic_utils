package generator

import (
	"fmt"
	"io"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

// copyMethodMakers makes DeepCopy (and related) methods for Go types,
// writing them to its codeWriter.
type copyMethodMaker struct {
	pkg *loader.Package
	*importsList
	*codeWriter
	debugMode bool
}

func newCopyMethodMaker(pkg *loader.Package, importsList *importsList, out io.Writer, debugMode bool) *copyMethodMaker {
	return &copyMethodMaker{
		pkg:         pkg,
		importsList: importsList,
		codeWriter: &codeWriter{
			Out: out,
		},
		debugMode: debugMode,
	}
}

//// +ioc:bean=true
//// +ioc:bean:name=testStruct
//// +ioc:bean:type=*testStruct
//type testStruct struct {
//}

/**
func init() {
	utils.Must(
		xapp.RegisterBeanDefinition(
			"testStruct",
			beans.MustNewBeanDefinition(
				reflect.TypeOf((*testStruct)(nil)),
			),
		),
	)
}

*/

// generateMethodsFor 根据扫描到的info生成代码
func (c *copyMethodMaker) generateMethodsFor(ctx *genall.GenerationContext, root *loader.Package, imports *importsList, infos []*markers.TypeInfo) {
	c.Line(`func init() {`)
	for _, info := range infos {
		if c.debugMode {
			fmt.Printf("[Scan Struct] %s.%s\n", root.PkgPath, info.Name)
			for key, v := range info.Markers {
				fmt.Printf("[Scan Struct %s Marker] with marker: key = %s, value = %+v\n", info.Name, key, v)
			}
		}

		if len(info.Markers["ioc:bean"]) == 0 {
			continue
		}
		if !info.Markers["ioc:bean"][0].(bool) {
			continue
		}
	}
	c.Line(`}`)

}
