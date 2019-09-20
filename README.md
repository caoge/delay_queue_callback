# delay-queue

基于redis的有序队列实现的延迟队列, 参考[有赞延迟队列设计](http://tech.youzan.com/queuing_delay)实现，但是没有采用文中的短轮询方式依赖客户端拉取，而是采用服务端回调的方式对客户端进行通知，遵循RESTFUL协议实现HTTP接口交互

## 应用场景
* 商品定时上架
* 订单超时未支付取消
* 拼团时间结束 关闭团购



## 依赖
* Redis



## 源码安装
* `go`语言版本1.11+
* `git clone https://github.com/ROOKIE20570/delay_queue_callback.git`
* `cd /path/to/cmd`
* `go build cmd.go`



## 运行
`./cmd` 

`参数 -c 指定配置文件的路径` 

##使用
调用添加任务接口添加任务，时间到达执行时间后，向客户端发起http请求，方法为PUT 携带信息，客户端需返回 ok 表示通知成功  否则会一直请求直到最大重试次数

### 携带信息

```json
{
  "topic": "inform",
  "id": "08019413123",
  "delay": 1567135500,
  "callback":"https://www.baidu.com",
  "body": {"example": true}
}
```
|  参数名 |     类型    |     含义     | 
|:-------:|:-----------:|:------------:|
|   topic  | string     |      当前任务类型                 | 
|   id     | string     |    当前任务唯一标识                   |                   
|   delay  | int        |    延迟时间 单位:秒   到时间后，回调callback参数提供的URL地址  |                   
|   callback  | string        |    回调URL  |                
|   body   | string     |    任务的额外内容 |             

## HTTP接口

URL `/job`

### 添加任务 

请求方法 `POST`

```json
{
  "topic": "inform",
  "id": "08019413123",
  "delay": 1567135500,
  "callback":"https://www.baidu.com",
  "body": {"example": true}
}
```
|  参数名 |     类型    |     含义     |
|:-------:|:-----------:|:------------:|
|   topic  | string     |      当前任务类型                 |
|   id     | string     |    当前任务唯一标识                   |
|   delay  | int        |    延迟时间 单位:秒   到时间后，回调callback参数提供的URL地址  |
|   callback  | string        |    回调URL  |
|   body   | string     |    任务的额外内容 |

```json
{
  "success": 1,
  "message": "添加新任务成功",
  "data": null
}
```

### 获取任务信息

请求方法 `GET`

`query参数`

|  参数名 |     类型    |     含义     |
|:-------:|:-----------:|:------------:|
|   topic  | string     |      当前任务类型                 |
|   id     | string     |    当前任务唯一标识                   |  


```json
{
  "success": 1,
  "message": "添加新任务成功",
  "data": {
            "topic": "inform",
            "id": "08019413123",
            "delay": 1567135500,
            "body": {"example": true},
            "job_sign":"inform-08019413123"
          }
}


```
|  参数名 |     类型    |     含义     |
|:-------:|:-----------:|:------------:|
|   topic  | string     |      当前任务类型                 |
|   id     | string     |    当前任务唯一标识                   |
|   delay  | int        |    延迟时间 单位:秒    |
|   callback  | string        |    回调URL  |
|   body   | string     |    任务的额外内容 |
|   job_sign   | string     |    全局唯一标识  客户端可忽略 |






### 删除任务
请求方法 DELETE  

`query参数`

|  参数名 |     类型    |     含义     |
|:-------:|:-----------:|:------------:|
|   topic  | string     |      当前任务类型                 |
|   id     | string     |    当前任务唯一标识                |  

返回值
```json
{
    "success": 1,
    "message": "删除job成功",
    "data": null
}
```


  
