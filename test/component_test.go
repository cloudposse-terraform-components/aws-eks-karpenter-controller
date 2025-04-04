package test

import (
	"testing"
	helper "github.com/cloudposse/test-helpers/pkg/atmos/component-helper"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/cloudposse/test-helpers/pkg/helm"
	"github.com/stretchr/testify/assert"
)

type ComponentSuite struct {
	helper.TestSuite
}

func (s *ComponentSuite) TestBasic() {
	const component = "eks/karpenter-controller/basic"
	const stack = "default-test"
	const awsRegion = "us-east-2"

	namespace := "kube-system"

	inputs := map[string]interface{}{}

	defer s.DestroyAtmosComponent(s.T(), component, stack, &inputs)
	options, _ := s.DeployAtmosComponent(s.T(), component, stack, &inputs)
	assert.NotNil(s.T(), options)

	metadataArray := []helm.Metadata{}

	atmos.OutputStruct(s.T(), options, "metadata", &metadataArray)

	assert.Equal(s.T(), len(metadataArray), 1)
	metadata := metadataArray[0]

	assert.Equal(s.T(), metadata.AppVersion, "1.3.2")
	assert.Equal(s.T(), metadata.Chart, "karpenter")
	assert.NotNil(s.T(), metadata.FirstDeployed)
	assert.NotNil(s.T(), metadata.LastDeployed)
	assert.Equal(s.T(), metadata.Name, "karpenter")
	assert.Equal(s.T(), metadata.Namespace, namespace)
	assert.NotNil(s.T(), metadata.Values)
	assert.Equal(s.T(), metadata.Version, "1.3.2")

	s.DriftTest(component, stack, &inputs)
}

func (s *ComponentSuite) TestEnabledFlag() {
	const component = "eks/karpenter-controller/disabled"
	const stack = "default-test"
	s.VerifyEnabledFlag(component, stack, nil)
}

func (s *ComponentSuite) SetupSuite() {
	s.TestSuite.InitConfig()
	s.TestSuite.Config.ComponentDestDir = "components/terraform/eks/karpenter-controller"
	s.TestSuite.SetupSuite()
}

func TestRunSuite(t *testing.T) {
	suite := new(ComponentSuite)
	suite.AddDependency(t, "vpc", "default-test", nil)
	suite.AddDependency(t, "eks/cluster", "default-test", nil)
	helper.Run(t, suite)
}
