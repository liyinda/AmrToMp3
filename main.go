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
 
func main() {
  flag.Parse()
  for {
    arm2mp3_control()
    time.Sleep(1e9)
  }
}
 
//amr转mp3控制函数
func arm2mp3_control() {
  //请求转换的工作目录
  WORKDIR := *filePath + "/work/"
  //备份目录
  BAKDIR := *filePath + "/bak/"
  //产品目录
  PRODUCTDIR := *filePath + "/audio/"

  //fmt.Println(PRODUCTDIR)

  //获取cpu个数
  cpuNum := runtime.NumCPU()
  //读取目录下的文件信息
  filelist := models.GetDirFiles(WORKDIR)
  if filelist == nil {
    return
  }
  //fmt.Println("filelist =", filelist)

  //获取当前日期
  today := time.Now().Format("20060102")
  //初始化判断是否存在日期目录，如果没有生成对应生产目录
  todaydir := PRODUCTDIR + today
  models.MkDir(todaydir)
  //初始化判断是否存在日期目录，如果没有生成对应备份目录 
  todaybakdir := BAKDIR + today 
  models.MkDir(todaybakdir)

  //将filelist的文件列表传入ch管道 
  ch :=make(chan string)

  //创建goroutine获取文件列表并将列表传递给管道另一端进行音频转换
  go func(mine []os.FileInfo) {
    for _, item := range mine {
      ch <- item.Name()
    }
  }(filelist)

  go func() {
    for i :=0; i<= cpuNum; i++ {
      file := <-ch

      //获取文件的内容
      wfilepath := WORKDIR + file
      bfilepath := BAKDIR + today + "/" + file
      pfilepath := PRODUCTDIR + today + "/" + file
  
      var rs bool
      var fileExt string
      //获取文件名后缀
      fileExt = path.Ext(file) 
      newfilepath := ""
      audiofilepath := ""

      //操作日志记录
      logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
      if logErr != nil {
        fmt.Println("Fail to find", *logFile, "AmrToMp3 start Failed")
        os.Exit(1)
      }
      defer logFile.Close()
      log.SetOutput(logFile)
      log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
  
      switch {
        case fileExt == ".amr" :
          //进行amr到mp3的转换
          newfilepath = strings.Replace(wfilepath, ".amr", ".mp3", -1)
          comeAndTo := [...]string{wfilepath, newfilepath}
          rs = models.Arm2mp3(comeAndTo)
          log.Printf("%v Amr format conversion success \n", newfilepath)
        case fileExt == ".wav" :
          //进行wav到mp3的转换
          newfilepath = strings.Replace(wfilepath, ".wav", ".mp3", -1)
          comeAndTo := [...]string{wfilepath, newfilepath}
          rs = models.Arm2mp3(comeAndTo)
          log.Printf("%v Wav format conversion success \n", newfilepath)
        default:
          //os.Remove(wfilepath)
          os.Rename(wfilepath, bfilepath)
          log.Printf("%v conversion fail and backup \n", bfilepath)
      }
  
      if rs {
        //转换成功之后，将原始文件移动到备份目录
        os.Rename(wfilepath, bfilepath)
        //生产mp3文件路径
        audiofilepath = strings.Replace(pfilepath, ".amr", ".mp3", -1)
        //转换成功之后，将生成文件移动到生产目录
        os.Rename(newfilepath, audiofilepath)
      }
    }
  }()

}
 
