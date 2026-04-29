# 🛡️ skill-guard

**Zero‑trust security scanner for AI Skills**  
*Scan third‑party skills before you run them – fast, offline, no dependencies.*

[![version](https://img.shields.io/badge/version-1.0-blue)](https://github.com/wjames2000/skill-guard)
[![go version](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)](https://golang.org)
[![license](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![platform](https://img.shields.io/badge/platform-linux%20|%20macOS%20|%20windows-lightgrey)](https://github.com/wjames2000/skill-guard/releases)

English | [中文](https://github.com/wjames2000/skill-guard/blob/main/README_ZH.md)

---

**skill-guard** is a command‑line security scanner purpose‑built for AI‑assisted development ecosystems.  
It detects hidden threats inside third‑party Skills (Markdown, JSON, YAML, Python, etc.) – leaked secrets, malicious commands, obfuscated payloads, and more.

- 🔒 **Zero trust** – assume every Skill is dangerous until proven safe.
- ⚡ **Zero dependencies** – a single native binary, no runtimes required.
- 📦 **30+ built‑in rules** – secrets, command injection, file tampering, network abuse, information theft.
- 🤖 **CI/CD ready** – structured JSON output + semantic exit codes.

---

## 🚀 Installation

### macOS (Homebrew)
```bash
brew install wjames2000/tap/skill-guard
```

### Linux / macOS (shell installer)

```bash
curl -sfL https://github.com/wjames2000/skill-guard/releases/latest/install.sh | sh
```

### Windows (Scoop)
```bash
scoop bucket add hpds.cc https://github.com/wjames2000/scoop-bucket.git
scoop install skill-guard
```

### Go install
```bash
go install github.com/wjames2000/skill-guard@latest
```

> Pre‑compiled binaries for Linux, macOS (Intel & Apple Silicon), and Windows are available on the [Releases page](https://github.com/wjames2000/skill-guard/releases).

---

## ⚙️ Quick Start

```bash
# Clone a Skill repository
git clone https://github.com/someuser/awesome-skill.git

# Scan it
skill-guard ./awesome-skill

# ✅ green  → safe to use
# ❌ red    → high‑risk items found – review carefully
```

---

## 📋 Features

### 🔍 Directory scanning
Recursively scans all text files while automatically skipping non‑essential directories like `.git`, `node_modules`, `__pycache__`, etc.

```bash
skill-guard                          # scan current directory
skill-guard ./downloaded-skills      # scan a specific directory
skill-guard ./skills ./community-skills   # scan multiple directories
skill-guard ./repo --ignore "tests/**" --ignore "examples/**"
```

### 🧠 Rule engine (30+ built‑in)

| Category                     | Detects                                                      |
| ---------------------------- | ------------------------------------------------------------ |
| **Secret leaks**             | AWS keys, GitHub tokens, private keys, DB connection strings |
| **Command execution**        | `os.system`, `exec()`, `eval()`, backtick shell invocations  |
| **Malicious file ops**       | `rm -rf`, arbitrary writes, overwriting system files         |
| **Network abuse**            | `curl | bash`, reverse shells, suspicious outbound requests  |
| **Info theft / obfuscation** | Base64 blobs, SSH key exfiltration, env‑var leaks            |

### ✏️ Custom rules
Load your own YAML or JSON rules to extend or override the built‑in set:
```bash
skill-guard ./skills --rules my-team-rules.yaml
```

### 📤 Output modes
```bash
skill-guard ./skills                    # colored terminal output (default)
skill-guard ./skills --json             # structured JSON for CI
skill-guard ./skills --quiet            # show only files with findings
skill-guard ./skills --severity high    # filter by severity (low, medium, high, critical)
```

### 🔧 Configuration file
```bash
skill-guard --config .skillguard.yaml
```
Store default flags, ignored paths, and rule files in a project‑level configuration.

### 🔒 File filtering
```bash
skill-guard --ext-include .py,.sh      # only scan these extensions
skill-guard --ext-exclude .md,.txt     # exclude specific extensions
skill-guard --ignore "vendor/**"       # ignore paths (glob patterns)
```

---

## 🔐 Security guarantees

- **Never executes files** – contents are only read and analysed.
- **Fully offline** – zero network calls during a scan.
- **No telemetry** – no data is ever sent anywhere.

---

## 🗺️ Roadmap

| Version | Focus                                                        |
| ------- | ------------------------------------------------------------ |
| v1.0    | CLI, directory traversal, core rule engine, 30+ built‑in rules |
| v1.1    | Custom rules, configuration file, JSON/Quiet output          |
| v1.2    | Performance & concurrency improvements, brew/scoop distribution |
| v2.0    | Remote rule updates, Lua scripting, AI‑assisted detection    |

---

## 🧪 CI/CD Integration (GitHub Actions example)

```yaml
- name: Scan skills for security risks
  run: |
    skill-guard ./skills --json --severity high
  continue-on-error: false
```

`skill-guard` exits with a non‑zero code when findings match the severity threshold, making it easy to block risky PRs.

---

## 🤝 Contributing

We welcome new rules, bug reports, and feature requests!  
Check out the [Contribution Guide](CONTRIBUTING.md) (coming soon) or open an issue.

---

## 📜 License

MIT © 2026 [wjames2000](https://github.com/wjames2000)

---

*skill‑guard — Scan before you trust.*