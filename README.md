# Terraform Provider Utilities

This repository contains utility functions, packages used in the other
in-house developed Terraform providers. Common and core functionality is
commited to this repository and used in other codebases.

This project is a component of Project Argo.

## Project Info

This project is developed, owned, and maintained by the SRE - Orchestration
pod at Wayfair.

This repository uses [`mkdocs`](https://www.mkdocs.org/) for documentation
and [`dep`](https://golang.github.io/dep/) for dependency management.
Dependencies are tracked as part of the repository.

## Requirements:

- [Terraform](https://www.terraform.io/downloads.html) >= 0.10.x
- [Golang](https://golang.org/doc/install) >= 1.8
- [Dep](https://golang.github.io/dep/docs/installation.html) >= 0.4.1
- [GNU Make](https://www.gnu.org/software/make/) >= 4.2.1

Follow the setup instructions provided on the install sections of their
respective websites. Windows environments should have a \*nix-style terminal
emulator installed such as [Cygwin](https://www.cygwin.com/) to be compatible
with the `makefile`.

## Repository Setup

After installing and configuring the toolchain listed in the `Requirements`
section:

1. Clone the repository with `ssh`:

    ```sh
    $ mkdir -p "${GOPATH}/src/github.com/wayfair"
    $ cd "${GOPATH}/src/github.com/wayfair"
    $ git clone git@github.com:wayfair/terraform-provider-utils.git
    ```

2. Include the utility repo in your provider:

    Update your `Gopkg.toml` to include the repo as a dependency:

    ```
    # Gopkg.toml

    [[constraint]]
      name = "github.com/wayfair/terraform-provider-utils"
      # constraint to a specific version, release here:
      version = "=1.0.0"
    ```

    Include the needed package(s) to your source files:

    ```
    // resource_example_foo.go
    //
    // Example provider, foo resource

    package example

    import (
        // Import a specific subpackage from the repo, in this case 'example'.
        // Here, we don't define a pseudonym, so the package is accessible by
        // 'example'.
        "github.com/wayfair/terraform-provider-utils/example"
        // Here we are importing the 'example2' package under the pseudonym
        // 'ex'. It is available in this file as 'ex' not 'example2'
        ex "github.com/wayfair/terraform-provider-utils/example2"
    )

    // ...
    ```

    Update your repository to include the constraint with `dep`:

    ```
    $> dep ensure
    ```

## Documentation

This repository uses [`mkdocs`](https://www.mkdocs.org/) for documentation.
Repository documentation can be found in the `doc` directory. Follow the
installation instructions on [`mkdocs`](https://www.mkdocs.org/#installation)
to get started.

The `Makefile` exposes a `godoc` target which can be used to generate and save
the project's Godoc to the local filesystem in `docs/godoc`. These pages are
used by `mkdocs` to generate the full project documentation. The `godoc` target
only saves the necessary package documentation for this repository and does
save the entire webroot.

To generate and view the entire repository's documentation:

```
$> make godoc
Generating godoc to docs/godoc...
Creating docs/godoc
godoc PID: [5084]
Sleeping while godoc initializes...
Downloading pages...
<...output truncated...>
done.
Killing godoc process [5084]

$> mkdocs serve
INFO    -  Building documentation...
INFO    -  Cleaning site directory
[I 160402 15:50:43 server:271] Serving on http://127.0.0.1:8000
[I 160402 15:50:43 handlers:58] Start watching changes
[I 160402 15:50:43 handlers:60] Start detecting changes
```

The documentation can then be viewed by accessing localhost in your favorite
browser or viewport.

Cleaning up the generated `godoc` can be done with the `clean-godoc` target.

```
$> make clean-godoc
Cleaning godoc files...
```
