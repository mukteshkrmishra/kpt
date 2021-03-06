// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by "mdtogo"; DO NOT EDIT.
package cfgdocs

var READMEShort = `Examine and modify configuration files`
var READMELong = `
Programmatically print and modify raw json or yaml Resource Configuration

| Command        | Description                                   |
|----------------|-----------------------------------------------|
| [annotate]     | set an annotation on one or more Resources    |
| [cat]          | print resources                               |
| [count]        | print resource counts                         |
| [create-setter]| create or modify a field setter               |
| [create-subst] | create or modify a field substitution         |
| [fmt]          | format configuration files                    |
| [grep]         | find resources by field value                 |
| [list-setters] | print available field setters                 |
| [set]          | set one or more field values                  |
| [tree]         | print resources as a tree                     |

**Data Flow**: local configuration or stdin -> kpt [cfg] -> local configuration or stdout

| Configuration Read From | Configuration Written To |
|-------------------------|--------------------------|
| local files or stdin    | local files or stdout    |
`
var READMEExamples = `
    # print the raw package contents
    $ kpt cfg cat helloworld

    # print the package using tree based structure
    $ kpt cfg tree helloworld --name --image --replicas
    helloworld
    ├── [deploy.yaml]  Deployment helloworld-gke
    │   ├── spec.replicas: 5
    │   └── spec.template.spec.containers
    │       └── 0
    │           ├── name: helloworld-gke
    │           └── image: gcr.io/kpt-dev/helloworld-gke:0.1.0
    └── [service.yaml]  Service helloworld-gke

    # only print Services
    $ kpt cfg grep "kind=Service" helloworld | kpt cfg tree --name --image --replicas
    .
    └── [service.yaml]  Service helloworld-gke

    #  list available setters
    $ kpt cfg list-setters helloworld replicas
        NAME          DESCRIPTION        VALUE    TYPE     COUNT   SETBY
      replicas   'helloworld replicas'   5       integer   1

    # set a high-level knob
    $ kpt cfg set helloworld replicas 3
    set 1 fields
`

var AnnotateShort = `Set an annotation on one or more Resources`
var AnnotateLong = `
  DIR:
    Path to local directory.
`
var AnnotateExamples = `
    # set an annotation on all Resources: 'key: value'
    kpt cfg annotate DIR --kv key=value

    # set an annotation on all Service Resources
    kpt cfg annotate DIR --kv key=value --kind Service

    # set an annotation on the foo Service Resource only
    kpt cfg annotate DIR --kv key=value --kind Service --name foo

    # set multiple annotations
    kpt cfg annotate DIR --kv key1=value1 --kv key2=value2
`

var CatShort = `Print resources`
var CatLong = `
    kpt cfg cat DIR

  DIR:
    Path to local directory.
`
var CatExamples = `
    # print Resource config from a directory
    kpt cfg cat my-dir/
`

var CountShort = `Print resource counts`
var CountLong = `
    kpt cfg count [DIR]

  DIR:
    Path to local directory.
`
var CountExamples = `
    # print Resource counts from a directory
    kpt cfg count my-dir/

    # print Resource counts from a cluster
    kubectl get all -o yaml | kpt cfg count
`

var CreateSetterShort = `Create or modify a field setter`
var CreateSetterLong = `
    kpt cfg create-setter DIR NAME VALUE

  DIR

    A directory containing Resource configuration.
    e.g. hello-world/

  NAME

    The name of the substitution to create.  This is both the name that will be given
    to the *set* command, and that will be referenced by fields.
    e.g. replicas

  VALUE

    The new value of the setter.
    e.g. 3

#### Field Setters

Field setters are OpenAPI definitions that define how fields may be modified programmatically
using the *set* command.  The OpenAPI definitions for setters are defined in a Kptfile
and referenced by fields which they set through an OpenAPI reference as a line comment
(e.g. # {"$ref":"#/definitions/..."}).

Setters may be manually created by editing the Kptfile, or programmatically created using the
` + "`" + `create-setter` + "`" + ` command.  The ` + "`" + `create-setter` + "`" + ` command will 1) create a new OpenAPI definition
for a setter in the Kptfile, and 2) identify all fields matching the setter value and create
an OpenAPI reference to the setter for each.

    # create or update a setter named replicas
    kpt create-setter hello-world/ replicas 3

Example setter definition in a Kptfile:

	openAPI:
	  definitions:
	    io.k8s.cli.setters.replicas:
	      x-k8s-cli:
	        setter:
	          name: "replicas"
	          value: "3"

This setter is named "replicas" and can be provided to the *set* command to change
all fields which reference it to the setter's value.

Example setter referenced from a field in a configuration file:

	kind: Deployment
	metadata:
	  name: foo
	spec:
	  replicas: 3  # {"$ref":"#/definitions/io.k8s.cli.setters.replicas"}

Setters may have types specified which ensure that the configuration is always serialized
correctly as yaml 1.1 -- e.g. if a string field such as an annotation or arg has the value
"on", then it would need to be quoted otherwise it will be parsed as a bool by yaml 1.1.

A type may be specified using the --type flag, and accepts string,integer,boolean as values.
The resulting OpenAPI definition looks like:

    # create or update a setter named version which sets the "version" annotation
    kpt create-setter hello-world/ version 3 --field "annotations.version" --type string

	openAPI:
	  definitions:
	    io.k8s.cli.setters.version:
	      x-k8s-cli:
	        setter:
	          name: "version"
	          value: "3"
	      type: string

And the configuration looks like:

	kind: Deployment
	metadata:
	  name: foo
	  annotations:
	    version: "3" # {"$ref":"#/definitions/io.k8s.cli.setters.version"}

Setters may be configured to accept enumeration values which map to different values set
on the fields.  For example setting cpu resources to small, medium, large -- and mapping
these to specific cpu values.  This may be done by manually modifying the Kptfile openAPI
definitions as shown here:

	openAPI:
	  definitions:
	    io.k8s.cli.setters.cpu:
	      x-k8s-cli:
	        setter:
	          name: "cpu"
	          value: "small"
	          # enumValues will replace the user provided key with the
	          # map value when setting fields.
	          enumValues:
	            small: "0.5"
	            medium: "2"
	            large: "4"

And the configuration looks like:

	kind: Deployment
	metadata:
	  name: foo
	spec:
	  template:
	    spec:
	      containers:
	      - name: foo
	    resources:
	      requests:
	        cpu: "0.5" # {"$ref":"#/definitions/io.k8s.cli.setters.cpu"}
`
var CreateSetterExamples = `
    # create a setter called replicas for fields matching "3"
    kpt cfg create-setter DIR/ replicas 3

    # scope creating setter references to a specified field
    kpt cfg create-setter DIR/ replicas 3 --field "replicas"

    # scope creating setter references to a specified field path
    kpt cfg create-setter DIR/ replicas 3 --field "spec.replicas"

    # create a setter called replicas with a description and set-by
    kpt cfg create-setter DIR/ replicas 3 --set-by "package-default" \
        --description "good starter value"

    # scope create a setter with a type.  the setter will make sure the set fields
    # always parse as strings with a yaml 1.1 parser (e.g. values such as 1,on,true wil
    # be quoted so they are parsed as strings)
    # only the final part of the the field path is specified
    kpt cfg create-setter DIR/ app nginx --field "annotations.app" --type string
`

var CreateSubstShort = `Create or modify a field substitution`
var CreateSubstLong = `
    kpt cfg create-subst DIR NAME VALUE --pattern PATTERN --value MARKER=SETTER

  DIR

    A directory containing Resource configuration.
    e.g. hello-world/

  NAME

    The name of the substitution to create.  This is simply the unique key which is referenced
    by fields which have the substitution applied.
    e.g. image-substitution

  VALUE

    The current value of the field that will have PATTERN substituted.
    e.g. nginx:1.7.9

  PATTERN

    A string containing one or more MARKER substrings which will be substituted
    for setter values.  The pattern may contain multiple different MARKERS,
    the same MARKER multiple times, and non-MARKER substrings.
    e.g. IMAGE_SETTER:TAG_SETTER

#### Field Substitutions

Field substitutions are OpenAPI definitions that define how fields may be modified programmatically
using the *set* command.  The OpenAPI definitions for substitutions are defined in a Kptfile
and referenced by fields which they set through an OpenAPI reference as a line comment
(e.g. # {"$ref":"#/definitions/..."}).

Substitutions may be manually created by editing the Kptfile, or programmatically created using the
` + "`" + `create-subst` + "`" + ` command.  The ` + "`" + `create-subst` + "`" + ` command will 1) create a new OpenAPI definition
for a substitution in the Kptfile, and 2) identify all fields matching the provided value and create
an OpenAPI reference to the substitution for each.

Field substitutions are computed by substituting setter values into a pattern.  They are
composed of 2 parts: a pattern and a list of values.

- The pattern is a string containing markers which will be replaced with 1 or more setter values.
- The values are pairs of markers and setter references.  The *set* command retrieves the values
  from the referenced setters, and replaces the markers with the setter values.
 
**The referenced setters MAY exist before creating the substitution, in which case the
existing setters are used instead of recreated.**

    # create or update a substitution + 2 setters
    # the substitution is derived from concatenating the image and tag setter values
    kpt create-subst hello-world/ image-tag nginx:1.7.9 \
      --pattern IMAGE_SETTER:TAG_SETTER \
      --value IMAGE_SETTER=image \
      --value TAG_SETTER=tag

If create-subst cannot infer the setter values from the VALUE + --pattern, and the setters
do not already exist, then it will throw and error, and the setters must be manually created
beforehand.

Example setter and substitution definitions in a Kptfile:

	openAPI:
	  definitions:
	    io.k8s.cli.setters.image:
	      x-k8s-cli:
	        setter:
	          name: "image"
	          value: "nginx"
	    io.k8s.cli.setters.tag:
	      x-k8s-cli:
	        setter:
	          name: "tag"
	          value: "1.7.9"
	    io.k8s.cli.substitutions.image-value:
	      x-k8s-cli:
	        substitution:
	          name: image-value
	          pattern: IMAGE_SETTER:TAG_SETTER
	          values:
	          - marker: IMAGE_SETTER
	            ref: '#/definitions/io.k8s.cli.setters.image'
	          - marker: TAG_SETTER
	            ref: '#/definitions/io.k8s.cli.setters.tag'

This substitution defines how a field value may be produced from the setters ` + "`" + `image` + "`" + ` and ` + "`" + `tag` + "`" + `
by replacing the pattern substring *IMAGE_SETTER* with the value of the ` + "`" + `image` + "`" + ` setter, and
replacing the pattern substring *TAG_SETTER* with the value of the ` + "`" + `tag` + "`" + ` setter.  Any time
either the ` + "`" + `image` + "`" + ` or ` + "`" + `tag` + "`" + ` values are changed via *set*, the substitution value will be
re-calculated for referencing fields.

Example substitution reference from a field in a configuration file:

	kind: Deployment
	metadata:
	  name: foo
	spec:
	  template:
	    spec:
	      containers:
	      - name: nginx
	        image: nginx:1.7.9 # {"$ref":"#/definitions/io.k8s.cli.substitutions.image-value"}

The ` + "`" + `image` + "`" + ` field has a OpenAPI reference to the ` + "`" + `image-value` + "`" + ` substitution definition.  When
the *set* command is called, for either the ` + "`" + `image` + "`" + ` or ` + "`" + `tag` + "`" + ` setter, the substitution will
be recalculated, and the ` + "`" + `image` + "`" + ` field updated with the new value.

**Note**: when setting a field through a substitution, the names of the setters are used
*not* the name of the substitution.  The name of the substitution is *only used in field
references*.
`
var CreateSubstExamples = `
    # Automatically create setters when creating the substitution, inferring the setter
    # values.
    #
    # 1. create a substitution derived from 2 setters.  The user will never call the
    #    substitution directly, instead it will be computed when the setters are used.
    kpt cfg create-subst DIR/ image-tag nginx:v1.7.9 \
      --pattern IMAGE_SETTER:TAG_SETTER \
      --value IMAGE_SETTER=nginx \
      --value TAG_SETTER=v1.7.9

    # 2. update the substitution value by setting one of the 2 setters it is computed from
    kpt cfg set tag v1.8.0


    # Manually create setters and substitution.  This is preferred to configure the setters
    # with a type, description, set-by, etc.
    #
    # 1. create the setter for the image name -- set the field so it isn't referenced
    kpt cfg create-setter DIR/ image nginx --field "none" --set-by "package-default"

    # 2. create the setter for the image tag -- set the field so it isn't referenced
    kpt cfg create-setter DIR/ tag v1.7.9 --field "none" --set-by "package-default"

    # 3. create the substitution computed from the image and tag setters
    kpt cfg create-subst DIR/ image-tag nginx:v1.7.9 \
      --pattern IMAGE_SETTER:TAG_SETTER \
      --value IMAGE_SETTER=nginx \
      --value TAG_SETTER=v1.7.9

    # 4. update the substitution value by setting one of the setters
    kpt cfg set tag v1.8.0
`

var FmtShort = `Format configuration files`
var FmtLong = `
Fmt will format input by ordering fields and unordered list items in Kubernetes
objects.  Inputs may be directories, files or stdin, and their contents must
include both apiVersion and kind fields.

- Stdin inputs are formatted and written to stdout
- File inputs (args) are formatted and written back to the file
- Directory inputs (args) are walked, each encountered .yaml and .yml file
  acts as an input

For inputs which contain multiple yaml documents separated by \n---\n,
each document will be formatted and written back to the file in the original
order.

Field ordering roughly follows the ordering defined in the source Kubernetes
resource definitions (i.e. go structures), falling back on lexicographical
sorting for unrecognized fields.

Unordered list item ordering is defined for specific Resource types and
field paths.

- .spec.template.spec.containers (by element name)
- .webhooks.rules.operations (by element value)
`
var FmtExamples = `
	# format file1.yaml and file2.yml
	kpt cfg fmt file1.yaml file2.yml

	# format all *.yaml and *.yml recursively traversing directories
	kpt cfg fmt my-dir/

	# format kubectl output
	kubectl get -o yaml deployments | kpt cfg fmt

	# format kustomize output
	kustomize build | kpt cfg fmt
`

var GrepShort = `Find resources by field value`
var GrepLong = `
    kpt cfg grep QUERY DIR

  QUERY:
    Query to match expressed as 'path.to.field=value'.
    Maps and fields are matched as '.field-name' or '.map-key'
    List elements are matched as '[list-elem-field=field-value]'
    The value to match is expressed as '=value'
    '.' as part of a key or value can be escaped as '\.'

  DIR:
    Path to local directory.
`
var GrepExamples = `
    # find Deployment Resources
    kpt cfg grep "kind=Deployment" my-dir/

    # find Resources named nginx
    kpt cfg grep "metadata.name=nginx" my-dir/

    # use tree to display matching Resources
    kpt cfg grep "metadata.name=nginx" my-dir/ | kpt cfg tree

    # look for Resources matching a specific container image
    kpt cfg grep "spec.template.spec.containers[name=nginx].image=nginx:1\.7\.9" my-dir/ | kpt cfg tree

###

[tutorial-script]: ../gifs/cfg-grep.sh`

var ListSettersShort = `List configured field setters`
var ListSettersLong = `
    kpt cfg list-setters DIR [NAME]

  DIR

    A directory containing a Kptfile.

  NAME

    Optional.  The name of the setter to display.
`
var ListSettersExamples = `
    # list the setters in the hello-world package
    kpt cfg list-setters hello-world/
      NAME     VALUE    SET BY    DESCRIPTION   COUNT  
    replicas   4       isabella   good value    1   

###

[tutorial-script]: ../gifs/cfg-set.sh`

var SetShort = `Set one or more field values using setters`
var SetLong = `
    kpt cfg set DIR NAME VALUE

  DIR

    A directory containing Resource configuration.
    e.g. hello-world/

  NAME

    The name of the setter
    e.g. replicas

  VALUE

    The new value to set on fields
    e.g. 3

#### Setters

The *set* command modifies configuration fields using setters defined as OpenAPI definitions
in a Kptfile.  Setters are referenced by fields using line commands on the fields.  Fields
referencing a setter will have their value modified to match the setter value when the *set*
command is called.

If multiple fields may reference the same setter, all of the field's values will be
changed when the *set* command is called for that setter.

The *set* command must be run on a directory containing a Kptfile with setter definitions.
The list of setters configured for a package may be found using ` + "`" + `kpt cfg list-setters` + "`" + `.

    kpt cfg set hello-world/ replicas 3

Example setter definition in a Kptfile:

	openAPI:
	  definitions:
	    io.k8s.cli.setters.replicas:
	      x-k8s-cli:
	        setter:
	          name: "replicas"
	          value: "3"

This setter is named "replicas" and can be provided to the *set* command to change
all fields which reference it to the setter's value.

Example setter referenced from a field in a configuration file:

	kind: Deployment
	metadata:
	  name: foo
	spec:
	  replicas: 3  # {"$ref":"#/definitions/io.k8s.cli.setters.replicas"}

#### Description

Setters may have a description of the current value.  This may be defined along with
the value by specifying the ` + "`" + `--description` + "`" + ` flag.

#### SetBy

Setters may record who set the current value.  This may be defined along with the
value by specifying the ` + "`" + `--set-by` + "`" + ` flag.  If unspecified the current value for
set-by will be cleared from the setter.

#### Substitutions

Substitutions define field values which may be composed of one or more setters substituted
into a string pattern.  e.g. setting only the tag portion of the ` + "`" + `image` + "`" + ` field.

Anytime set is called for a setter used by a substitution, it will also modify the fields
referencing that substitution.

See ` + "`" + `kpt cfg create-subst` + "`" + ` for more information on substitutions.
`
var SetExamples = `
    # set replicas to 3 using the 'replicas' setter
    kpt cfg set hello-world/ replicas 3

    # set the replicas to 5 and include a description of the value
    kpt cfg set hello-world/ replicas 5 --description "need at least 5 replicas"

    # set the replicas to 5 and record who set this value
    kpt cfg set hello-world/ replicas 5 --set-by "mia"

    # set the tag portion of the image field to '1.8.1' using the 'tag' setter
    # the tag setter is referenced as a value by a substitution in the Kptfile
    kpt cfg set hello-world/ tag 1.8.1

###

[tutorial-script]: ../gifs/cfg-set.sh`

var TreeShort = `Print resources as a tree`
var TreeLong = `
kpt cfg tree may be used to print Resources in a directory or cluster, preserving structure

Args:

  DIR:
    Path to local directory directory.

Resource fields may be printed as part of the Resources by specifying the fields as flags.

kpt cfg tree has build-in support for printing common fields, such as replicas, container images,
container names, etc.

kpt cfg tree supports printing arbitrary fields using the '--field' flag.

By default, kpt cfg tree uses Resource graph structure if any relationships between resources (ownerReferences)
are detected, as is typically the case when printing from a cluster. Otherwise, directory graph structure is used. The
graph structure can also be selected explicitly using the '--graph-structure' flag.
`
var TreeExamples = `
    # print Resources using directory structure
    kpt cfg tree my-dir/

    # print replicas, container name, and container image and fields for Resources
    kpt cfg tree my-dir --replicas --image --name

    # print all common Resource fields
    kpt cfg tree my-dir/ --all

    # print the "foo"" annotation
    kpt cfg tree my-dir/ --field "metadata.annotations.foo"

    # print the "foo"" annotation
    kubectl get all -o yaml | kpt cfg tree \
      --field="status.conditions[type=Completed].status"

    # print live Resources from a cluster using owners for graph structure
    kubectl get all -o yaml | kpt cfg tree --replicas --name --image

    # print live Resources with status condition fields
    kubectl get all -o yaml | kpt cfg tree \
      --name --image --replicas \
      --field="status.conditions[type=Completed].status" \
      --field="status.conditions[type=Complete].status" \
      --field="status.conditions[type=Ready].status" \
      --field="status.conditions[type=ContainersReady].status"

###

[tutorial-script]: ../gifs/cfg-tree.sh`
