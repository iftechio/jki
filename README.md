# jki - JiKe Image utils

## 1. 安装

请到 https://github.com/bario/jki/releases 页面下载

下载后赋予可执行权限, 然后复制到 `PATH` 中, 比如:
```
$ chmod +x jki_0.0.3_darwin_amd64
# cp jki_0.0.3_darwin_amd64 /usr/local/bin/jki
```

## 2. 使用方法

### 2.1. 生成配置

#### 2.1.1 保存默认配置

```
$ jki config init --save
```

#### 2.1.2 修改配置

```
$ jki config edit
```

#### 2.1.3 查看配置

```
$ jki config view
default-registry: ali
registries:
- name: ali
  aliyun:
    username: <USERNAME>
    password: <PASSWORD>
    region: cn-hangzhou
    namespace: <YOUR NAMESPACE>
- name: aws-tokyo
  aws:
    access_key: <YOUR ACCESS KEY>
    secret_access_key: <YOUR SECRET ACCESS KEY>
    region: ap-northeast-1
    account_id: <YOUR ACCOUNT ID> # 注意填写的 account id 两边要有双引号
```

#### 2.1.4 检查配置正确性

```
$ jki config check
```
如果配置语法没问题的话会打印 `OK!`

### 2.2. 构建镜像

默认情况下会把构建出来的镜像推送到配置里指定的 `default-registry`
```
$ jki build
```

指定 Dockerfile 路径:
```
$ jki build -f <YOUR Dockerfile>
```

指定构建目录:
```
$ jki build <YOUR DIR>
```

指定镜像的名字:
```
$ jki build --image-name <IMAGE NAME>
```

指定要推送的 registry (以上面的配置为例):
```
$ jki build --registry aws-tokyo
```

### 2.3. 跨云服务商复制镜像

以上面的配置为例:
```
$ jki cp <YOUR ACCOUNT ID>.dkr.ecr.ap-northeast-1.amazonaws.com/foo:bar
```
会把 `<YOUR ACCOUNT ID>.dkr.ecr.ap-northeast-1.amazonaws.com/foo:bar` 该镜像复制到 `ali` 对应的 registry 上


指定目标 registry:
```
$ jki cp k8s.gcr.io/etcd:3.3.10 aws-tokyo
```
会把 `k8s.gcr.io/etcd:3.3.10 ` 该镜像复制到 `aws-tokyo` 对应的 registry 上
