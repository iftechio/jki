package config

const defaultConfig = `default-registry: ali
registries:
- name: ali
  aliyun:
    # 如果使用 access key 的话这里就不用设置
    # 这里的用户名、密码请访问 https://cr.console.aliyun.com/cn-hangzhou/instances/credentials 获取
    #username: user
    #password: passwd

    # 使用子账号登录后访问 https://usercenter.console.aliyun.com/#/manage/ak 创建自己的 access key
    # 如果出现 user not exist 错误, 需要到 https://cr.console.aliyun.com 开通服务
    access_key: foo
    secret_access_key: bar

    # 这个 namespace 指的是容器镜像服务里的命名空间
    # 可以到 https://cr.console.aliyun.com/cn-hangzhou/instances/namespaces 查看
    namespace: test

    region: cn-hangzhou
#- name: ali-ee
#  aliyun_ee:
#    #username: user
#    #password: passwd
#
#    access_key: foo
#    secret_access_key: bar
#    
#    namespace: test
#    region: cn-hangzhou
#	
#    # 企业实例id
#    instance_id: cri-123456
#
#    # optional 企业版可自己配host
#    # 如果不填将从 aliyun 的 api 取host, 需要提供access_key, secret_access_key
#    #instance_host: abc-registry.cn-hangzhou.cr.aliyuncs.com

#- name: aws-tokyo
#  aws:
#    access_key: foo
#    secret_access_key: bar
#    region: ap-northeast-1
#    account_id: "12345"
#- name: aws-bj
#  aws:
#    access_key: foo
#    secret_access_key: bar
#    region: cn-north-1
#    account_id: "45678"
#- name: mine
#  dockerhub:
#    username: foo
#    password: bar
`
