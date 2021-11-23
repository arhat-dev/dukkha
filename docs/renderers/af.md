# Archive File Renderer

```yaml
foo@af: path/to/your/archive:in/archive/path
```

Extract file content from archive.

## Config Options

__NOTE:__ Configuration is required to activate this renderer.

```yaml
renderers:
  af:
    # enable local file caching for extracted files
    # disable to always extract from archive at runtime
    enable_cache: true
    cache_item_size_limit: 64KB
    cache_max_age: 10h
    cache_size_limit: 32MB
```

## Supported value types

- String: SCP style path `<path-to-archive>:<in-archive-path>` or `<>` to use first regular file in archive

  ```yaml
  # example 1: without in archive path
  foo@af: path/to/archive.tar.gz

  # example 2: with in archive path
  bar@af: path/to/archive.zip.gz:foo.yaml
  ```

- Valid archive file extraction spec in yaml

  ```yaml
  foo@af:
    # path to the target archive
    archive: path/to/the/archive
    path: in/archive/path

    # NOTE: currently not implemented
    password: password for encrypted rar/zip
  ```

## Supported Attributes

- `cached-file`: Return local file path to cached file instead of fetched content.
- `allow-expired`: Allow to use previously extracted file when it's not possible to extract from original archive (e.g. archive deleted)

## Supported Archive Formats

- Plain archive files:
  - `tar`
  - `zip`: with internal compression using
    - `deflate`
    - `zstd`
    - `bzip2`
    - `xz`
    - `lzma`
- Compressed files (including compressed archive files) using following compression methods:
  - `gzip`
  - `bzip2`
  - `xz`
  - `lzma`
  - `lz4`
  - `deflate`

## Suggested Use Cases

Work with archives but do not directly touch them in shell.
