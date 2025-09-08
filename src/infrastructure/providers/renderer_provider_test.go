package providers

import (
	"testing"

	"gormgoskeleton/src/application/shared/settings"

	"github.com/stretchr/testify/assert"
)

func TestRenderBaseProvider(t *testing.T) {
	assert := assert.New(t)
	type TestData struct {
		Name string
		Age  int
	}

	renderer := RendererBase[TestData]{}
	data := TestData{
		Name: "John Doe",
		Age:  30,
	}
	settings.AppSettingsInstance.TemplatesPath = "../../application/shared/templates/"
	templatePath := settings.AppSettingsInstance.TemplatesPath + "test/render_template.gohtml"
	expectedOutput := "Hello, John Doe! You are 30 years old.\n"

	result, err := renderer.Render(templatePath, data)

	assert.Nil(err)
	assert.Equal(expectedOutput, result)
}
