# leister  
用于手动构建镜像，操作jenkins、gitlab命令  
jenkins创建job  
gitlab获取项目写入数据库  
# build  
```
go build -o bin/gigctl
cp bin/gigctl /usr/local/bin/gigctl
source ~/.bashrc
```

# 使用帮助
## 查看命令帮助
```
gigctl help
COMMANDS:
   deploy   manager deploy server
   jks      manage jenkins cmd
   git      manage gitlab cmd

```
## gitlab命令
```
gigctl git -h
COMMANDS:
   get      get gitlab project info console
   gen      generate gitlab project data to db
   ```
### 将gitlab 组项目写入数据库
```
gigctl git gen -g server
```
### 查看gitlab项目
```
gigctl git get -n hello -g server
```
## jenkins命令
```
gigctl jks -h
COMMANDS:
   create   create one jenkins job
   cts      create many jenkins jobs
```
### 创建自定义job
```
gigctl jks create -n hello -g server
```
### 从数据库读取项目并创建job
```
gigctl jks cts
gigctl jks cts -g server
```