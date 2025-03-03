package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractIIN(t *testing.T) {
	data := `SERIALNUMBER=IIN741012301763,CN=БОРЗИЛОВ АЛЕКСЕЙ,C=KZ,2.5.4.42=#0c10d09fd095d0a2d0a0d09ed092d098d0a7,2.5.4.4=#0c10d091d09ed0a0d097d098d09bd09ed092`
	iin := ExtractIIN(data)
	assert.Equal(t, "741012301763", iin)

}
