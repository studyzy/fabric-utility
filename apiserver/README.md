
# 启动apiserver


##编译apiserver

apiserver目录：

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --mod=vendor

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo --mod=vendor


## 编译镜像

apiserver/目录：

docker build -t apiserver:v0.6 .


## 启动容器

demofabric目录：

  docker docker run -d --name apiserver-${orgname} -p 100${i}:1001 \ 
   -v ${prj_dir}/artifacts:/opt/apiserver/artifacts -v ${prj_dir}/config:/opt/apiserver/config
   apiserver:v0.6 apiserver -orgName ${orgname} -uerName Admin







