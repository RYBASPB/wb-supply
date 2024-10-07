package transformData

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TransformSuite struct {
	suite.Suite
}

type CopySuite struct {
	suite.Suite
}

func (s *TransformSuite) TestMoscow(t provider.T) {
	t.NewStep("Rename warehouse test for Тула")
	transformedWarehouse := renameWarehouse("Тула")
	assert.Equal(t, "Москва", transformedWarehouse, "Неправильная трансформация имени склада")
}

func (s *TransformSuite) TestSPb(t provider.T) {
	t.NewStep("Rename warehouse test for СПб")
	transformedWarehouse := renameWarehouse("Санкт-Петербург Уткина Заводь")
	assert.Equal(t, "СПб", transformedWarehouse, "Неправильная трансформация имени склада")
}

func (s *CopySuite) TestCopyArticles(t provider.T) {

}

func TestTransform(t *testing.T) {
	suite.RunSuite(t, new(TransformSuite))
	suite.RunSuite(t, new(CopySuite))
}
