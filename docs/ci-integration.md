# CI/CD 集成指南

## GitHub Actions

```yaml
# .github/workflows/security-scan.yml
name: Security Scan

on:
  push:
    paths:
      - '**/*.py'
      - '**/*.sh'
      - '**/*.yaml'
      - '**/*.json'
  pull_request:

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: wjames2000/skill-guard@v0.2.0
        with:
          path: .
          severity: high
          format: sarif
```

## Pre-commit Hook

在项目根目录创建 `.pre-commit-config.yaml`：

```yaml
repos:
  - repo: https://github.com/wjames2000/skill-guard
    rev: v0.2.0
    hooks:
      - id: skill-guard
        args: ["--severity", "high", "--quiet"]
```

然后运行：

```bash
pip install pre-commit
pre-commit install
```

## GitLab CI

```yaml
security-scan:
  stage: test
  script:
    - curl -sfL https://github.com/wjames2000/skill-guard/releases/latest/download/install.sh | sh
    - skill-guard . --severity high --json
  artifacts:
    reports:
      json: gl-sast-report.json
```

## Jenkins Pipeline

```groovy
stage('Security Scan') {
    steps {
        sh 'curl -sfL https://github.com/wjames2000/skill-guard/releases/latest/download/install.sh | sh'
        sh 'skill-guard . --severity high --json --output sast-report.json'
    }
    post {
        always {
            archiveArtifacts artifacts: 'sast-report.json'
        }
    }
}
```
