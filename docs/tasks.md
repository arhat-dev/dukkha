# Tasks

## Config Declaration

- Option 1: without tool name: `<tool-kind>:<task-kind>`
  - tasks defined under this declaration are admitted by all tools with same `<tool-kind>`
- Option 2: With tool name: `<tool-kind>:<tool-name>:<task-kind>`
  - tasks defined under this declaration are only admitted by the tool with same `<tool-kind>` and `<tool-name>`

Example:

```yaml
workflow:run:
# all workflow tool can run task `foo`
- name: foo

workflow:local:run:
# only workflow tool named `local` can run task `bar`
- name: bar
```

## Common Task Options

- `name: string`: required task name

- `env: []Env`: define task specific envrionment variables

- `matrix`
  - `kernel: []string`: special vector for cross platform tasks
  - `arch: []string`: special vector for cross platform tasks
  - `libc: []string`: special vector for cross platform tasks
  - `exclude: []map[string][]string`: exclude matched matrix entries
  - `include: []map[string][]string`: include extra vectors

- `hooks`
  - `before: []Action`: run actions before task start.
  - `before:matrix: []Action`: run actions before each task matrix run.
  - `after:matrix:success: []Action`: run actions after each successful task matrix run.
  - `after:matrix:failure: []Action`: run actions when task matrix run failed.
  - `after:matrix: []Action`: run actions when task matrix run finished, no matter failed or succeeded.
  - `after:success: []Action`: run actions after all task matrix succeeded.
  - `after:failure: []Action`: run actions after all task matrix finished but some errored.
  - `after: []Action`: run actions after all task matrix run finished, regardless of failure.

And `Action` is defined as:

- `name: string`: action name
- `task: string`: reference to other task
  - task reference format: `<tool-kind>{:<tool-name>}:<task-kind>(<another_task_name>{, <matrix-spec> })`
    - where `<matrix-sepc>` is the task matrix yaml
- `cmd: []string`: raw cmd to run
- `idle: any`: do nothing, serves as a placeholder so you can use rendering suffix for non action operations.
- `shell: string`: run script in embedded bash.
- `next: string`: name of the action as next step.
- `env: []Env`: action specific environment vairables.
- `chdir: string`: change work directory for this action
- `continue_on_error: bool`: continue next action even when this action failed

Example:

```yaml
workflow:run:
- name: example

  env:
  - name: ENV_NAME
    value: env-value

  matrix:
    kernel:
    - linux
    - openbsd
    arch:
    - amd64
    - arm64

    foo:
    - bar
    - woo

    exclude:
    # exclude by partial matching
    - foo:
      - bar
      arch:
      - amd64
    # or full matching
    - foo:
      - woo
      kernel:
      - linux
      arch:
      - amd64

    include:
    - foo:
      - gee

  # task hooks
  hooks:
    before:
    # use the embedded bash to run commands (with env expansion)
    - shell: |-
        echo "${}"

    # use a specific shell to run commands, that shell must be configured in `shells` section
    - shell:python: |-
        print("bar")

    before:matrix:
    # run another task, limit to the same matrix as current matrix
    - task: workflow:run(bar)

    after:matrix:success:
    # do nothing
    - idle: {}

    # do nothing but resolve rendering suffix
    - idle@tpl: |-
        {{ dukkha.Set "foo" "bar" }}

    after:matrix:failure:
    # run raw command (no env expansion) in /home (as specified by chdir)
    - cmd:
      - dukkha
      - render
      - foo.yml
      chdir: /home

    after:matrix:
    # define name for your task
    - name: foo

    # infinit loop is usually not what we want
    # use rendering suffix to do conditional next
    # - next: foo
    - idle@tpl: |-
        {{- dukkha.Set "foo" "done" -}}
      next@tpl: |-
        {{- if ne values.foo "done" -}}
          foo
        {{- end -}}

    after:success:
    # run another predefined task, all matrix
    - task: workflow:run(foo, {})

    # run commands/tasks when the task execution failed
    after:failure: []

    after: []
```
