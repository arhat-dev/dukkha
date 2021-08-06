# Kubernetes Manifest Examples

Generate kubernetes manifests using `dukkha render`

```bash
# in same dir as this README.md
dukkha -c ./dukkha-config.yaml \
  render ./source \
  -o ./build/kubernetes-example

# then you can find generated manifests in ./build/kubernetes-example
```
