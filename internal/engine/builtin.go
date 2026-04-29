package engine

import pkgtypes "github.com/wjames2000/skill-guard/pkg/types"

func BuiltinRules() []*pkgtypes.Rule {
	return []*pkgtypes.Rule{
		// === 密钥泄露类（7 条）===
		{
			ID: "SKL-001", Name: "硬编码 AWS Access Key",
			Severity: "Critical", Pattern: `(?i)AKIA[A-Z0-9]{16}`,
			Keywords: []string{"AKIA"},
		},
		{
			ID: "SKL-002", Name: "私钥文件泄露",
			Severity: "Critical", Pattern: `-----BEGIN (RSA|OPENSSH|EC|DSA) PRIVATE KEY-----`,
			Keywords: []string{"BEGIN PRIVATE KEY"},
		},
		{
			ID: "SKL-009", Name: "硬编码 GitHub Token",
			Severity: "Critical", Pattern: `(?i)(ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9]{36}`,
			Keywords: []string{"ghp_", "gho_", "ghu_", "ghs_", "ghr_"},
		},
		{
			ID: "SKL-010", Name: "硬编码 Google API Key",
			Severity: "Critical", Pattern: `(?i)AIza[0-9A-Za-z\-_]{35}`,
			Keywords: []string{"AIza"},
		},
		{
			ID: "SKL-011", Name: "数据库连接串含密码",
			Severity: "High", Pattern: `(mysql|postgres|mongodb)://[^:]+:[^@]+@`,
			Keywords: []string{"mysql://", "postgres://", "mongodb://"},
		},
		{
			ID: "SKL-012", Name: "硬编码 Slack Token",
			Severity: "High", Pattern: `xox[baprs]-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-z0-9]{32}`,
			Keywords: []string{"xoxb-", "xoxp-", "xoxa-", "xoxr-", "xoxs-"},
		},
		{
			ID: "SKL-013", Name: "JWT Token 硬编码",
			Severity: "Medium", Pattern: `eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}`,
			Keywords: []string{"eyJ"},
		},

		// === 命令执行类（7 条）===
		{
			ID: "SKL-003", Name: "可疑命令执行（Python）",
			Severity: "High", Pattern: `os\.system\s*\(`,
			Keywords: []string{"os.system"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-014", Name: "可疑命令执行（subprocess）",
			Severity: "High", Pattern: `subprocess\.(call|Popen|run)\s*\(`,
			Keywords: []string{"subprocess."}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-015", Name: "Shell 命令执行（Node.js）",
			Severity: "High", Pattern: `child_process\.exec|execSync|execFile`,
			Keywords: []string{"child_process"}, FileTypes: []string{".js", ".ts"},
		},
		{
			ID: "SKL-016", Name: "eval 执行（Python/JS）",
			Severity: "High", Pattern: `\beval\s*\(`,
			Keywords: []string{"eval("}, FileTypes: []string{".py", ".js", ".ts"},
		},
		{
			ID: "SKL-017", Name: "反引号 Shell 执行",
			Severity: "Medium", Pattern: "`[a-z]{2,10}\\s+[-a-zA-Z0-9]",
			Keywords: []string{"`"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-018", Name: "exec 系统调用（Go）",
			Severity: "Medium", Pattern: `exec\.Command\s*\(`,
			Keywords: []string{"exec.Command"}, FileTypes: []string{".go"},
		},
		{
			ID: "SKL-019", Name: "动态代码执行（Node.js）",
			Severity: "Medium", Pattern: `new Function\s*\(`,
			Keywords: []string{"new Function"}, FileTypes: []string{".js", ".ts"},
		},

		// === 恶意文件操作类（6 条）===
		{
			ID: "SKL-007", Name: "过于宽松的文件权限设置",
			Severity: "Low", Pattern: `chmod\s+777`,
			Keywords: []string{"chmod 777"}, FileTypes: []string{".sh", ".py"},
		},
		{
			ID: "SKL-020", Name: "递归删除操作",
			Severity: "Critical", Pattern: `rm\s+(-rf|-fr|--recursive)`,
			Keywords: []string{"rm -rf"}, FileTypes: []string{".sh", ".py"},
		},
		{
			ID: "SKL-021", Name: "任意文件写入（Python）",
			Severity: "High", Pattern: `open\s*\([^)]*\s*['\"][^'\"]*['\"]\s*,\s*['\"]w['\"]`,
			Keywords: []string{`open(`, `,"w"`, `,'w'`}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-022", Name: "覆盖系统文件",
			Severity: "High", Pattern: `(>.*|\s+tee\s+)(/etc/|/usr/|/boot/)`,
			Keywords: []string{"/etc/", "/usr/", "/boot/"}, FileTypes: []string{".sh"},
		},
		{
			ID: "SKL-023", Name: "删除根目录",
			Severity: "High", Pattern: `rm\s+-rf\s+/\s*$`,
			Keywords: []string{"rm -rf /"}, FileTypes: []string{".sh", ".py"},
		},
		{
			ID: "SKL-024", Name: "fs.writeFile 危险写入",
			Severity: "Medium", Pattern: `fs\.writeFileSync?\s*\(`,
			Keywords: []string{"fs.writeFile"}, FileTypes: []string{".js", ".ts"},
		},

		// === 网络请求滥用类（6 条）===
		{
			ID: "SKL-004", Name: "下载并执行脚本",
			Severity: "High", Pattern: `curl.*\|.*(bash|sh|python)`,
			Keywords: []string{"curl |bash", "curl |sh"}, FileTypes: []string{".sh", ".py"},
		},
		{
			ID: "SKL-025", Name: "wget 下载执行",
			Severity: "High", Pattern: `wget.*-O.*\|.*(bash|sh)`,
			Keywords: []string{"wget"}, FileTypes: []string{".sh"},
		},
		{
			ID: "SKL-026", Name: "反向 Shell",
			Severity: "Critical", Pattern: `bash\s+-i\s*>&\s*/dev/tcp/`,
			Keywords: []string{"/dev/tcp/"}, FileTypes: []string{".sh"},
		},
		{
			ID: "SKL-027", Name: "Python 反向 Shell",
			Severity: "High", Pattern: `socket\.socket.*connect\s*\([^)]*\)`,
			Keywords: []string{"socket.socket", ".connect("}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-028", Name: "urllib 请求外部地址",
			Severity: "Medium", Pattern: `urllib\.request\.urlopen\s*\(`,
			Keywords: []string{"urllib.request"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-029", Name: "requests 请求外部地址",
			Severity: "Medium", Pattern: `requests\.(get|post|put|delete)\s*\(`,
			Keywords: []string{"requests.get", "requests.post"}, FileTypes: []string{".py"},
		},

		// === 信息窃取与混淆类（6 条）===
		{
			ID: "SKL-006", Name: "Base64 编码可疑命令",
			Severity: "Medium", Pattern: `echo\s+[A-Za-z0-9+/=]{20,}\s*\|.*base64.*-d`,
			Keywords: []string{"base64 -d"}, FileTypes: []string{".sh"},
		},
		{
			ID: "SKL-030", Name: "读取 SSH 私钥",
			Severity: "High", Pattern: `cat\s+~/.ssh/`,
			Keywords: []string{"~/.ssh"}, FileTypes: []string{".sh", ".py"},
		},
		{
			ID: "SKL-031", Name: "读取环境变量（批量）",
			Severity: "Medium", Pattern: `env\|grep\|export\s+[A-Z]`,
			Keywords: []string{"env |", "export "}, FileTypes: []string{".sh"},
		},
		{
			ID: "SKL-032", Name: "Hex 编码可疑字符串",
			Severity: "Medium", Pattern: `echo\s+-e\s*['\"]\\x[0-9a-f]{2}`,
			Keywords: []string{"echo -e"}, FileTypes: []string{".sh"},
		},
		{
			ID: "SKL-033", Name: "上传敏感文件",
			Severity: "Medium", Pattern: `curl\s+.*-F\s*['\"]file=@`,
			Keywords: []string{"-F file=@"}, FileTypes: []string{".sh", ".py"},
		},
		{
			ID: "SKL-008", Name: "使用 eval（通用）",
			Severity: "Low", Pattern: `\beval\s*\(`,
			Keywords: []string{"eval("}, FileTypes: []string{".sh", ".js", ".ts", ".py"},
		},
	}
}
