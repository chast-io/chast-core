# CHAST Demo

[![GitHub](https://img.shields.io/badge/DockerHub-0db7ed?logo=Docker\&logoColor=white)](https://hub.docker.com/r/rjenni/chast-demo)

<!-- TOC -->
* [CHAST Demo](#chast-demo)
  * [Run](#run)
  * [Recipies](#recipies)
    * [Rearrange Class Memebers](#rearrange-class-memebers)
  * [Build](#build)
    * [Push](#push)
<!-- TOC -->

## Run
```bash
docker run --rm -it --privileged rjenni/chast-demo:latest
```
* `--privileged` is required for the `unionfs-fuse` tool to work.
* `--rm` is recommended to remove the container after it has been stopped.
* `-it` is required to run the container in interactive mode.

*Note: Linux Kernel > v4.18.0 is required for the `unionfs-fuse` tool to work.*

* Exit the container with `exit`.

## Recipies
All recipies are located in the `/recipies` directory.
List them with:
```bash
ls -l /recipes
```

For checking the files as a tree under a specific folder, use:
```bash
tree <folder>
```

### Rearrange Class Memebers

Run Tests:
```bash
chast test refactoring /recipes/class_to_record/class_to_record.chast.yml
```

### Python 2 to Pyhton 3
> This is only a very limited demo using the `comby` tool.

Run Tests:
```bash
chast test refactoring /recipes/python2to3/python2to3.chast.yml
```

### Rearrange Class Memebers

Run Tests:
```bash
chast test refactoring /recipes/rearrange_class_members/rearrange_class_members.chast.yml
```

### Remove Double Negation

Run Tests:
```bash
chast test refactoring /recipes/remove_double_negation/remove_double_negation.chast.yml
```

## Build

Build:
```bash
# Run in `cli/` directory:
go build -o chast main.go
```

Copy the newest built binary from the cli directory into the demo folder.
```bash
cp ../../cli/chast .
```

Build the image:
```bash
 docker build -t rjenni/chast-demo . 
```

### Push
Push the image to the docker hub:
```bash
docker push rjenni/chast-demo
```
