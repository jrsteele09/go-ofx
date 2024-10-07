package ofx_test

import (
	"os"
	"testing"

	"github.com/jrsteele09/go-ofx/ofx"
	"github.com/stretchr/testify/assert"
)

func TestOFX_XML_Parsing(t *testing.T) {
	file, err := os.Open("./testdata/statement.xml")
	assert.NoError(t, err)
	transactions, err := ofx.Parse(file)
	assert.NoError(t, err)
	assert.Len(t, transactions, 2)
	defer file.Close()
}
