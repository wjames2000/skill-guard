# 🛡️ skill-guard

**AI 技能的零信任安全扫描器**  
*在使用第三方技能之前先行扫描 —— 快速、离线、无依赖。*

[![version](https://img.shields.io/badge/版本-1.0-blue)](https://github.com/wjames2000/skill-guard)
[![go version](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)](https://golang.org)
[![license](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![platform](https://img.shields.io/badge/平台-linux%20|%20macOS%20|%20windows-lightgrey)](https://github.com/wjames2000/skill-guard/releases)

---

**skill-guard** 是一款面向 AI 辅助开发场景的命令行安全扫描工具。  
它可以检测第三方技能文件（Markdown、JSON、YAML、Python 等）中隐藏的安全风险，包括密钥泄露、恶意命令、混淆载荷等。

- 🔒 **零信任** —— 默认不信任任何第三方技能，扫描后方可使用
- ⚡ **零依赖** —— 单一原生二进制文件，无需任何运行时环境
- 📦 **内置 30+ 规则** —— 覆盖密钥泄露、命令注入、文件篡改、网络滥用、信息窃取
- 🤖 **CI/CD 友好** —— 结构化 JSON 输出 + 语义化退出码

---

## 🚀 安装

### macOS（Homebrew）
```bash
brew install wjames2000/tap/skill-guard
```

### Linux / macOS（一键安装脚本）

```bash
curl -sfL https://github.com/wjames2000/skill-guard/releases/latest/install.sh | sh
```

### Windows（Scoop）
```bash
scoop bucket add hpds.cc https://github.com/wjames2000/scoop-bucket.git
scoop install skill-guard
```

### Go 安装
```bash
go install github.com/wjames2000/skill-guard@latest
```

> 预编译的 Linux、macOS（Intel 与 Apple Silicon）以及 Windows 二进制文件可在 [Releases 页面](https://github.com/wjames2000/skill-guard/releases) 下载。

---

## ⚙️ 快速上手

```bash
# 克隆一个技能仓库
git clone https://github.com/someuser/awesome-skill.git

# 进行安全扫描
skill-guard ./awesome-skill

# ✅ 绿色  → 安全，可放心使用
# ❌ 红色  → 发现高危风险，请仔细审查
```

---

## 📋 核心功能

### 🔍 目录扫描
递归扫描所有文本文件，并自动跳过 `.git`、`node_modules`、`__pycache__` 等无关目录。

```bash
skill-guard                          # 扫描当前目录
skill-guard ./downloaded-skills      # 扫描指定目录
skill-guard ./skills ./community-skills   # 同时扫描多个目录
skill-guard ./repo --ignore "tests/**" --ignore "examples/**"
```

### 🧠 规则检测引擎（内置 30+ 规则）

| 类别               | 检测目标                                    |
| ------------------ | ------------------------------------------- |
| **密钥泄露**       | AWS 密钥、GitHub 令牌、私钥、数据库连接串   |
| **命令执行**       | `os.system`、`exec()`、`eval()`、反引号执行 |
| **恶意文件操作**   | `rm -rf`、任意文件写入、覆盖系统文件        |
| **网络请求滥用**   | `curl | bash`、反向 Shell、可疑的外部请求   |
| **信息窃取与混淆** | Base64 混淆、SSH 私钥窃取、环境变量泄露     |

### ✏️ 自定义规则
通过 YAML 或 JSON 文件加载自定义规则，覆盖或扩充内置规则集：
```bash
skill-guard ./skills --rules my-team-rules.yaml
```

### 📤 输出模式
```bash
skill-guard ./skills                    # 终端彩色输出（默认）
skill-guard ./skills --json             # 适合 CI 的结构化 JSON
skill-guard ./skills --quiet            # 仅显示存在问题的文件
skill-guard ./skills --severity high    # 按严重级别过滤（low、medium、high、critical）
```

### 🔧 配置文件
```bash
skill-guard --config .skillguard.yaml
```
在项目级配置文件中预设默认参数、忽略路径和规则文件。

### 🔒 文件过滤
```bash
skill-guard --ext-include .py,.sh      # 仅扫描指定扩展名的文件
skill-guard --ext-exclude .md,.txt     # 排除指定扩展名的文件
skill-guard --ignore "vendor/**"       # 忽略特定路径（glob 模式）
```

---

## 🔐 安全声明

- **不执行被扫描文件** —— 仅读取内容进行分析，绝不执行或解释
- **完全离线** —— 扫描过程中不产生任何网络请求
- **不收集数据** —— 不会将任何使用数据或扫描结果发送至外部

---

## 🗺️ 路线图

| 版本 | 重点内容                                       |
| ---- | ---------------------------------------------- |
| v1.0 | CLI 框架、目录遍历、核心规则引擎、30+ 内置规则 |
| v1.1 | 自定义规则、配置文件、JSON/Quiet 输出          |
| v1.2 | 性能与并发优化，brew/scoop 分发                |
| v2.0 | 规则远程更新、Lua 脚本规则、AI 辅助检测        |

---

## 🧪 CI/CD 集成（GitHub Actions 示例）

```yaml
- name: 扫描技能中的安全风险
  run: |
    skill-guard ./skills --json --severity high
  continue-on-error: false
```

当发现符合严重级别阈值的问题时，`skill-guard` 会以非零状态码退出，便于阻止有风险的 PR。

---

## 🤝 参与贡献

我们欢迎新的规则、Bug 报告和功能建议！  
请参阅[贡献指南](CONTRIBUTING.md)（即将推出）或直接提交 Issue。

---

## 📜 许可证

MIT © 2026 [wjames2000](https://github.com/wjames2000)

---

*skill‑guard —— 扫描之后再信任。*
