package main

import (
	"fmt"
	"github.com/halcx/magic_utils/gencli/generator"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

func main() {
	var (
		// allGenerators maintains the list of all known generators, giving
		// them names for use on the command line.
		// each turns into a command line option,
		// and has options for output forms.
		allGenerators = map[string]genall.Generator{
			"register": generator.Generator{},
		}

		allOutputRules = map[string]genall.OutputRule{
			"dir":       genall.OutputToDirectory(""),
			"none":      genall.OutputToNothing,
			"stdout":    genall.OutputToStdout,
			"artifacts": genall.OutputArtifacts{},
		}

		// optionsRegistry contains all the marker definitions used to process command line options
		optionsRegistry = &markers.Registry{}
	)

	for genName, gen := range allGenerators {
		// make the generator options marker itself
		defn := markers.Must(markers.MakeDefinition(genName, markers.DescribesPackage, gen))
		if err := optionsRegistry.Register(defn); err != nil {
			panic(err)
		}
		if helpGiver, hasHelp := gen.(genall.HasHelp); hasHelp {
			if help := helpGiver.Help(); help != nil {
				optionsRegistry.AddHelp(defn, help)
			}
		}

		// make per-generation output rule markers
		for ruleName, rule := range allOutputRules {
			ruleMarker := markers.Must(markers.MakeDefinition(fmt.Sprintf("output:%s:%s", genName, ruleName), markers.DescribesPackage, rule))
			if err := optionsRegistry.Register(ruleMarker); err != nil {
				panic(err)
			}
			if helpGiver, hasHelp := rule.(genall.HasHelp); hasHelp {
				if help := helpGiver.Help(); help != nil {
					optionsRegistry.AddHelp(ruleMarker, help)
				}
			}
		}
	}

	// make "default output" output rule markers
	for ruleName, rule := range allOutputRules {
		ruleMarker := markers.Must(markers.MakeDefinition("output:"+ruleName, markers.DescribesPackage, rule))
		if err := optionsRegistry.Register(ruleMarker); err != nil {
			panic(err)
		}
		if helpGiver, hasHelp := rule.(genall.HasHelp); hasHelp {
			if help := helpGiver.Help(); help != nil {
				optionsRegistry.AddHelp(ruleMarker, help)
			}
		}
	}

	// add in the common options markers
	if err := genall.RegisterOptionsMarkers(optionsRegistry); err != nil {
		panic(err)
	}

	// add in the common options markers
	if err := genall.RegisterOptionsMarkers(optionsRegistry); err != nil {
		panic(err)
	}

	var rawOpts []string

	if len(rawOpts) == 0 {
		rawOpts = []string{"register", "paths=./gencli/test"}
	}

	if len(rawOpts) == 1 {
		rawOpts = append(rawOpts, "register")
	}

	rt, err := genall.FromOptions(optionsRegistry, rawOpts)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	if len(rt.Generators) == 0 {
		_ = fmt.Errorf("no generators specified")
	}

	if hadErrs := rt.Run(); hadErrs {
		_ = fmt.Errorf("not all generators ran successfully")
	}
}
