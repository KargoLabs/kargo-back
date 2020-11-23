package urlValidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidURL(t *testing.T) {
	c := require.New(t)

	c.True(IsValidURL("https://www.google.com/maps/place/Intec/@18.4880037,-69.9646837,17z/data=!3m1!4b1!4m5!3m4!1s0x8eaf8a18c3d58be3:0x1455e0381fcfb97d!8m2!3d18.4880037!4d-69.962495"))
	c.False(IsValidURL("ocampo.org"))
}

func TestIsValidURLOfGivenDomain(t *testing.T) {
	c := require.New(t)

	c.True(IsValidURLOfGivenDomain("https://maps.app.goo.gl/BNbVqV8PfoYAn1q58", "https://maps.app.goo.gl/"))
	c.False(IsValidURLOfGivenDomain("https://maps.app.goo.gl/BNbVqV8PfoYAn1q58", "https://www.google.com/maps/place/"))
}
