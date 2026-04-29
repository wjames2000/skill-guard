# Homebrew 官方仓库提交指南

## 前提条件

1. Release 已发布且可下载
2. 已通过 `shasum -a 256` 获取各平台二进制 SHA256

## 提交流程

### 1. 获取 SHA256

```bash
# 构建所有平台
bash scripts/build.sh

# 计算 SHA256
shasum -a 256 dist/skill-guard-darwin-amd64
shasum -a 256 dist/skill-guard-darwin-arm64
shasum -a 256 dist/skill-guard-linux-amd64
shasum -a 256 dist/skill-guard-linux-arm64
shasum -a 256 dist/skill-guard-windows-amd64.exe
```

### 2. 提交至 homebrew-core

```bash
brew tap homebrew/core
brew create --set-version 0.2.0 \
  https://github.com/wjames2000/skill-guard/releases/download/v0.2.0/skill-guard-darwin-amd64
```

### 3. 编辑 Formula

编辑生成的 Formula 文件，参考 `contrib/homebrew/skill-guard.rb`。

### 4. 提交 PR

```bash
# Fork homebrew-core
git clone https://github.com/YOUR_USERNAME/homebrew-core
cd homebrew-core
git checkout -b skill-guard-0.2.0

# 创建 Formula
cp /path/to/formula Formula/s/skill-guard.rb

# 提交
git add Formula/s/skill-guard.rb
git commit -m "skill-guard 0.2.0: security scanner for AI skill files"
git push origin skill-guard-0.2.0
```

### 审核要求

- Formula 必须通过 `brew audit --strict skill-guard`
- 必须通过 `brew test skill-guard`
- SHA256 必须与 Release 附件完全一致

### 安装验证

```bash
brew install skill-guard
skill-guard --version
skill-guard --help
```
