# UC

### 用户中心

## Feature
+ [x] 邮箱注册/登录
+ [x] 手机号注册/登录
+ [ ] 发送短信验证码
+ [x] facebook 登录
+ [x] google 登录

## Install
#### 加载依赖
```go
go get -u github.com/jiebutech/uc
```

#### 初始化数据表结构
```sql
source ./schema.sql
```

## 自定义相关功能需要实现的接口
+ [x] 用户表 => github.com/jiebutech/uc/model.UserEntity
+ [x] 第三方登录表 => github.com/jiebutech/uc/model.OauthUserEntity
+ [x] 用户相关操作 => github.com/jiebutech/uc/model.UserResource
+ [x] 验证码发送器 => github.com/jiebutech/uc/sender.SmsCodeSender
+ [x] 邮箱验证码标题及内容模板 => github.com/jiebutech/uc/sender.CodeMessenger
+ [x] 验证码缓存器 => github.com/jiebutech/uc/cache.Cache

## 扩展第三方登录实现
1. 实现github.com/jiebutech/uc/oauth.Oauth接口, 完成登录注册的基本业务流程
2. 实现github.com/jiebutech/uc/oauth.OauthLoginType接口, 完成登录类型的相关定义

## Example

[用户中心接口实现](https://github.com/jiebutech/jin/tree/example/uc)
