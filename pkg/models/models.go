package models
 
import (
    "os/exec"
    "os"
    "fmt"
    "io/ioutil"
)
 
//执行命令函数: 不能感知命令的执行信息，只返回是否执行成功
//commandName 命名名称，如cat，ls，git等
//params 命令参数，如ls -l的-l，git log 的log等
func execCommand(commandName string, params []string) bool {
    cmd := exec.Command(commandName, params...)
 
    err := cmd.Start()
    if err != nil {
        return false
    }
 
    err = cmd.Wait()
    if err != nil {
        return false
    }
 
    return true
}


//如果文件夹不存在，创建文件夹
func MkDir(path string) {
  exist, err := pathExists(path)
  if err != nil {
    fmt.Printf("get dir error![%v]\n", err)
    return
  }
  if exist {
    fmt.Printf("has dir![%v]\n", path)
    //continue
  } else {
    err := os.Mkdir(path, os.ModePerm)
    if err != nil {
      fmt.Printf("mkdir failed![%v]\n", err)
    } else {
      fmt.Printf("mkdir success![%v]\n", path)
    }
  }
}


//判断文件夹是否存在
func pathExists(path string) (bool, error) {
  _, err := os.Stat(path)
  if err == nil {
      return true, nil
  }
  if os.IsNotExist(err) {
      return false, nil
  }
  return false, err
}


//遍历目录下的文件
//path 遍历的目录
//返回文件记录的切片
//cpunum等于系统核心数   
func GetDirFiles(path string) []os.FileInfo {
  //最大处理文件数为100
  var filelist [100]os.FileInfo
  var i int = 0

  dir, err := ioutil.ReadDir(path)
  if err != nil {
    return nil
  }

  for _, file := range dir {
    if file.IsDir() {
      continue
    }

    if i >= 100 {
      break
    }

    filelist[i] = file
    i++
  }

  if i == 0 {
    return nil
  }

  return filelist[0:i]
}

//进行amr到mp3音频格式的转换
//farr 文件信息：0arm文件名；1mp3文件名
//返回是否成功：true成功，false失败
func Arm2mp3(farr [2]string) bool {
  //执行命令ffmpeg -i 1.amr 1.mp3
  params := [...]string{"-i", farr[0], farr[1]}
  //fmt.Print(params[0:3])
  rs := execCommand("ffmpeg", params[0:3])

  return rs
}
