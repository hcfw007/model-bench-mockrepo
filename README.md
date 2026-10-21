# model-bench-mockrepo

**中文** | [English](README.en.md)

[model-bench](https://github.com/hcfw007/model-bench) 的测试床。

一个脱敏的多服务仓库，用来复现真实工作里的任务形状：跨服务定位、最小编辑、
证据取证、SOP 执行、subagent 协作。它不是任何真实业务的代码。

## 为什么单独一个仓

主仓跑分时用 `git worktree` 给每个 case 做隔离，所以测试床需要自己的 git 历史。
其中一个 case 要求「只靠 git 历史找出引入问题的那个 commit」，
因此这里刻意保留了带缺陷的提交链。

## 结构

```
go-projects/
  some-service/          主服务：任务创建、确认、轮询、打标
    internal/biz/bulk_send_v1/   Kratos 风格三层里的业务层
  auth-service/          权限组、成员组、NEW 角标
  contact-data-service/  客户档案与标签
ts-projects/
  api-bff/               BFF：素材校验、确认项投影
  web-console/           前端页面
.tasks/bulk-send-v1/     跨会话任务档案（STATUS / record / 用例树）
fixtures/                脱敏日志、mongo 文档、卡死脚本
```

## 预埋了什么

每个服务里都有若干个与真实缺陷同构的问题，分布在：

- 收尾条件：部分对象失败时任务能否正常收敛
- 自动打标：配置落库后有没有真正的执行点
- 确认接口：已取消任务返回的错误是否可区分
- 前端列表：是否消费后端已经算好的可确认性字段
- 权限与角标：注册表缺项、分页截断
- 任务档案：两处记录互相矛盾（用来测「是否会发现并停下来问」）

具体判据在主仓的 `cases/` 里，这里不展开。

## 验证

```bash
go build ./...
go test ./... -count=1
```

Go 侧零外部依赖，标准库即可。

## 约束

- 不配置任何真实远端，不允许 `git push`
- 不调用任何真实消息或群发接口
- 所有 ID、日志、客户名均为构造数据
