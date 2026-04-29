# 🛡️ skill-guard

**Zero‑trust security scanner for AI Skills**  
*Scan third‑party skills before you run them — fast, offline, no dependencies.*

[![version](https://img.shields.io/badge/version-v0.2.0-blue)](https://github.com/wjames2000/skill-guard/releases)
[![go version](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)](https://golang.org)
[![license](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![tests](https://img.shields.io/badge/tests-71%20passing-brightgreen)](https://github.com/wjames2000/skill-guard/actions)
[![gosec](https://img.shields.io/badge/gosec-0%20issues-brightgreen)](https://github.com/wjames2000/skill-guard)
[![platform](https://img.shields.io/badge/platform-linux%20|%20macOS%20|%20windows-lightgrey)](https://github.com/wjames2000/skill-guard/releases)

English | [中文](README_ZH.md)

---

skill-guard detects hidden threats inside third-party AI Skills (`.md`, `.json`, `.yaml`, `.py`, `.sh`, etc.) — leaked secrets, malicious commands, obfuscated payloads, and more.

- 🔒 **Zero trust** — assume every Skill is dangerous until proven safe
- ⚡ **Zero dependencies** — single native binary, no runtimes required
- 📦 **100 built‑in rules** — 12 categories: secrets, injections, crypto, supply chain, containers, etc.
- 🤖 **AI‑assisted scanning** — optional LLM verification to reduce false positives
- 🚀 **CI/CD ready** — JSON / SARIF output, semantic exit codes, GitHub Action

---

## 🚀 Installation

### macOS (Homebrew)
```bash
brew install wjames2000/tap/skill-guard
```

### Linux / macOS (one‑line installer)
```bash
curl -sfL https://github.com/wjames2000/skill-guard/releases/latest/download/install.sh | sh
```

### Windows (Scoop)
```bash
scoop bucket add wjames2000 https://github.com/wjames2000/scoop-bucket.git
scoop install skill-guard
```

### Go install
```bash
go install github.com/wjames2000/skill-guard@latest
```

> Pre‑compiled binaries for Linux (amd64/arm64), macOS (amd64/arm64), and Windows (amd64) on the [Releases page](https://github.com/wjames2000/skill-guard/releases).

---

## ⚙️ Quick Start

```bash
# Scan a directory of skills
skill-guard ./downloaded-skills

# Scan with AI assistance to reduce false positives
skill-guard ./skills --ai

# Output SARIF for GitHub Code Scanning
skill-guard ./skills --sarif

# Export all 100 built‑in rules
skill-guard rules export yaml > my-rules-backup.yaml

# Initialize a project with pre-commit hook
cd my-project && skill-guard init
```

---

## 📋 Features

### 🔍 Smart Directory Scanning
Recursively scans text files, auto‑skips `.git`, `node_modules`, `vendor`, etc.
Respects `.gitignore` rules with `--gitignore` flag. Skips symlinks.

```bash
skill-guard                              # scan current directory
skill-guard ./path/to/skills             # scan a specific directory
skill-guard ./a ./b                      # scan multiple directories
skill-guard --gitignore                  # respect .gitignore rules
```

### 🧠 100 Built‑in Rules (12 Categories)

| Category | Rules | Examples |
|----------|-------|---------|
| **Secret leaks** | 17 | AWS key, GitHub token, JWT, Stripe, GCP SA, Azure, Discord, npm, Slack, SendGrid |
| **Command injection** | 7 | `os.system`, `subprocess`, `child_process`, `eval`, `exec`, `new Function` |
| **Code injection** | 8 | SQLi, NoSQLi, XSS, SSTI, PHP code exec, shell=True |
| **Malicious file ops** | 6 | `rm -rf`, arbitrary writes, chmod 777, system file overwrite |
| **Network abuse** | 6 | `curl | bash`, reverse shell, wget pipe, SSRF |
| **Info theft / obfuscation** | 6 | Base64, Hex, SSH key read, env leak, file upload |
| **Config risks** | 8 | Debug mode, CORS `*`, weak TLS, CSRF disabled, insecure cookies |
| **Crypto & auth** | 8 | MD5/SHA1, ECB mode, JWT secret, hardcoded cert, weak RSA |
| **Supply chain** | 7 | `pip install git+https`, `"*"` deps, npm postinstall, curl pip |
| **Info disclosure** | 8 | Stack traces, debug endpoints, hardcoded IPs, SMTP creds, logging |
| **Container risks** | 8 | Privileged, Docker socket, host network, `USER root`, exposed SSH |
| **Backdoor / C2** | 10 | C2 callbacks, crontab writes, pickle deserialize, yaml.load, sudo pipe |

### 🤖 AI‑Assisted Scanning
Pass findings through a local LLM for semantic verification, reducing false positives:
```bash
skill-guard ./skills --ai                              # default: gemma-4-26b-a4b-it
skill-guard ./skills --ai --ai-model llama3.2           # custom model
skill-guard ./skills --ai --ai-endpoint https://my-api/v1  # custom endpoint
```

### ✏️ Custom Rules & Rule Marketplace
```bash
# Load custom rules
skill-guard ./skills --rules my-rules.yaml

# Create a new rule template
skill-guard rules new

# Test a rule against sample input
skill-guard rules test my-rule.yaml "dangerous_code()"

# Browse rule marketplace
skill-guard update info
skill-guard update list
skill-guard update install <url>
```

### 📤 Output Modes
```bash
skill-guard ./skills                              # colored terminal (default)
skill-guard ./skills --json                       # structured JSON for CI
skill-guard ./skills --quiet                      # only file paths with issues
skill-guard ./skills --summary                    # statistics only
skill-guard ./skills --sarif                      # SARIF 2.1 (GitHub Code Scanning)
skill-guard ./skills --no-color                   # plain text for CI logs
skill-guard ./skills --output report.json         # write JSON to file
skill-guard ./skills --severity high              # filter: low/medium/high/critical
```

### 🔧 Configuration File
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
skill-guard init          # auto-install pre-commit hook
```
Or use with pre-commit framework:
```yaml
repos:
  - repo: https://github.com/wjames2000/skill-guard
    rev: v0.2.0
    hooks:
      - id: skill-guard
        args: ["--severity", "high", "--quiet"]
```

### 🐚 Shell Completion
```bash
skill-guard completion bash   # generate bash completion
skill-guard completion zsh    # generate zsh completion
skill-guard completion fish   # generate fish completion
```

---

## 🧪 CI/CD Integration

### GitHub Actions
```yaml
- uses: wjames2000/skill-guard@v0.2.0
  with:
    path: .
    severity: high
    format: sarif
```

### Standalone
```yaml
- name: Security scan
  run: |
    skill-guard ./skills --json --severity high
```

Exits with code `1` when issues found — blocks risky PRs naturally.

---

## 📊 Performance

| Scenario | Time | Memory |
|----------|------|--------|
| 100 rules × 1 file | ~0.14 ms | 68 KB |
| 100 files scan | ~7.9 ms | 7 MB |
| 500 files scan | ~44 ms | 35 MB |
| 50 MB binary | ~0.045 s (skipped) | — |

---

## 🔐 Security Guarantees

- **Never executes files** — read‑only analysis
- **Fully offline** — zero network calls during scan (`update` command excluded)
- **No telemetry** — no data ever sent anywhere
- **gosec: 0 issues** — static analysis pass

---

## 🗺️ Roadmap

| Version | Focus |
|---------|-------|
| v0.1.0 | CLI, directory traversal, 32 rules, concurrent scan, exit codes |
| v1.1 | Quality: scanner tests, panic protection, benchmarks, Windows compat, symlink, progress bar, `.gitignore` |
| v1.2 | Experience: shell completion, `--no-color`, `--output`, `--summary`, SARIF |
| v2.0-α | Ecosystem: rule marketplace, AI detection, `skill-guard init` |
| v2.0-β | 100 rules (12 categories), 71 tests, gosec 0 |
| v2.1 | Quality: 66 new rule tests, supplements, stress benchmarks |
| v2.2 | Rules: `rules export\|new\|test`, version management, contribution guide |
| v2.3 | Platform: GitHub Action, pre-commit hooks, Homebrew, CI docs |
| v2.4 | Engine: LUA scripting, YARA, semantic analysis |
| v3.0 | Platform: Web marketplace, team management, API service |

---

## 🤝 Contributing

- [Rule Contribution Guide](docs/规则贡献指南.md)
- [CI/CD Integration Guide](docs/ci-integration.md)
- [Homebrew Submission Guide](docs/homebrew-submission.md)
- Open an [issue](https://github.com/wjames2000/skill-guard/issues) or PR

---

## 📜 License

MIT © 2026 [wjames2000](https://github.com/wjames2000)

---

*skill-guard — Scan before you trust.*
