package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkIDExtracted(t *testing.T) {
	link := "id:my-article"
	id, ok := ExportLinkID(link)
	assert.Equal(t, id, "my-article")
	assert.Equal(t, ok, true)
}
