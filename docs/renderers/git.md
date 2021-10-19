# Git Renderer

```yaml
foo@git: my-org/foo.git/foo.yaml@master
```

Fetch file content from your git ssh repo as the field value.

## Config Options

__NOTE:__ Configuration is required to activate this renderer, ssh config is required to make it work with string input like `my-org/foo.git/foo.yaml@master`.

```yaml
renderers:
  # no options
  git:
    # cache config
    # enable local cache, disable to always fetch from remote
    enable_cache: true
    cache_max_age: 1h

    # git ssh config
    # git ssh user, defaults to git
    user: foo
    # git ssh service host
    host: example.com
    # git ssh service port, defaults to 22
    port: 60022
    # public host key for remote host verification
    # will skip host verification if not set
    host_key: ""
    # git ssh private key
    private_key: ""
    # git ssh password, not effective if private_key is set
    password: ""
```

## Supported value types

- String: SCP style URL with optional `host:port` part and `@<ref>` suffix (when you have configured ssh options of your git renderer in renderer config)

```yaml
# example 1: without host:port part (using that configured in renderer config) and with @<ref>
foo@git: my-org/foo.git/foo.yaml@master

# example 2: with host:port but without @<ref> (using HEAD)
bar@git: my-domain.com:1022:my-org/foo.git/foo.yaml

# example 3: with host but not port
# the port will defaults to 22 rather than using the port configured in renderer config
woo@git: my-domain.com:my-org/foo.git/foo.yaml@dev
```

- Valid git fetch spec in yaml (you can omit ssh options if you have configured them in renderer config and you don't want to override it)

```yaml
# git ssh settings
ssh:
  # git ssh config
  # git ssh user, defaults to git
  user: foo
  # git ssh service host
  host: example.com
  # git ssh service port, defaults to 22
  port: 60022
  # public host key for remote host verification
  # will skip host verification if not set
  host_key: ""
  # git ssh private key
  private_key: ""
  # git ssh password, not effective if private_key is set
  password: ""

# repo name with .git suffix
repo: my-org/foo.git
# target fetch path in repo
path: foo.yaml
# git object ref, usually a branch/tag name
# defaults to HEAD (default branch in remote)
ref: master
```

## Suggested Use Cases

Access file content in a ssh only repo.
