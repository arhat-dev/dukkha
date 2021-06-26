# docker

## Task `docker:build`

```yaml
docker:<docker-tool-name>:build:
- name: example-image
  # images_names of the build output
  # if not set, will use the `name` value as `image`
  image_names:
  - image: example.com/image:tag-amd64
    manifest: example.com/image:manifest-tag

  # docker build <positional-arg>
  context: "."

  # docker build -f
  dockerfile: path/to/dockerfile

  # extra docker build args
  extraArgs: []
```

## Task `docker:push`
