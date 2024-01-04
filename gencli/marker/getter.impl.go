package marker

import (
	"sigs.k8s.io/controller-tools/pkg/markers"
)

// 在此处定义新的Defination
type EnableBeanDefinition struct {
}

func NewEnableBeanDefinition() *EnableBeanDefinition {
	return &EnableBeanDefinition{}
}

func (m *EnableBeanDefinition) GetMarkerDefinition() *markers.Definition {
	return markers.Must(markers.MakeDefinition("ioc:bean", markers.DescribesType, true))
}

type EnableBeanNameDefinition struct {
}

func NewEnableBeanNameDefinition() *EnableBeanNameDefinition {
	return &EnableBeanNameDefinition{}
}

func (m *EnableBeanNameDefinition) GetMarkerDefinition() *markers.Definition {
	return markers.Must(markers.MakeDefinition("ioc:bean:name", markers.DescribesType, ""))
}

type EnableBeanTypeDefinition struct {
}

func NewEnableBeanTypeDefinition() *EnableBeanTypeDefinition {
	return &EnableBeanTypeDefinition{}
}

func (m *EnableBeanTypeDefinition) GetMarkerDefinition() *markers.Definition {
	return markers.Must(markers.MakeDefinition("ioc:bean:type", markers.DescribesType, ""))
}
