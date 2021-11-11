# Define Same Kind Task in Different Files

Question: How can I have my tasks definitions distributed to individual files and run all tasks at project root?

## TL;DR

- Option 1: Add rendering suffix `@` to your task kind line. (e.g. `workflow:run@`)
- Option 2: Do not include these files, use patch spec (`!`). (e.g. `workflow:run@!`)

## Detailed Answers

### Option 1 - Make use of rendering suffix's special handling of inline map items

Task definitions are actually inline map items in dukkha config, when its field key has a rendering suffix, you are able to define same key multiple times in the same doc and merge their items together.

First, let's have a look at tasks in single file, sample valid task `foo` and `bar` of kind `workflow:run`:

```yaml
workflow:run:
- name: foo
- name: bar
```

But following doc is invalid due to there are duplicate keys `workflow:run` (yaml syntax error)

```yaml
workflow:run:
- name: foo

workflow:run:
- name: bar
```

Rendering suffix as its design purpose is to give get fields rendered after unmarshaling, so every fields with rendering suffix will survive even it's duplicated (since they can have different renderers), and if the target field (in compiled code) is a inline map, its items get merged together, so you can have follow valid dukkha config:

```yaml
# notice the `@` suffix
workflow:run@:
- name: foo

workflow:run@:
- name: bar
```

And when it comes to the multi-file scenario, this behavior stays the same as dukkha merges all files before resolving, so just include these files.

### Option 2 - Use Patch Spec

Given you have a task definition `foo` in `foo.yaml`:

```yaml
# filename: foo.yaml
workflow:run:
- name: foo
```

You want to run it along with root config (`.dukkha.yaml`) already having `bar` defined:

```yaml
# filename: .dukkha.yaml
workflow:run:
- name: bar
```

Then you can include that task `foo` using patch spec

```yaml
# filename: .dukkha.yaml

# NOTICE the suffix `@!`, `@` is the start of renering suffix, and `!` means doing patching, there are no renderer being used in this case so it's empty
workflow:run@!:
  # most of the case you would't like task `foo`'s fields get resolved at patch time, just disable resolving (defaults to true for common data patching)
  resolve: false

  value:
  - name: bar

  merge:
  - value@file: foo.yaml
    # disable resolving of task foo, for the same reason above
    resolve: false
```
