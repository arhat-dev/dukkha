# Tools

## Common tool config options

All tools have a `arhat.dev/dukkha/pkg/tools.BaseTool` embedded

```yaml
tools:
  # tool_kind is the kind name of the tool
  <tool_kind>:
  - name: <tool_name> # custom tool name

    # set extra environment variables when running this tool
    env: []
    # - ENV_NAME=env_value

    # cmd to run this tool
    # e.g. [ssh, remote-host, do, something]
    cmd: []
```
