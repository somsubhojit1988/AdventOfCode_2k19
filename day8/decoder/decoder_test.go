package decoder

import (
	"fmt"
	"testing"
)

func TestImgCreate(t *testing.T) {
	str := "123456789012"
	h, w := 2, 3
	img, err := CreateImage(str, h, w)
	if err != nil {
		t.Error(err)
	}
	t.Log(img)
}

func TestImgDecode(t *testing.T) {
	str, h, w := "0222112222120000", 2, 2
	mat, err := CreateImage(str, h, w)
	if err != nil {
		t.Error(err)
	}
	t.Logf(fmt.Sprintf("%s", mat.Decode()))
}
