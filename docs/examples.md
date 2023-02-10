
Check out [Workflow Recipes](https://github.com/bitrise-io/workflow-recipes#-key-based-caching-beta) for platform-specific examples!

#### Restore and save cache using a key that includes a checksum

```yaml
steps:
- restore-cache@1:
    inputs:
    - key: npm-cache-{{ checksum "package-lock.json" }}

# Build steps

- save-cache@1:
    inputs:
    - key: npm-cache-{{ checksum "package-lock.json" }}
    - paths: node_modules
```

#### Use fallback key when exact cache match is not available

```yaml
steps:
- restore-cache@1:
    inputs:
    - key: |-
        npm-cache-{{ checksum "package-lock.json" }}
        npm-cache-
```

The Step will look up the first key (`npm-cache-a982ee8f` for example). If there is no exact match for this key (because there is only a cache archive for `npm-cache-233ad571`), then it will look up any cache archive whose key starts with `npm-cache-`.

#### Separate caches for each OS and architecture

Cache is not guaranteed to work across different Bitrise Stacks (different OS or same OS but different CPU architecture). If a workflow runs on different stacks, it's a good idea to include the OS and architecture in the cache key:

```yaml
steps:
- restore-cache@1:
    inputs:
    - key: |-
        {{ .OS }}-{{ .Arch }}-npm-cache-{{ checksum "package-lock.json" }}
        {{ .OS }}-{{ .Arch }}-npm-cache-
```
