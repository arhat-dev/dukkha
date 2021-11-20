# archive

Manage archives in declarative yaml config without external dependencies.

## Config

```yaml
tools:
  archive:
  - name: <name your archive tool>
  - name: <another archive tool>
```

## Supported Tasks

### Task `archive:create`

Create archives

```yaml
archive:create:
- name: foo
  matrix:
    # this task recognizes `kernel` matrix values for archive format defaulting
    kernel:
    - linux
    - windows

  # archive format, `zip` or `tar`
  #
  # defaults to zip when matrix.kernel is windows, otherwise tar
  #format: tar

  # output path of the created archive
  output: some-archive.tag.gz

  compression:
    # enable compression
    #
    # defaults to false
    enabled: true

    # method of compression
    #
    # for format tar, one of [gzip, deflate, bzip2, zstd, lzma, xz, zstd]
    # for format zip, one of [deflate, bzip2, zstd, lzma, xz, zstd]
    #
    # defaults to defalte when format is zip
    # defaults to gzip when format is tar
    #method: gzip

    # level of compression
    #
    # defaults to "5"
    level: "5"

  # files to archive and their paths in final archive
  files:
  - # local path of the file to be added, path glob ** and * is supported
    from: path/to/**/file
    # path in archive, add `/` suffix if it should be a directory
    #
    # when your `from` matches multiple files, this must be with `/` suffix
    to: somewhere/
```
