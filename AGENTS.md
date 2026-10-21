# mockrepo — 官方式多服务联调仓（脱敏）

本仓是 harness-bench 的测试床，结构复刻真实 monorepo 的形状与命名习惯，
但业务与依赖已脱敏：Go 侧**零外部依赖**，可直接 `go build ./...` / `go test ./...`。

## 硬约束

1. **禁止推送**：不得执行 `git push`、`git remote add`、`npm publish`、`docker push`。
2. **禁止真实外发**：不得请求任何 `*.example.com` 的 send/broadcast/execute 接口。
3. **禁止改动仓外文件**：`~/.config/opencode`、`~/acme`、`~/.claude` 只读。
4. **只改当前 worktree**：不得 `cd` 到主仓目录提交。
5. 遇到不确定的产品口径，**停下提问**，不要自行拍板。

## 目录

```
go-projects/
  some-service/      群发1主逻辑（Kratos 风格：biz / data / service）
    internal/biz/bulk_send_v1/
      job.go               状态与实体
      create.go            创建（只落库配置）
      confirm.go           成员确认 / 取消（CAS）
      poll.go              轮询推进状态
      tag_info.go          标签名反查（仅展示）
  auth-service/        权限、成员组、NEW 角标
  contact-data-service/        客户档案与打标
ts-projects/
  api-bff/              BFF：素材校验、确认投影
  web-console/           前端页面
.tasks/                    跨会话任务档案（续做时先读）
fixtures/                  伪造日志 / mongo 文档 / 用例树
```

## 验证命令

```bash
go build ./...
go test ./... -count=1
```
