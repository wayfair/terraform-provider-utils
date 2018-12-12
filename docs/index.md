# Terraform Provider Utilities

This repository contains utility functions, packages used in the other
in-house developed Terraform providers. Common and core functionality is
commited to this repository and used in other codebases.

This project is a component of Project Argo.

## Project Info

This project is developed, owned, and maintained by the SRE - Orchestration
pod at Wayfair.

## Getting Started

To include this repository in your project, you first need to have a
dependency management tool in place. For Terraform provider development, we
use [`dep`](https://golang.github.io/dep/) as it is widely adopted and part of
the official golang tooling. If your repository is not already using `dep`,
follow the getting started guide.

This process assumes you are using `dep` for dependency management.

### Include the Utility Repository

Update your `Gopkg.toml` file so it includes the repository:

```
# Gopkg.toml
[[constraint]]
      name     = "github.com/wayfair/terraform-provider-utils"
      # constraint to a specific version, release here:
      version = "=1.0.0"

```

### Include the Desired Packages in Source

Now you can include the package(s) to your source files. Notice the `.git`
in the import. See the `NOTE` in the above section.

```
// example/main.go
//
// Example repository, main.go source file

package example

import (
    "github.com/wayfair/terraform-provider-utils/example"
)

func main() {
  // calling the 'Foo' function from the 'example' package
  example.Foo()
}
```

### Update Dependencies

Last, inform `dep` to update your dependencies:

```
$> dep ensure
```

Dep will constrain the dependencies, pull in the used packages, and manage
them in `vendor/`.
