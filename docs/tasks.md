# Tasks

All tasks using default tool will be available to all other tools with same tool kind

## Common Task Config Options

All tasks have a `arhat.dev/dukkha/pkg/tools.BaseTask` embedded

- `tool_kind`: The tool implementation name (e.g. `docker`, `buildah`)
  - aka: Sub-directory names in `${project_root}/pkg/tools`
- `tool_name`: Custom name configured in `tools.<tool_kind>[*].name`
  - This is optional if the task is available for all tools with `<tool_kind>`
- `task_kind`: Task kind handled by the tool

```yaml
<tool_kind>{:<tool_name>}:<task_kind>:
- name: <task_name>
  # task hooks
  hook:
    before:
    # use the embedded shell to run commands
    - shell: |-
        echo "foo"

    # use a specific shell to run commands
    - shell:<shell_name>: |-
        print("bar")

    # run another predefined task
    - task: <tool_kind>{:<tool_name>}:<task_kind>:<another_task_name>

    # run commands/tasks after a successful task execution
    after:success: []

    # run commands/tasks when the task execution failed
    after:failure: []

    # run before/after matrix execution instead of the whole task
    before:matrix: []
    after:matrix:success: []
    after:matrix:failure: []
```
