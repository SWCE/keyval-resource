
[![Docker Stars](https://img.shields.io/docker/stars/swce/keyval-resource.svg?style=plastic)](https://registry.hub.docker.com/v2/repositories/swce/keyval-resource/stars/count/)
[![Docker pulls](https://img.shields.io/docker/pulls/swce/keyval-resource.svg?style=plastic)](https://registry.hub.docker.com/v2/repositories/swce/keyval-resource)
[![Docker build status](https://img.shields.io/docker/build/swce/keyval-resource.svg)](https://github.com/swce/keyval-resource)
[![Docker Automated build](https://img.shields.io/docker/automated/swce/keyval-resource.svg)](https://github.com/swce/keyval-resource)

# Concourse CI Key Value Resource

Implements a resource that passes key values between jobs without using any external resource such as git/s3 etc.

## Thanks

This resource was implemented based on the [time resource](https://github.com/concourse/time-resource)

## Source Configuration

``` YAML
resource_types:
  - name: keyval
    type: docker-image
    source:
      repository: swce/keyval-resource
      
resources:
  - name: keyval
    type: keyval
```

#### Parameters

*None.*

## Behavior

### `check`: Produce a single dummy key

This is a version-less resource so `check` behavior is no-op

### `in`: Report the given time.

Fetches the given key values and stores them in the `keyval.properties` file.
The format is of a `.properties` file, e.g. `"<key>=<value>"`.

#### Parameters

*None.*

### `out`: Consumes the given properties file

``` YAML
- put: keyval
  params:
    file: keyvalout/keyval.properties
```

Reads the given properties file and sets them for the next job.

#### Parameters
- file - the properties file to read the key values from

## Examples

```YAML
resource_types:
  - name: keyval
    type: docker-image
    source:
      repository: swce/keyval-resource

resources:
  - name: keyval
    type: keyval

jobs:

  - name: build
    plan:
      - task: build
        file: tools/tasks/build/task.yml
      - put: keyval
        params:
          file: keyvalout/keyval.properties

  - name: test-deploy
    plan:
      - aggregate:
        - get: keyval
          passed:
          - build
      - task: test-deploy
        file: tools/tasks/task.yml
```

The build job writes all the key values it needs to pass along (e.g. artifact id) in the `keyvalout/keyval.properties` file. 
The test-deploy can read the data from the `keyval/keyval.properties` file and use them as it pleases. 

## CI suggestions

### Auto export the keys

You can add the following bash script to the **start** of every job to auto export the passed key values, if they exist. 
The script assumes that the resource folder is `keyval`. 

* Don't forget to source the script so it's exports will be passed along

```bash
#!/bin/bash

props="${ROOT_FOLDER}/keyval/keyval.properties"
if [ -f "$props" ]
then
  echo "Reading passed key values"
  while IFS= read -r var
  do
    if [ ! -z "$var" ]
    then
      echo "Adding: $var"
      export "$var"
    fi
  done < "$props"
fi

```

### Auto pass the keys

You can add the following bash script to the **end** of every job to auto pass specific environment variables as key values to the next job. 
The script only passes environment variables that start with `PASSED_`. 
The script assumes that the resource out file is `keyvalout/keyval.properties`:

e.g. 
```YAML
- put: keyval
  params:
    file: keyvalout/keyval.properties
``` 

```bash
#!/bin/bash

propsDir="${ROOT_FOLDER}/keyvalout"
propsFile="${propsDir}/keyval.properties"
if [ -d "$propsDir" ]
then
  touch "$propsFile"
  echo "Setting key values for next job in ${propsFile}"
  while IFS='=' read -r name value ; do
    if [[ $name == 'PASSED_'* ]]; then
      echo "Adding: ${name}=${value}"
      echo "${name}=${value}" >> "$propsFile"
    fi
  done < <(env)
fi

```

## Development

### Prerequisites

* golang is *required* - version 1.9.x is tested; earlier versions may also
  work.
* docker is *required* - version 17.06.x is tested; earlier versions may also
  work.
* godep is used for dependency management of the golang packages.

### Running the tests

The tests have been embedded with the `Dockerfile`; ensuring that the testing
environment is consistent across any `docker` enabled platform. When the docker
image builds, the test are run inside the docker container, on failure they
will stop the build.

Run the tests with the following command:

```sh
docker build -t keyval-resource .
```

### Contributing

Please make all pull requests to the `master` branch and ensure tests pass
locally.
