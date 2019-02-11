/*
Copyright 2018 The AmrToMp3 Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/


package main
 
import (
  "os"
  "strings"
  "time"
  "fmt"
  "path"
  "runtime"
  "flag"
  "log"
  models "github.com/liyinda/AmrToMp3/pkg/models"
)

var (
  filePath = flag.String("path", ".", "Audio Conversion Path")
  logFileName = flag.String("log", "audio.log", "Log File Name")
)

func putFile(fileChan chan string, WORKDIR string) {

  //将目录中的文件取出来放入管道fileChan 
  filelist := models.GetDirFiles(WORKDIR)
  if filelist == nil {
    //return
    fmt.Println("return nil")
  }
  for _, item := range filelist {
    fileChan<- item.Name()
  }
  close(fileChan)
} 

func audioConvert(fileChan chan string, amrChan chan string, exitChan chan bool) {
  //请求转换的工作目录
  WORKDIR := *filePath + "/work/"
  //备份目录
  BAKDIR := *filePath + "/bak/"
  //产品目录
  PRODUCTDIR := *filePath + "/audio/"
  //获取当前日期
  TODAY := time.Now().Format("20060102")

  for {
    file, ok := <-fileChan
    if !ok {
      break
    }
    //获取文件的内容
    wfilepath := WORKDIR + file
    bfilepath := BAKDIR + TODAY + "/" + file
    pfilepath := PRODUCTDIR + TODAY + "/" + file

    var rs bool
    var fileExt string
    //获取文件名后缀
    fileExt = path.Ext(file) 
    newfilepath := ""
    audiofilepath := ""

    switch {
      case fileExt == ".amr" :
        //进行amr到mp3的转换
        newfilepath = strings.Replace(wfilepath, ".amr", ".mp3", -1)
        comeAndTo := [...]string{wfilepath, newfilepath}
        rs = models.Arm2mp3(comeAndTo)
      case fileExt == ".wav" :
        //进行wav到mp3的转换
        newfilepath = strings.Replace(wfilepath, ".wav", ".mp3", -1)
        comeAndTo := [...]string{wfilepath, newfilepath}
        rs = models.Arm2mp3(comeAndTo)
      default:
        //os.Remove(wfilepath)
        os.Rename(wfilepath, bfilepath)
    }
    if rs {
      //转换成功之后，将原始文件移动到备份目录
      os.Rename(wfilepath, bfilepath)
      //生产mp3文件路径
      audiofilepath = strings.Replace(pfilepath, ".amr", ".mp3", -1)
      //转换成功之后，将生成文件移动到生产目录
      os.Rename(newfilepath, audiofilepath)
      //将转换成功的值传入amrChan管道
      amrChan<- audiofilepath
    }
  }
  //amrChan完成后将true传入exitChan管道
  fmt.Println("amrChan done")
  exitChan<- true

}

func main() {
  //调用运行参数
  flag.Parse()

  //请求转换的工作目录
  WORKDIR := *filePath + "/work/"
  //备份目录
  BAKDIR := *filePath + "/bak/"
  //产品目录
  PRODUCTDIR := *filePath + "/audio/"

  for {
    //获取cpu数量
    CPUNUM := runtime.NumCPU()
    //获取当前日期
    TODAY := time.Now().Format("20060102")
    //初始化判断是否存在日期目录，如果没有生成对应生产目录
    todaydir := PRODUCTDIR + TODAY
    models.MkDir(todaydir)
    //初始化判断是否存在日期目录，如果没有生成对应备份目录 
    todaybakdir := BAKDIR + TODAY
    models.MkDir(todaybakdir)
    //开启三个协程
    fileChan := make(chan string, 1000)
    amrChan := make(chan string, 2000)
    exitChan := make(chan bool, CPUNUM)
  
    //将目录中文件放入管道fileChan中
    go putFile(fileChan, WORKDIR)
  
    //开启cpu核数个协程，将fileChan中的文件名称取出放入amrChan中做转化，
    //amrChan协程完成后放入exitChan安全退出
    for i := 1; i <= CPUNUM; i++ {
      go audioConvert(fileChan, amrChan, exitChan)
    }
  
    go func(){
      for i := 1; i <= CPUNUM; i++ {
        <-exitChan
      }
      close(amrChan)
    }()
  
    for {
      res, ok := <-amrChan
      if !ok {
        break
      }
      //操作日志记录
      logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
      if logErr != nil {
        fmt.Println("Fail to find", *logFile, "AmrToMp3 start Failed")
        os.Exit(1)
      }
      defer logFile.Close()
      log.SetOutput(logFile)
      log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
      //写入日志
      log.Printf("%v Amr format conversion success \n", res)
      fmt.Printf("%v Amr format conversion success \n", res)
      
    }
    time.Sleep(1e9)
  }
}
 
