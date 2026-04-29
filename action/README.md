# skill-guard GitHub Action

在 CI/CD 流水线中自动扫描 AI 技能文件的安全风险。

## 使用方法

```yaml
jobs:
  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: wjames2000/skill-guard@v0.2.0
        with:
          path: ./skills
          severity: high
          format: sarif
```

## 参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `path` | `.` | 扫描路径 |
| `severity` | `high` | 最低严重级别 |
| `format` | `sarif` | 输出格式 (terminal/json/sarif) |
| `args` | `""` | 额外参数 |
| `version` | `latest` | skill-guard 版本 |

## 输出

| 输出 | 说明 |
|------|------|
| `exit_code` | 0=无风险, 1=发现风险 |
| `issues_count` | 发现的问题数 |
