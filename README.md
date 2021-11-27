# goldmark-media-tags

[![Go Reference](https://pkg.go.dev/badge/github.com/IsSkyfalls/goldmark-media-tags.svg)](https://pkg.go.dev/github.com/IsSkyfalls/goldmark-media-tags)

A media tags extension for [GoldMark](https://github.com/yuin/goldmark/).

Inspired by
flexmark-java [\<flexmark-ext-media-tags\>](https://github.com/vsch/flexmark-java/wiki/Extensions#media-tags) (java)

## Supports

`<media>`, `<audio>` and `<picture>.

```markdown
!v[This is a video](https://example.org/test.webm)
!a[And this is an audio](https://example.org/test.mp3)
!p[And this a picture](https://example.org/test.png)
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
    <picture>
        <img alt="And this a picture" src="https://example.org/test.png">
    </picture>
</p>
```

## Usage

check [extension_test.go](./extension_test.go)




