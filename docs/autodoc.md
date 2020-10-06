# Autodoc

The `autodoc` package contains the engine to automatically generate `mkdocs`
style documentation and configuration files for a Terraform provider. The
engine parses the schema maps and schema definitions inside the provider,
resources, and data sources to generate this documentation. In addition,
`autodoc` provides a basic metadata/tagging feature to allow for more
fine-grained details and override certain behaviors.

## Command Line Arguments

The autodoc engine recognizes the following arguments:

* `-help` Show usage information and exit
* `-provider` The name of the Terraform provider. Defaults to `"Terraform Provider"`.
* `-root` The root output directory. `mkdocs.yml` will be generated here
    and the documentation directory will fall under this path. Defaults to
    the current working directory.
* `-docs-dir` Name of the documentation directory for `mkdocs.yml`. It will
    be placed under the value supplied for `-root`. The `mkdocs.yml` `docs_dir`
    will be set to this value. Defaults to `docs`.
* `-templates-dir` Path to the templates directory. Template files used to
    generate the documentation will be searched from this directory. Defaults
    to `templates`.
* `-template-ext` File extension for templates. Defaults to `.template`.

## Output Files

The following files are generated. Assuming `/` is the project root and
`docs` is the documentation directory.

* `/mkdocs.yml` The `mkdocs` configuration file
* `/docs/index.md` Documentation index. This is the "landing page" to your
    documentation.
* `/docs/godoc.md` A container page for the project's `godoc`
* `/docs/resources/*.md` A documentation file is generated for each resource.
    The file name will correspond to the name of the resource in the
    `Provider.Schema.ResourcesMap`.
* `/docs/datasources/*.md` A documentation file is generated for each data
    source. The file name will correspond to the name of the data source in
    the `Provider.Schema.DataSourcesMap`.

## Templates

`autodoc` utilizes the `text/template` package from golang stdlib in order
to generate documentation. `autodoc` makes the following template associations
to produce output:

* `mkdocs.yml.template` => `mkdocs.yml`
* `godoc.md.template` => `godoc.md`
* `resource.md.template` => `docs/resources/*.md`
* `datasource.md.template` => `docs/datasources/*.md`

These 4 templates are required for the engine to function properly.

### Template Data

`autodoc` makes the following data available to in your templates:

#### `mkdocs.yml`

* `DocsDir` The name of the documentation directory. This is the value
    supplied to `-docs-dir`.
* `Resources` The list of resource names. These are the keys to the
    `Provider.Schema.ResourcesMap`.
* `DataSources` The list of data source names. These are the keys to the
    `Provider.Schema.DataSourcesMap`.

#### Provider, Resources, & Data Sources Documentation

* `Constants` A map of exposed constants:
    * `TypeProvider` Denotes a provider schema
    * `TypeResource` Denotes a resource schema
    * `TypeDataSource` Denotes a data source schema
* `SchemaType` The type of schema that is being documented. This will be
    one of the `TypeXxx` constants in `Constants`.
* `Name` The name of the provider, resource, or data source. If this is a
    provider, the name will be the value supplied to `-provider`. Otherwise,
    this will be the key to `Provider.Schema.ResourcesMap` or
    `Provider.Schema.DataSourcesMap` for resources and data sources
    respectively.
* `Meta` The parsed metadata information for this schema. `Meta` has the
    following attributes available:
    * `Uncreatable` Boolean, whether or not this resource supports create
    * `Undeletable` Boolean, whether or not this resource supports delete
    * `Immutable` Boolean, whether or not this resource supports update
    * `Summary` The parsed summary information for this schema
* `Attributes` List of exported schema attributes. Each attribute has the
    following properties available:
    * `Name` The name of the attribute. This is the key to
        `Provider.Schema` for the provider, and `Resource.Schema` for
        resources and data sources.
    * `Type` The parsed, markdown escaped, formatted type. For simple types
        like `schema.TypeInt` it will be the markdown string `schema.TypeInt`.
        For complex types like sets, lists, or maps it will be an escaped string
        indicating the element type as well. For example
        `schema.TypeSet of schema.TypeInt`.
    * `Description` The description of the attribute with metadata tags stripped
* `Arguments` List of schema arguments. Each argument has the following
    properties available:
    * `Name` The name of the attribute. This is the key to
        `Provider.Schema` for the provider, and `Resource.Schema` for
        resources and data sources.
    * `Type` The parsed, markdown escaped, formatted type. For simple types
        like `schema.TypeInt` it will be the markdown string `schema.TypeInt`.
        For complex types like sets, lists, or maps it will be an escaped string
        indicating the element type as well. For example
        `schema.TypeSet of schema.TypeInt`.
    * `Description` The description of the attribute with metadata tags stripped
    * `Example` An example value for this argument. If no example was provided,
        it will be the empty string.
    * `Optional` Boolean, whether or not this an optional argument.
    * `ForceNew` Boolean, whether or not this argument forces a destroy and
        recreation of the resource.
    * `ConflictsWith` List of any conflicting arguments

## Metadata Attributes and Tagging

`autodoc` supports metadata and tagging, much like `javadoc`, `sphinx`,
and other documentation tools. Tags are read and parsed from the
`Description` of a schema.

Some tags allow "arguments" which should follow the tag separated by a
whitespace character on either side. Tags are separated by a whitespace
character.

### The Meta Attribute

The meta attribute, `__meta__` is a Schema that is added to a schema map
(ie: a `Provider.Schema`, `Resource.Schema`, etc) and applies resource-level
meta information.

The meta attribute should be defined as `Computed: true` to ensure it cannot
be affected or treated as an argument to the HCL.

* `@SUMMARY value` Provides a summary of the resource. This is displayed in the
    `Description` section of a resource, data source.
* `@NOTCREATABLE` This object cannot be created. This is sometimes the case when
    Terraform manages system-generated objects. By default, `autodoc` assumes
    the resource can be created.
* `@NOTDELETABLE` This object cannot be deleted. This is sometimes the case when
    Terraform manages system-generated objects. By default, `autodoc` assumes
    the resource can be deleted.
* `@IMMUTABLE` This object cannot be updated. This is sometimes the case when
    Terraform manages system-generated objects. By default, `autodoc` assumes
    the resource can be updated.

Example utilization:

```
package example

// ...

func resourceExampleFoo() *schema.Resource {
  return &schema.Resource{

    // ...

    Schema: map[string]*schema.Schema{

      "__meta__": &schema.Schema{
        Type:        schema.TypeBool,
        Computed:    true,
        Description: "@SUMMARY This is a Foo. Foo are example objects. @IMMUTABLE",
      },

      // ...

    },

    // ...

  }
}
```

### Attribute and Argument Tags

`autodoc` makes a distinction between arguments and attributes. Attributes are
the exported properties of a resource that can be accessed by other objects
in the HCL. This is accomplished by calling `ResourceData.Set(key, val)`.

Arguments are provided to the resource block in HCL and are needed
to properly generate, modify the resource. An entry in the `Provider.Schema` or
`Resource.Schema` can be both an argument and an exported attribute, just an
attribute (ie: a `Computed` entry), or just an argument (ie: an `Optional` or
`Required` entry that is not saved with `ResourceData.Set(key, val)`).

`autodoc` can parse the following tags in the description of an argument
or attribute:

* `@UNEXPORTED` This property is not exported and available to other resources,
    data sources. By default, `autodoc` assumes the property is exported. This
    will override that behavior.
* `@EXAMPLE value` Provides an example value for an argument. This will show up
    in the `Examples` section of the documentation to show how to properly
    use this resouce/data source.

## Getting Started

This tool was designed to be plug and play with little disruption. However,
there is some work to get things setup properly:

* Include the util repo and the `autodoc` package as part of your dependencies
    if you have not already done so. There is information on the
    [index](index.md) that will help you with that.

* Set up your templates directory and write your documentation templates (or
    copy them from another repository to get started). Look at the `Template
    Data` section above to get a reference on what is exposed to you in your
    templates.

* Create the `autodoc` binary specific to your project.

    Each provider will need to create a simple `main()` routine that passes
    the provider refernce to the engine. Following go conventions, you will
    create a `/cmd/autodoc/main.go` in your repository.

    Golang does not allow for late linking and produces static binaries. As
    a design decision, we would either have to include all providers as part
    of this package, or have all providers create a binary that uses the
    engine. We opted for the latter.

    The contents of this file is rather trivial:

```
// Package main contains the main goroutine for the autodoc command-line
// application.  This application uses the autodoc engine to create mkdocs
// style documentation for the Terraform provider.
//
// For more information on the autodoc tool, its arguments, etc see:
// pkg/github.com/wayfair/terraform-provider-utils/autodoc
package main

import (
  "fmt"
  "os"

  // include your provider here:
  "github.com/jsmith/terraform-provider-example/example
  // include the autodoc package:
  autodoc "github.com/wayfair/terraform-provider-utils/autodoc"

  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func main() {

  // Use the provider function to get information on the provider's schema,
  // resources, and data sources. The Provider() function returns a
  // *terraform.ResourceProvider (interface) which will need to be type asserted
  // to a *schema.Provider (struct)
  resourceProvider := example.Provider()
  provider := resourceProvider.(*schema.Provider)

  // Start the autodoc engine, this will return a slice of non-nil errors (if
  // any were encountered during the run)
  errors := autodoc.Document(provider)
  if len(errors) != 0 {
    for _, err := range errors {
      fmt.Println(err)
    }
    os.Exit(autodoc.ExitError)
  }
  os.Exit(autodoc.ExitSuccess)

}
```

* Update your resources, data sources, and provider with the necessary
    information. All of metadata and tags should be the `Description` of your
    Schemas. If you have documentation in the form of comments, move them to
    the description:

```
// The name of the Foo
"name": &schema.Schema{
  Type:     schema.TypeString,
  Required: true,
},
```

Becomes:

```
"name": &schema.Schema{
  Type:        schema.TypeString,
  Required:    true,
  Description: "The name of the Foo",
},
```

* Add a makefile target to generate your Godoc for your project. There are
    different approaches to this problem, you can look at other providers
    for examples on how to do this.

* Optionally, update your `.gitignore` to not track documentation as part
    of the project. If you pipeline is using the `autodoc` tool to create
    and publish the documentation on push, then this may of some use to keep
    the codebase clean.

```
# .gitignore

# autodoc files
mkdocs.yml
docs/*.md
docs/godoc/**
docs/resources/*.md
docs/datasources/*.md
```
