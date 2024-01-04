package marker

import "sigs.k8s.io/controller-tools/pkg/markers"

type DefinitionGetter interface {
	GetMarkerDefinition() *markers.Definition
}
