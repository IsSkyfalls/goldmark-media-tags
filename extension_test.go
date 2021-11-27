package media

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"os"
	"testing"
)

func TestExtensionWithDefaults(t *testing.T) {
	md := goldmark.New(goldmark.WithExtensions(
		WithDefaults(),
	))
	src := []byte(`
!v[Train](https://upload.wikimedia.org/wikipedia/commons/5/5d/KkStB_Class_310.webm)
!a[THE Anthem](https://upload.wikimedia.org/wikipedia/commons/transcoded/5/59/State_Anthem_of_the_Soviet_Union_%281984%29.wav/State_Anthem_of_the_Soviet_Union_%281984%29.wav.mp3)
!p[Blossoms](https://upload.wikimedia.org/wikipedia/commons/b/b8/Blossoming-Blackberries-P1400842_%2837396534302%29.jpg)
`)
	err := md.Convert(src, os.Stdout)
	assert.NoError(t, err)
}

func TestExtensionWithOptions(t *testing.T) {
	md := goldmark.New(goldmark.WithExtensions(
		WithOptions(Options{
			MediaControls: true,
			MediaAutoplay: true,
			MediaLoop:     true,
			MediaMuted:    true,
			MediaPreload:  Metadata,
		}),
	))
	src := []byte(`
!v[Train](https://upload.wikimedia.org/wikipedia/commons/5/5d/KkStB_Class_310.webm)
!a[THE Anthem](https://upload.wikimedia.org/wikipedia/commons/transcoded/5/59/State_Anthem_of_the_Soviet_Union_%281984%29.wav/State_Anthem_of_the_Soviet_Union_%281984%29.wav.mp3)
!p[Blossoms](https://upload.wikimedia.org/wikipedia/commons/b/b8/Blossoming-Blackberries-P1400842_%2837396534302%29.jpg)
`)
	err := md.Convert(src, os.Stdout)
	assert.NoError(t, err)
}
