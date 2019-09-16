package config

const defaultConfig = `default-registry: ali
registries:
- name: ali
  aliyun:
    # 这里的用户名、密码请访问 https://cr.console.aliyun.com/cn-hangzhou/instances/credentials 获取
    #username: user
    #password: passwd

    # 使用子账号登录后访问 https://usercenter.console.aliyun.com/#/manage/ak 创建自己的 access key
    # 如果出现 user not exist 错误, 需要到 https://cr.console.aliyun.com 开通服务
    access_key: foo
    secret_access_key: bar
    region: cn-hangzhou
    namespace: test
- name: aws-tokyo
  aws:
    access_key: foo
    secret_access_key: bar
    region: ap-northeast-1
    account_id: "12345"
- name: aws-bj
  aws:
    access_key: foo
    secret_access_key: bar
    region: cn-north-1
    account_id: "45678"
`
