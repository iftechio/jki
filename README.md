# jki - JiKe Image utils
[![Build Status](https://travis-ci.org/iftechio/jki.svg?branch=master)](https://travis-ci.org/iftechio/jki)

## 0. Features

1. 构建 Docker image（当 Docker daemon 支持 buildkit 的时候自动启用 buildkit）并方便地 push 到不同的云厂商 registry 上
2. 在不同的云厂商 registry 间复制 image
3. 自动替换修复 Deployment 不能访问的镜像

## 1. 安装最新版本

```

curl -s https://api.github.com/repos/iftechio/jki/releases/latest \
  | grep browser_download_url \
  | grep darwin \
  | cut -d '"' -f 4 \
  | wget -i -

tar --strip-components=1 -xf jki*

sudo mv jki /usr/local/bin/
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
    # 如果使用 access key 的话这里就不用设置
    # 用户名、密码请用子账号访问 https://cr.console.aliyun.com/cn-hangzhou/instances/credentials 获取
    #username: <YOUR USERNAME>
    #password: <YOUR PASSWORD>

    # 使用子账号登录后访问 https://usercenter.console.aliyun.com/#/manage/ak 创建自己的 access key
    # 如果出现 user not exist 错误, 需要到 https://cr.console.aliyun.com 开通服务
    access_key: <YOUR ACCESS KEY ID>
    secret_access_key: <YOUR ACCESS KEY SECRET>

    # 这个 namespace 指的是容器镜像服务里的命名空间
    # 可以到 https://cr.console.aliyun.com/cn-hangzhou/instances/namespaces 查看
    namespace: <REGISTRY NAMESPACE>

    region: cn-hangzhou
- name: aws-tokyo
  aws:
    access_key: <YOUR ACCESS KEY>
    secret_access_key: <YOUR SECRET ACCESS KEY>
    region: ap-northeast-1
    account_id: <YOUR ACCOUNT ID> # 注意填写的 account id 两边要有双引号
- name: mine
  dockerhub:
    username: foo
    password: bar
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

更多选项可以参考 `jki build -h`

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

会把 `k8s.gcr.io/etcd:3.3.10` 该镜像复制到 `aws-tokyo` 对应的 registry 上

自动复制最新的 tag (仅限于 AWS ECR 跟阿里云容器镜像服务):
```
$ jki cp <YOUR ACCOUNT ID>.dkr.ecr.ap-northeast-1.amazonaws.com/foo
```

会查询 `<YOUR ACCOUNT ID>.dkr.ecr.ap-northeast-1.amazonaws.com/foo` 该 image 最新的 tag, 然后复制到 `ali` 对应的 registry

### 2.4. 自动替换修复 deployment 不能访问的镜像

执行命令后会逐个提示替换无法现在的镜像

```
$ jki transferimage --namespace default
Transfer Deployment/foo1 gcr.io/foo1:bar(y/n)?
>y
Transfer gcr.io/foo1:bar to xxx.dkr.ecr.ap-northeast-1.amazonaws.com/foo1:bar
```
