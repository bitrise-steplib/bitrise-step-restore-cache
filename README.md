# Restore Cache (Beta)

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-restore-cache?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-restore-cache/releases)

Restores build cache using a cache key. This Step needs to be used in combination with **Save Cache**.

<details>
<summary>Description</summary>

Restores build cache using a cache key. This Step needs to be used in combination with **Save Cache**.

#### About key-based caching

Key-based caching is a concept where cache archives are saved and restored using a unique cache key. One Bitrise project can have multiple cache archives stored simultaneously, and the **Restore Cache Step** downloads a cache archive associated with the key provided as a Step input.

Caches can become outdated across builds when something changes in the project (for example, a dependency gets upgraded to a new version). In this case, a new (unique) cache key is needed to save the new cache contents. This is possible if the cache key is dynamic and changes based on the project state (for example, a checksum of the dependency lockfile is part of the cache key). If you use the same dynamic cache key when restoring the cache, the Step will download the most relevant cache archive available.

Key-based caching is platform-agnostic and can be used to cache anything by carefully selecting the cache key and the files/folders to include in the cache.

#### Templates

The Step requires a string key to use when downloading a cache archive. In order to always download the most relevant cache archive for each build, the cache key input can contain template elements. The Step evaluates the key template at runtime and the final key value can change based on the build environment or files in the repo.

The following variables are supported in cache keys input:

- `cache-key-{{ .Branch }}`: Current git branch the build runs on
- `cache-key-{{ .CommitHash }}`: SHA-256 hash of the git commit the build runs on
- `cache-key-{{ .Workflow }}`: Current Bitrise workflow name (eg. `primary`)
- `{{ .Arch }}-cache-key`: Current CPU architecture (`amd64` or `arm64`)
- `{{ .OS }}-cache-key`: Current operating system (`linux` or `darwin`)

Functions available in a template:

`checksum`: This function takes one or more file paths and computes the SHA256 [checksum](https://en.wikipedia.org/wiki/Checksum) of the file contents. This is useful for creating unique cache keys based on files that describe content to cache.

Examples of using `checksum`:
- `cache-key-{{ checksum "package-lock.json" }}`
- `cache-key-{{ checksum "**/Package.resolved" }}`
- `cache-key-{{ checksum "**/*.gradle*" "gradle.properties" }}`

`getenv`: This function returns the value of an environment variable or an empty string if the variable is not defined.

Examples of `getenv`:
- `cache-key-{{ getenv "PR" }}`
- `cache-key-{{ getenv "BITRISEIO_PIPELINE_ID" }}`

#### Key matching and fallback keys

The most straightforward use case is that a cache archive is downloaded and restored if the provided key matches a cache archive uploaded previously using the Save Cache Step. Stored cache archives are scoped to the Bitrise project. Builds can restore caches saved by any previous Workflow run on any Bitrise Stack.

It's possible to define more than one key in the cache keys input. You can specify additional keys by listing one key per line. The list is in priority order, so the Step will first try to find a match for the first key you provided, and if there is no cache stored for the key, it will move on to find a match for the second key (and so on).

In addition to listing multiple keys, each key can be a prefix of a saved cache key and still get a matching cache archive. For example, the key `my-cache-` can match an existing archive saved with the key `my-cache-a6a102ff`.

We recommend configuring the keys in a way that the first key is an exact match to a checksum key, and to use a more generic prefix key as a fallback:

```
inputs:
  key: |
    npm-cache-{{ checksum "package-lock.json" }}
    npm-cache-
```

#### Related steps

[Save cache](https://github.com/bitrise-steplib/bitrise-step-save-cache/)

</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

⚠️ **Beta status**: While this Step is in beta, everyone can use it without restrictions, quotas or costs.

### Examples

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


## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `key` | Keys used for restoring a cache archive. One cache key per line in priority order.  The key supports template elements for creating dynamic cache keys. These dynamic keys change the final key value based on the build environment or files in the repo in order to create new cache archives.  See the Step description for more details and examples. | required |  |
| `verbose` | Enable logging additional information for troubleshooting. | required | `false` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-restore-cache/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-restore-cache/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

**Note:** this step's end-to-end tests (defined in `e2e/bitrise.yml`) are working with secrets which are intentionally not stored in this repo. External contributors won't be able to run those tests. Don't worry, if you open a PR with your contribution, we will help with running tests and make sure that they pass.


Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
