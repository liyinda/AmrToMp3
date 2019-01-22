# AmrToMp3

AmrToMp3基本功能是将amr或wav格式的音频文件通过转换成为音频通用格式，程序使用goroutine并行处理增大
程序处理性能。环境依赖linux系统调用ffmpeg，如果没有对应运行环境可以下载docker镜像直接运行。
镜像地址为：


## 目录
* [环境](#环境)
* [下载](#下载)
* [编译](#编译)
  * [build binary](#build-binary)
  * [build docker image](#build-docker-image)
* [运行](#运行)
  * [run binary](#run-binary)
  * [run docker image](#run-docker-image)
* [运行参数&google-authenticator手机端](#运行参数&google-authenticator手机端)
* [用户后台与google-authenticator对接](#户后台与google-authenticator对接)
  * [Redmine](#redmine)
  * [Zabbix](#zabbix)



## 环境

* [ffmpeg](http://ffmpeg.org/)
* [docker](http://www.docker.com/)


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
make docker
DOCKER 部署方式作者会尽快补充
docker pull 空:latest
```

## 运行
``` shell
1）确保redis服务运行正常，并确保redis存储是持久化化配置。
./redis-server /etc/redis/6379.conf

2）生成用户秘钥和google-authenticator二维码
mkdir jpg
./createGoogleCode [用户名]

3）运行认证接口服务端
./verificationGoogleCode

4）测试接口访问是否正常
curl "http://127.0.0.1:8082/get?issuser=[用户名]&code=[google验证码]"
如返回ok表示返回正常
如返回error表示返回异常
```
### run docker
```
DOCKER 部署方式作者会尽快补充
docker pull 空:latest
```

## 运行参数&google-authenticator手机端

### 可根据自身环境更改运行参数
``` shell
./verificationGoogleCode -h
Usage of ./verificationGoogleCode:
  -http.address string
        Address on HTTP Listen . (default ":8082")
  -log string
        Log file name (default "authenticator.log")
  -redis.address string
        Address on Redis Server . (default "127.0.0.1:6379")

```

### 手机下载google-authenticator客户端
iphone手机和android手机都有对应的客户端，请大家自行下载

![image](https://github.com/liyinda/google-authenticator/blob/master/jpg/google-authenticator.jpg)


## 用户后台与google-authenticator对接

### Redmine
vi app/views/account/login.html.erb
``` shell 
添加
14 <tr>
15     <td style="text-align:right;"><label for="code">Google验证码:</label></td>
16     <td style="text-align:left;"><%= text_field_tag 'code', nil, :tabindex => '3' %></td>
17 </tr>

```

vi app/controllers/account_controller.rb
``` shell 
添加
1   require "open-uri"

192   def password_authentication
193 
194     uri = 'http://[google-authenticator服务端地址]/get?issuser=' + params[:username] + '&code=' + params[:code]
195     html_response = nil
196     open(uri) do |http|
197     html_response = http.read
198     end
199 
200     if html_response == 'ok'
201 
202     user = User.try_to_login(params[:username], params[:password], false)
203     if user.nil?
204       invalid_credentials
205     elsif user.new_record?
206       onthefly_creation_failed(user, {:login => user.login, :auth_source_id => user.auth_source_id })
207     else
208       # Valid user
209       if user.active?
210         successful_authentication(user)
211       else
212         handle_inactive_user(user)
213       end
214     end
215 
216     else
217         redirect_to(:action => 'login')
218     end
219 
220 
221   end

```

![image](https://github.com/liyinda/google-authenticator/blob/master/jpg/redmine.jpg)

### Zabbix
vi include/views/general.login.php
``` shell 
添加
55         ->addItem([new CLabel(_('Password'), 'password'), (new CTextBox('password'))->setType('password')])
56         ->addItem([
57                 new CLabel(_('Google Code'), 'code'),
58                 (new CTextBox('code'))->setAttribute('', ''),
59                 $error
60         ])

```

vi index.php 
``` shell 
添加
65 if (isset($_REQUEST['enter']) && $_REQUEST['enter'] == _('Sign in')) {
66         // try to login
67         $autoLogin = getRequest('autologin', 0);
68         //print_r($_REQUEST);
69         $authflag=file_get_contents("http://[google-authenticator服务端地址]/get?issuser=".getRequest('name', '')."&code=".getRequest('code', ''));
70         //echo "http://[google-authenticator服务端地址]/get?issuser=".getRequest('name', '')."&code=".getRequest('code', '');
71         if ($authflag=='ok'){}else{
72             echo 'Google验证码错误'; header('Refresh: 2; url=http://zabbix.org/');exit;
73         }
74         //echo getRequest('code', '');

```


![image](https://github.com/liyinda/google-authenticator/blob/master/jpg/zabbix.jpg)


更多后台对接改造等您实现
