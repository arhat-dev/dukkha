# Arbitrary File Rendering

## Run

```bash
# with jq
dukkha render ./srouce -f json | jq '."foo"' -r
# or with yq
dukkha render ./source | yq '."foo"' -r
```
