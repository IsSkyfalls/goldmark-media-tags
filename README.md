# goldmark-media-tags

A media tags extension for [GoldMark](https://github.com/yuin/goldmark/).

Inspired by
flexmark-java [<flexmark-ext-media-tags>](https://github.com/vsch/flexmark-java/wiki/Extensions#media-tags) (java)

## Supports

`<media>` and `<audio>`.

```markdown
!v[This is a video](https://example.org/test.webm)
!a[And this is an audio](https://example.org/test.mp3)
```

Renders to:

```html
<p>
    <video controls>
        <source src="https://example.org/test.webm">
    </video>
    <audio controls>
        <source src="https://example.org/test.mp3">
    </audio>
</p>
```

## Usage

check [extension_test.go](./extension_test.go)




