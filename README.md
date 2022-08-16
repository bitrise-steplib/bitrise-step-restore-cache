# Restore Cache (Beta)

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-restore-cache?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-restore-cache/releases)

Restores build cache using a cache key. This step needs to be used in combination with **Save Cache**.

<details>
<summary>Description</summary>

#### About key-based caching

Key-based caching is a concept where cache archives are saved and restored using a unique cache key. One Bitrise project can have multiple cache archives stored at the same time, and the **Restore Cache** step downloads a cache archive associated with the key provided as a step input.

Caches can become outdated across builds when something changes in the project (for example, a dependency gets upgraded to a new version). In this case, a new (unique) cache key is needed to save the new cache contents. This is possible if the cache key is dynamic and changes based on the project state (for example, a checksum of the dependency lockfile is part of the cache key). If the same dynamic cache key is used when restoring the cache, it will download the most relevant cache archive available.

Key-based caching is not platform specific and can be used to cache anything by carefully selecting the cache key and the files/folders to include in the cache.

#### Templates

The Step requires a string key to use when downloading a cache archive. In order to always download the most relevant cache archive for each build, the cache key input can contain template elements. The Step evaluates the cache key at runtime and the final key value can change based on the build environment or files in the repo.  

The following variables are supported in keys:

- `cache-key-{{ .Branch }}`: Current git branch the build runs on
- `cache-key-{{ .CommitHash }}`: SHA-256 hash of the git commit the build runs on
- `cache-key-{{ .Workflow }}`: Current Bitrise workflow name (eg. `primary`)
- `{{ .Arch }}-cache-key`: Current CPU architecture (`amd64` or `arm64`)
- `{{ .OS }}-cache-key`: Current operating system (`linux` or `darwin`)

Functions available in a template:

`checksum`: this function takes one or more file paths and computes the SHA256 checksum of the file contents. This is useful for creating unique cache keys based on files that describe the cached content.

Examples of using `checksum`:
- `cache-key-{{ checksum "package-lock.json" }}`
- `cache-key-{{ checksum "**/Package.resolved" }}`
- `cache-key-{{ checksum "**/*.gradle*" "gradle.properties" }}`

`getenv`: this function returns the value of an environment variable or an empty string if the variable is not defined.

Examples of `getenv`:
- `cache-key-{{ getenv "PR" }}`
- `cache-key-{{ getenv "BITRISEIO_PIPELINE_ID" }}`

#### Key matching and fallback keys

In the simplest case, a cache archive gets downloaded and restored if the provided key matches a cache archive uploaded previously using the Save Cache step. Stored cache archives are scoped to the Bitrise project. Builds can restore caches saved by any previous workflow run on any Bitrise Stack.

It's possible to define more than one key for the Step. By listing one key per line, the Step will try to look up the first key, and if there is no cache stored for the first key, it will look up the second key (and so on).

In addition to listing multiple keys, each key can be a prefix of a saved cache key and still get a matching cache archive. For example, the key `my-cache-` can match an existing archive saved with the key `my-cache-a6a102ff`.

It's a common practice to configure the keys in a way that the first key is an exact match to a checksum key, and a more generic prefix key is used as a fallback:

```
inputs:
  key: |
    npm-cache-{{ checksum "package-lock.json" }}
    npm-cache-
```

</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `key` | One cache key per line in priority order.  The key supports template elements for creating dynamic cache keys. These dynamic keys change the final key value based on the build environment or files in the repo in order to create new cache archives.  See the step description for more details and examples. | required |  |
| `verbose` | Enable logging additional information for troubleshooting. | required | `false` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-restore-cache/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-restore-cache/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
