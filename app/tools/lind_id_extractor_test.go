package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLinkIDExtracted(t *testing.T) {
	link := "id:my-article"
	id, ok := ExportLinkID(link)
	assert.Equal(t, id, "my-article")
	assert.Equal(t, ok, true)
}
