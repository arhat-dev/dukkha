# Include Remote Files

Question: There is a top level option `include` in dukkha config, but only supports including local files, how can I include remote files?

## TL;DR

Create a local file, within which write a top level virtual key (`__`) with some renderer able to fetch remote content, then include that local file in your dukkha config.

## Detailed Answer

Given you have a yaml file stored at some http server, its url is `https://example.com/my-config.yaml`, and you want to include it in local dukkha config

1. First create a yaml file locally, say `./remote/my-config.yaml`

2. Write a virtual key and tell it how to render the whole yaml file, use `str` type hint to prevent dukkha from rendering it before task execution.

   ```yaml
   # filename: ./remote/my-config.yaml
   __@http?str: https://example.com/my-config.yaml
   ```

3. Include the newly created local yaml file in your dukkha config

   ```yaml
   include:
   - ./remote/my-config.yaml
   ```

The reason why we prefer this way is you can easily find out what's actually included using `dukkha render`, as long as it uses the same rendering context (same set of dukkha config, regardless of the file to included included or not), it will render the same result `dukkha run` consumes.

There is another way to do this if you don't want to create extra local files: use virtual key and renderer attribute `cached-file` directly in `include`, by this way, you lose the ability to inspect what's getting included using `dukkha render` directly (but you can still find these files by path), so it's not recommended, but you can do it anyway.

```yaml
include:
- __@http#cached-file: https://example.com/my-config.yaml
- __@git#cached-file: my-repo.git/my-file
```

## Tips: Always Verify Checksum of remote config

```yaml
__@T?str:
   value@http?str: https://gist.githubusercontent.com/arhatbot/d1f27e2b6d7e41a7c9d0a6ef7e39a921/raw/1e014333a3d78ac1139bc4cab9a68685e5080685/echo.sh
   ops:
   - checksum:
      data@tpl?str: |-
         {{- VALUE -}}
      kind: sha256
      # yes, you can host dynamic checksum if your file is dynamic
      # but make sure their renderer having same caching policy
      # so it won't get false positive rejection
      sum@http: https://gist.githubusercontent.com/arhatbot/d1f27e2b6d7e41a7c9d0a6ef7e39a921/raw/f36f55f83d2bb0d3d45e79e48360bb5c35826048/echo.sh.sha256
```

## Ext. Similar Solutions In Other Software

- Jenkins: use plugin [remote file](https://plugins.jenkins.io/remote-file/) and configure using web ui, pom.xml, plus some DSL
- GitLab CI: the `include` option in yaml config supports `local`, combination of `project`:`ref`:`file`, `remote` (http), `template`
- GitHub CI: Actions, `fromJsom`, or call other workflows with github script
