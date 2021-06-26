# Tasks

## Common Task Config Options

All tasks have a `arhat.dev/dukkha/pkg/tools.BaseTask` embedded

- `tool_kind`: The tool implementation name (e.g. `docker`, `buildah`)
  - aka: Sub-directory names in `${project_root}/pkg/tools`
- `tool_name`: Custom name configured in `tools.<tool_kind>[*].name`
  - This is optional if you mean the default tool (the first tool configured in `tools.<tool_kind>`)
- `task_kind`: Task kind handled by this tool

```yaml
<tool_kind>{:<tool_name>}:<task_kind>:
- name: <task_name>
  # task hooks
  hook:
    before:
    # use the default shell to run commands
    - shell: |-
        echo "foo"

      # run after matrix execution instead of the whole task
      # defaults to true
      per_matrix_run: true

    # use a specific shell to run commands
    - shell:<shell_name>: |-
        print("bar")

    # run another predefined task
    - task: <tool_kind>{:<tool_name>}:<task_kind>:<another_task_name>

    # run commands/tasks after a successful task execution
    after:success:
    - task: <tool_kind>{:<tool_name>}:<task_kind>:<another_task_name>
    - shell: echo "Done."

    # run commands/tasks when the task execution failed
    after:failure:
    - task: <tool_kind>{:<tool_name>}:<task_kind>:<another_task_name>
    - shell: echo "Failed."

```
