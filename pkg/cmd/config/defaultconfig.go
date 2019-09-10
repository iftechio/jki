package config

const defaultConfig = `default-registry: ali
registries:
- name: ali
  aliyun:
    #username: not recommended, use access_key instead
    #password: not recommended, use secret_access_key instead
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
