# spider

## python
打包已安装的依赖包

    pip freeze >requirements.txt
安装依赖包

    pip install -r requirements.txt  

### oss
对于Windows和Mac OS X系统，由于安装Python的时候会将Python依赖的头文件一并安装，因此您无需安装python-devel。
对于CentOS、RHEL、Fedora系统，请执行以下命令安装python-devel：
    
    yum install python-devel
对于Debian，Ubuntu系统，请执行以下命令安装python-devel：

    apt-get install python-dev
安装oss sdk
    pip install oss2
    
----


## go 

### infrastructure
```
go version go1.12.6 darwin/amd64
ide goland
```

### add package example (proxy may error cannot find package "golang.org/x/net/html")
```
dep ensure -add golang.org/x/net/html
```

### install package
```
dep ensure
```
### run 
```
go run go/src/main.go
```
