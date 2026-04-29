# 🛡️ skill-guard

**AI 技能的零信任安全扫描器**  
*在使用第三方技能之前先行扫描 — 快速、离线、无依赖。*

[![version](https://img.shields.io/badge/版本-v0.2.0-blue)](https://github.com/wjames2000/skill-guard/releases)
[![go version](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)](https://golang.org)
[![license](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![tests](https://img.shields.io/badge/测试-71%20通过-brightgreen)](https://github.com/wjames2000/skill-guard/actions)
[![gosec](https://img.shields.io/badge/gosec-0%20issues-brightgreen)](https://github.com/wjames2000/skill-guard)
[![platform](https://img.shields.io/badge/平台-linux%20|%20macOS%20|%20windows-lightgrey)](https://github.com/wjames2000/skill-guard/releases)

[English](README.md) | 中文

---

skill-guard 检测第三方 AI 技能文件（`.md`、`.json`、`.yaml`、`.py`、`.sh` 等）中的安全威胁——密钥泄露、恶意命令、混淆载荷等。

- 🔒 **零信任** — 默认不信任任何第三方技能，扫描后方可使用
- ⚡ **零依赖** — 单一原生二进制文件，无需任何运行时环境
- 📦 **100 条内置规则** — 覆盖 12 大安全类别
- 🤖 **AI 辅助检测** — 可选 LLM 验证，降低误报
- 🚀 **CI/CD 友好** — JSON/SARIF 输出、语义化退出码、GitHub Action

---

## 🚀 安装

### macOS（Homebrew）
```bash
brew install wjames2000/tap/skill-guard
```

### Linux / macOS（一键安装）
```bash
curl -sfL https://github.com/wjames2000/skill-guard/releases/latest/download/install.sh | sh
```

### Windows（Scoop）
```bash
scoop bucket add wjames2000 https://github.com/wjames2000/scoop-bucket.git
scoop install skill-guard
```

### Go 安装
```bash
go install github.com/wjames2000/skill-guard@latest
```

> 预编译的 Linux (amd64/arm64)、macOS (amd64/arm64)、Windows (amd64) 二进制文件可在 [Releases 页面](https://github.com/wjames2000/skill-guard/releases) 下载。

---

## ⚙️ 快速上手

```bash
# 扫描技能目录
skill-guard ./downloaded-skills

# 启用 AI 辅助检测
skill-guard ./skills --ai

# 输出 SARIF 格式（GitHub Code Scanning 兼容）
skill-guard ./skills --sarif

# 导出全部 100 条内置规则
skill-guard rules export yaml > my-rules-backup.yaml

# 一键初始化 pre-commit hook
cd my-project && skill-guard init
```

---

## 📋 核心功能

### 🔍 智能目录扫描
递归扫描文本文件，自动跳过 `.git`、`node_modules`、`vendor` 等目录。支持 `.gitignore` 规则和符号链接跳过。

```bash
skill-guard                              # 扫描当前目录
skill-guard ./path/to/skills             # 扫描指定目录
skill-guard ./a ./b                      # 同时扫描多个目录
skill-guard --gitignore                  # 遵循 .gitignore 规则
```

### 🧠 100 条内置规则（12 大类别）

| 类别 | 数量 | 检测目标 |
|------|------|---------|
| **密钥泄露** | 17 | AWS Key、GitHub Token、JWT、Stripe、GCP、Azure、Discord、npm、Slack |
| **命令执行** | 7 | `os.system`、`subprocess`、`child_process`、`eval`、`exec` |
| **代码注入** | 8 | SQLi、NoSQLi、XSS、SSTI、PHP 代码执行、shell=True |
| **恶意文件操作** | 6 | `rm -rf`、任意写入、chmod 777、系统文件覆盖 |
| **网络请求滥用** | 6 | `curl \| bash`、反向 Shell、wget pipe、SSRF |
| **信息窃取与混淆** | 6 | Base64、Hex、SSH 密钥读取、环境变量泄露 |
| **配置风险** | 8 | Debug 模式、CORS `*`、弱 TLS、CSRF 禁用、不安全 Cookie |
| **加密与认证** | 8 | MD5/SHA1、ECB 模式、JWT 密钥、硬编码证书、弱 RSA |
| **供应链风险** | 7 | pip 直装 GitHub、通配符依赖、npm postinstall |
| **信息泄露** | 8 | 堆栈暴露、调试端点、硬编码 IP、SMTP 凭据 |
| **容器风险** | 8 | 特权容器、Docker socket、host 网络、USER root |
| **后门与 C2** | 10 | C2 回调、crontab 写入、pickle 反序列化、yaml.load |

### 🤖 AI 辅助检测
通过本地 LLM 对匹配结果进行语义验证，显著降低误报率：

```bash
skill-guard ./skills --ai                              # 默认: gemma-4-26b-a4b-it
skill-guard ./skills --ai --ai-model llama3.2           # 自定义模型
skill-guard ./skills --ai --ai-endpoint https://my-api/v1  # 自定义端点
```

### ✏️ 自定义规则与规则市场
```bash
# 加载自定义规则
skill-guard ./skills --rules my-rules.yaml

# 创建规则模板
skill-guard rules new

# 测试规则匹配效果
skill-guard rules test my-rule.yaml "dangerous_code()"

# 浏览规则市场
skill-guard update info
skill-guard update list
skill-guard update install <url>
```

### 📤 多种输出模式
```bash
skill-guard ./skills                              # 终端彩色输出（默认）
skill-guard ./skills --json                       # 结构化 JSON
skill-guard ./skills --quiet                      # 仅显示有问题的文件
skill-guard ./skills --summary                    # 仅统计数据
skill-guard ./skills --sarif                      # SARIF 2.1 格式
skill-guard ./skills --no-color                   # 无颜色纯文本
skill-guard ./skills --output report.json         # 写入 JSON 文件
skill-guard ./skills --severity high              # 严重级别过滤
```

### 🔧 配置文件
```yaml
# .skillguard.yaml
severity: "high"
ext-include: [".py", ".sh", ".yaml"]
ignore: ["tests/**"]
ai_enabled: true
ai_model: "gemma-4-26b-a4b-it"
```

### 🔔 Pre-commit Hook
```bash
skill-guard init          # 一键安装 git pre-commit hook
```

或集成 pre-commit 框架：
```yaml
repos:
  - repo: https://github.com/wjames2000/skill-guard
    rev: v0.2.0
    hooks:
      - id: skill-guard
        args: ["--severity", "high", "--quiet"]
```

### 🐚 Shell 自动补全
```bash
skill-guard completion bash   # bash 补全
skill-guard completion zsh    # zsh 补全
skill-guard completion fish   # fish 补全
```

---

## 🧪 CI/CD 集成

### GitHub Actions
```yaml
- uses: wjames2000/skill-guard@v0.2.0
  with:
    path: .
    severity: high
    format: sarif
```

### 独立使用
```yaml
- name: 安全扫描
  run: |
    skill-guard ./skills --json --severity high
```

扫描发现风险时退出码为 `1`，可自然阻断有风险的 PR。

---

## 📊 性能指标

| 场景 | 耗时 | 内存 |
|------|------|------|
| 100 规则 × 1 文件 | ~0.14 ms | 68 KB |
| 100 文件扫描 | ~7.9 ms | 7 MB |
| 500 文件扫描 | ~44 ms | 35 MB |
| 50 MB 二进制文件 | ~0.045 s（自动跳过） | — |

---

## 🔐 安全声明

- **不执行被扫描文件** — 只读分析
- **完全离线** — 扫描过程无网络请求（`update` 命令除外）
- **不收集数据** — 不会发送任何使用数据或扫描结果
- **gosec: 0 issues** — 静态安全分析通过

---

## 🗺️ 路线图

| 版本 | 重点 |
|------|------|
| v0.1.0 | CLI、目录遍历、32 条规则、并发扫描、退出码 |
| v1.1 | 质量加固：Scanner 测试、panic 保护、基准、Windows 兼容、symlink、进度条、.gitignore |
| v1.2 | 体验增强：Shell 补全、--no-color、--output、--summary、SARIF |
| v2.0-α | 生态扩展：规则市场、AI 检测、skill-guard init |
| v2.0-β | 100 条规则（12 类别）、71 测试、gosec 0 |
| v2.1 | 质量深化：66 条新规则测试、补充测试、压力基准 |
| v2.2 | 规则生态：rules export/new/test、版本管理、贡献指南 |
| v2.3 | 平台集成：GitHub Action、pre-commit hooks、Homebrew |
| v2.4 | 高级引擎：LUA 脚本、YARA、语义分析 |
| v3.0 | 平台化：规则市场 Web、团队管理、API 服务 |

---

## 🤝 参与贡献

- [规则贡献指南](docs/规则贡献指南.md)
- [CI/CD 集成指南](docs/ci-integration.md)
- [Homebrew 提交指南](docs/homebrew-submission.md)
- 提交 [Issue](https://github.com/wjames2000/skill-guard/issues) 或 PR

---

## 📜 许可证

MIT © 2026 [wjames2000](https://github.com/wjames2000)

---

*skill-guard — 扫描之后再信任。*
