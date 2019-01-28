# AmrToMp3

AmrToMp3基本功能是将amr或wav格式的音频文件通过转换成为音频通用格式，程序使用goroutine并行处理增大
程序处理性能。环境依赖linux系统调用ffmpeg，如果没有对应运行环境可以下载docker镜像直接运行。
镜像地址为：docker.io/liyinda/alpine_amrtomp3


## 目录
* [环境](#环境)
* [下载](#下载)
* [编译](#编译)
  * [build binary](#build-binary)
  * [build docker image](#build-docker-image)
* [运行](#运行)
  * [run binary](#run-binary)
  * [run docker image](#run-docker-image)
* [运行参数](#运行参数)



## 环境

* [ffmpeg version 2.8.15](http://ffmpeg.org/)
* [docker-ce version 18.09.0](http://www.docker.com/)


## 下载

Binary can be downloaded from [Releases](https://github.com/liyinda/AmrToMp3/releases) page.

## 编译

### build binary

``` shell
go get  github.com/liyinda/AmrToMp3
go build main.go
```
### build docker image
``` shell
#go version 1.9.7 & docker-ce 18.09.0
docker build --network="bridge"  -t "docker.io/liyinda/amrtomp3" .
docker pull docker.io/liyinda/alpine_amrtomp3 
```

## 运行
``` shell
1）目录含义
./work #音频转换工作目录
./bak #音频文件备份目录
./audio #音频文件转换后生产目录

2）运行
nohup go run main.go &

```
### run docker
```
docker run docker.io/liyinda/alpine_amrtomp3
```

## 运行参数

``` shell
./main -h
Usage of ./main:
  -log string
        Log File Name (default "audio.log")
  -path string
        Audio Conversion Path (default ".")

```
