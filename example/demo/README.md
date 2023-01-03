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

## Recipies
All recipies are located in the `/recipies` directory.

### Rearrange Class Memebers

Run Tests:
```bash
chast test refactoring /recipes/rearrange_class_members/rearrange_class_members.chast.yml
```


## Build

Build:
```bash
go build ../../cli/main.go
```

Copy the newest built binary from the cli directory into the demo folder.
```bash
cp ../cli/main .
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
