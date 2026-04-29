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

		// === 补全原规则（1 条）===
		{
			ID: "SKL-005", Name: "硬编码密码赋值",
			Severity: "Medium", Pattern: `(password|passwd|pwd)\s*[:=]\s*["'][^"'\s]{6,}`,
			Keywords: []string{"password =", "password=", "passwd =", "pwd ="},
		},

		// === 密钥泄露类扩展（10 条）===
		{
			ID: "SKL-034", Name: "硬编码 Azure 订阅密钥",
			Severity: "Critical", Pattern: `(?i)[a-f0-9]{32}`,
			Keywords: []string{"subscription", "azure"},
		},
		{
			ID: "SKL-035", Name: "硬编码 Discord Bot Token",
			Severity: "Critical", Pattern: `(?i)[MN][A-Za-z\d]{23}\.[\w-]{6}\.[\w-]{27}`,
			Keywords: []string{"discord", "token"},
		},
		{
			ID: "SKL-036", Name: "硬编码 Telegram Bot Token",
			Severity: "Critical", Pattern: `[0-9]{8,10}:[A-Za-z0-9_-]{35}`,
			Keywords: []string{"telegram", "bot"},
		},
		{
			ID: "SKL-037", Name: "硬编码 Stripe API Key",
			Severity: "Critical", Pattern: `(?i)sk_live_[0-9a-z]{24,}`,
			Keywords: []string{"sk_live"},
		},
		{
			ID: "SKL-038", Name: "硬编码 Twilio API Key",
			Severity: "Critical", Pattern: `(?i)SK[a-f0-9]{32}`,
			Keywords: []string{"twilio"},
		},
		{
			ID: "SKL-039", Name: "硬编码 Heroku API Key",
			Severity: "High", Pattern: `(?i)[hH][eE][rR][oO][kK][uU].*[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`,
			Keywords: []string{"heroku", "api"},
		},
		{
			ID: "SKL-040", Name: "硬编码 SendGrid API Key",
			Severity: "High", Pattern: `SG\.[a-zA-Z0-9_-]{22}\.[a-zA-Z0-9_-]{43}`,
			Keywords: []string{"SG."},
		},
		{
			ID: "SKL-041", Name: "硬编码阿里云 AccessKey",
			Severity: "High", Pattern: `(?i)LTAI[a-zA-Z0-9]{12,}`,
			Keywords: []string{"LTAI"},
		},
		{
			ID: "SKL-042", Name: "GCP 服务账号密钥泄露",
			Severity: "Critical", Pattern: `"type":\s*"service_account"`,
			Keywords: []string{"service_account", "private_key_id"},
		},
		{
			ID: "SKL-043", Name: "硬编码 npm Token",
			Severity: "High", Pattern: `(?i)npm_[a-z0-9]{36}`,
			Keywords: []string{"npm_"},
		},

		// === 代码注入类（8 条）===
		{
			ID: "SKL-044", Name: "SQL 注入风险（Python）",
			Severity: "High", Pattern: `execute\s*\(.*['"]SELECT|execute\s*\(.*['"]INSERT|execute\s*\(.*['"]DELETE`,
			Keywords: []string{"execute(", "cursor.execute"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-045", Name: "SQL 注入风险（Node.js）",
			Severity: "High", Pattern: `db\.(query|execute)\s*\(['"][^'"]*\+`,
			Keywords: []string{"db.query", "db.execute"}, FileTypes: []string{".js", ".ts"},
		},
		{
			ID: "SKL-046", Name: "Jinja2 模板注入风险",
			Severity: "Medium", Pattern: `render_template_string\s*\(`,
			Keywords: []string{"render_template_string"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-047", Name: "NoSQL 注入（MongoDB）",
			Severity: "Medium", Pattern: `\$where`,
			Keywords: []string{"$where"}, FileTypes: []string{".js", ".ts"},
		},
		{
			ID: "SKL-048", Name: "XSS 风险（JavaScript）",
			Severity: "Medium", Pattern: `innerHTML\s*=|outerHTML\s*=|document\.write\s*\(`,
			Keywords: []string{"innerHTML", "outerHTML"}, FileTypes: []string{".js", ".ts"},
		},
		{
			ID: "SKL-049", Name: "命令注入（shell=True）",
			Severity: "High", Pattern: `subprocess\.(call|Popen|run)\s*\(.*shell\s*=\s*True`,
			Keywords: []string{"shell=True"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-050", Name: "用户输入传入 eval（Node.js）",
			Severity: "High", Pattern: `eval\s*\(.*req\.|eval\s*\(.*body\.|eval\s*\(.*params\.`,
			Keywords: []string{"eval(req.", "eval(body."}, FileTypes: []string{".js", ".ts"},
		},
		{
			ID: "SKL-051", Name: "PHP 代码执行函数",
			Severity: "High", Pattern: `(exec|system|passthru|shell_exec)\s*\(`,
			Keywords: []string{"exec(", "system("}, FileTypes: []string{".php"},
		},

		// === 配置风险类（8 条）===
		{
			ID: "SKL-052", Name: "Flask Debug 模式开启",
			Severity: "Low", Pattern: `debug\s*=\s*True`,
			Keywords: []string{"debug=True"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-053", Name: "CORS 配置过于宽松",
			Severity: "Medium", Pattern: `Access-Control-Allow-Origin:\s*\*`,
			Keywords: []string{"Access-Control-Allow-Origin: *"},
		},
		{
			ID: "SKL-054", Name: "弱 TLS 协议配置",
			Severity: "Medium", Pattern: `ssl_protocols\s+.*TLSv1\b|ssl_protocols\s+.*SSLv`,
			Keywords: []string{"ssl_protocols"},
		},
		{
			ID: "SKL-055", Name: "CSRF 保护被禁用",
			Severity: "Medium", Pattern: `csrf\.exempt|csrf_exempt|WTF_CSRF_ENABLED\s*=\s*False`,
			Keywords: []string{"csrf.exempt", "csrf_exempt"},
		},
		{
			ID: "SKL-056", Name: "硬编码 Flask SECRET_KEY",
			Severity: "High", Pattern: `SECRET_KEY\s*=\s*['"][^'"]+['"]`,
			Keywords: []string{"SECRET_KEY"},
		},
		{
			ID: "SKL-057", Name: "Django DEBUG 模式开启",
			Severity: "Low", Pattern: `DEBUG\s*=\s*True`,
			Keywords: []string{"DEBUG = True"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-058", Name: "Cookie 安全标记未启用",
			Severity: "Medium", Pattern: `cookie\.secure\s*=\s*false|Secure\s*=\s*false`,
			Keywords: []string{"Secure=false"},
		},
		{
			ID: "SKL-059", Name: "生产环境暴露 API 文档",
			Severity: "Medium", Pattern: `api-docs|swagger\.json|openapi\.json`,
			Keywords: []string{"swagger.json", "openapi.json"},
		},

		// === 加密与认证类（8 条）===
		{
			ID: "SKL-060", Name: "弱密码哈希（MD5）",
			Severity: "High", Pattern: `hashlib\.md5\b`,
			Keywords: []string{"hashlib.md5"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-061", Name: "弱密码哈希（SHA1）",
			Severity: "Medium", Pattern: `hashlib\.sha1\b`,
			Keywords: []string{"hashlib.sha1"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-062", Name: "硬编码 JWT Secret",
			Severity: "Critical", Pattern: `jwt\.(encode|decode|sign)\s*\(.*['"][A-Za-z0-9!@#$%^&*()_+=-]{10,}`,
			Keywords: []string{"jwt.encode", "jwt.sign"},
		},
		{
			ID: "SKL-063", Name: "硬编码 SSL 证书",
			Severity: "Medium", Pattern: `-----BEGIN CERTIFICATE-----`,
			Keywords: []string{"BEGIN CERTIFICATE"},
		},
		{
			ID: "SKL-064", Name: "弱 RSA 密钥长度（1024）",
			Severity: "High", Pattern: `openssl\s+genrsa\s+-out\s+\S+\s+1024\b`,
			Keywords: []string{"genrsa 1024"},
		},
		{
			ID: "SKL-065", Name: "硬编码加密密钥",
			Severity: "Critical", Pattern: `encryption_key\s*=\s*['"][A-Za-z0-9!@#$%^&*()_+=-]{16,}`,
			Keywords: []string{"encryption_key"},
		},
		{
			ID: "SKL-066", Name: "使用 ECB 加密模式",
			Severity: "Medium", Pattern: `AES\.MODE_ECB|DES\.MODE_ECB`,
			Keywords: []string{"MODE_ECB"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-067", Name: "硬编码 Basic Auth 凭据",
			Severity: "High", Pattern: `Authorization:\s*Basic\s+[A-Za-z0-9+/=]{20,}`,
			Keywords: []string{"Authorization: Basic"},
		},

		// === 依赖与供应链风险类（7 条）===
		{
			ID: "SKL-068", Name: "依赖固定为 latest 版本",
			Severity: "Low", Pattern: `"[^"]+":\s*"\*"|'[^']+':\s*'\*'`,
			Keywords: []string{`": "*"`, `': '*'`},
		},
		{
			ID: "SKL-069", Name: "从 GitHub 直接安装依赖",
			Severity: "Medium", Pattern: `pip\s+install\s+git\+https|npm\s+install\s+https://github`,
			Keywords: []string{"pip install git+https", "npm install https://github"},
		},
		{
			ID: "SKL-070", Name: "curl 管道安装（供应链风险）",
			Severity: "High", Pattern: `curl.*\|.*pip\s+install`,
			Keywords: []string{"curl | pip install"},
		},
		{
			ID: "SKL-071", Name: "npm postinstall 脚本风险",
			Severity: "Medium", Pattern: `"postinstall":\s*"[^"]*"`,
			Keywords: []string{`"postinstall"`},
		},
		{
			ID: "SKL-072", Name: "Base64 解码后执行（Python）",
			Severity: "High", Pattern: `exec\(.*base64\.(b64decode|decodestring)`,
			Keywords: []string{"exec(base64"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-073", Name: "导入不安全的网络模块",
			Severity: "Medium", Pattern: `import\s+(telnetlib|ftplib|poplib)`,
			Keywords: []string{"import telnetlib"},
		},
		{
			ID: "SKL-074", Name: "使用已弃用的 ioutil 包",
			Severity: "Low", Pattern: `ioutil\.`,
			Keywords: []string{"ioutil."}, FileTypes: []string{".go"},
		},

		// === 信息泄露类（8 条）===
		{
			ID: "SKL-075", Name: "堆栈跟踪信息泄露",
			Severity: "Medium", Pattern: `traceback\.print_exc|print_exc\(\)`,
			Keywords: []string{"traceback.print_exc"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-076", Name: "生产环境调试端点",
			Severity: "Low", Pattern: `@app\.route\(['"]/(debug|test)`,
			Keywords: []string{"/debug", "/test"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-077", Name: "硬编码内网 IP 地址",
			Severity: "Medium", Pattern: `(10\.\d{1,3}\.\d{1,3}\.\d{1,3}|172\.(1[6-9]|2\d|3[01])\.\d{1,3}\.\d{1,3}|192\.168\.\d{1,3}\.\d{1,3})`,
			Keywords: []string{"10.", "192.168."},
		},
		{
			ID: "SKL-078", Name: "硬编码邮件 SMTP 凭据",
			Severity: "Medium", Pattern: `smtp\.(login|password)\s*=\s*['"][^'"]+`,
			Keywords: []string{"smtp.login", "smtp.password"},
		},
		{
			ID: "SKL-079", Name: "硬编码 LDAP 凭据",
			Severity: "High", Pattern: `ldap\.(bind|simple_bind)\s*\(.*['"][^'"]+['"]`,
			Keywords: []string{"ldap.bind", "ldap.simple_bind"},
		},
		{
			ID: "SKL-080", Name: "暴露 .env 文件",
			Severity: "Low", Pattern: `\.env`,
			Keywords: []string{},
		},
		{
			ID: "SKL-081", Name: "生产环境 Debug 日志开启",
			Severity: "Low", Pattern: `logging\.(basicConfig|setLevel)\(.*DEBUG`,
			Keywords: []string{"logging.DEBUG"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-082", Name: "硬编码社交媒体 API 密钥",
			Severity: "High", Pattern: `(?i)(facebook|twitter|instagram|linkedin).*(secret|token|key)\s*[:=]\s*['"][^'"]+`,
			Keywords: []string{"facebook_secret", "twitter_token"},
		},

		// === 容器与基础设施风险类（8 条）===
		{
			ID: "SKL-083", Name: "容器以 root 用户运行",
			Severity: "Medium", Pattern: `USER\s+root\b`,
			Keywords: []string{"USER root"},
		},
		{
			ID: "SKL-084", Name: "特权容器模式",
			Severity: "High", Pattern: `privileged:\s*true`,
			Keywords: []string{"privileged: true"},
		},
		{
			ID: "SKL-085", Name: "容器使用主机网络",
			Severity: "Medium", Pattern: `network_mode:\s*host\b`,
			Keywords: []string{"network_mode: host"},
		},
		{
			ID: "SKL-086", Name: "挂载 Docker 套接字",
			Severity: "High", Pattern: `/var/run/docker\.sock`,
			Keywords: []string{"docker.sock"},
		},
		{
			ID: "SKL-087", Name: "Dockerfile 从远程 URL 下载文件",
			Severity: "Medium", Pattern: `ADD\s+https?://`,
			Keywords: []string{"ADD http", "ADD https"},
		},
		{
			ID: "SKL-088", Name: "Kubernetes 共享主机 PID/网络",
			Severity: "High", Pattern: `hostPID:\s*true|hostNetwork:\s*true`,
			Keywords: []string{"hostPID: true", "hostNetwork: true"},
		},
		{
			ID: "SKL-089", Name: "Dockerfile 暴露 SSH 端口",
			Severity: "Low", Pattern: `EXPOSE\s+22\b`,
			Keywords: []string{"EXPOSE 22"},
		},
		{
			ID: "SKL-090", Name: "硬编码容器仓库登录凭据",
			Severity: "Critical", Pattern: `(docker|podman)\s+login\s+-u\s+\S+\s+-p\s+\S+`,
			Keywords: []string{"docker login -u"},
		},

		// === 行为监控与后门类（10 条）===
		{
			ID: "SKL-091", Name: "C2 回调模式（IP:Port）",
			Severity: "Critical", Pattern: `(http|https)://(\d{1,3}\.){3}\d{1,3}:[0-9]{4,5}`,
			Keywords: []string{"http://"},
		},
		{
			ID: "SKL-092", Name: "DNS 数据外泄（nslookup）",
			Severity: "High", Pattern: `nslookup\s+`,
			Keywords: []string{"nslookup "}, FileTypes: []string{".sh"},
		},
		{
			ID: "SKL-093", Name: "下载至共享内存目录",
			Severity: "High", Pattern: `wget.*-O\s+/dev/shm/|curl.*-o\s+/dev/shm/`,
			Keywords: []string{"/dev/shm/"},
		},
		{
			ID: "SKL-094", Name: "修改 crontab 定时任务",
			Severity: "High", Pattern: `crontab\s+-e\b|crontab\s+.*\.tmp`,
			Keywords: []string{"crontab -e"}, FileTypes: []string{".sh"},
		},
		{
			ID: "SKL-095", Name: "SSRF 漏洞（请求内网地址）",
			Severity: "Medium", Pattern: `requests\.(get|post)\s*\(.*['"](http://localhost|http://127|http://10\.|http://172\.|http://192\.)`,
			Keywords: []string{"http://localhost", "http://127."}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-096", Name: "Pickle 反序列化风险",
			Severity: "High", Pattern: `pickle\.loads|pickle\.load\(`,
			Keywords: []string{"pickle.loads"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-097", Name: "不安全的 YAML 加载",
			Severity: "High", Pattern: `yaml\.load\(`,
			Keywords: []string{"yaml.load("}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-098", Name: "XML 外部实体注入（XXE）",
			Severity: "Medium", Pattern: `xml\.parsers\.expat|xml\.dom\.minidom`,
			Keywords: []string{"xml.parsers"}, FileTypes: []string{".py"},
		},
		{
			ID: "SKL-099", Name: "sudo curl/wget 管道执行",
			Severity: "High", Pattern: `sudo.*curl.*\|.*bash|sudo.*wget.*\|.*sh`,
			Keywords: []string{"sudo curl", "sudo wget"}, FileTypes: []string{".sh"},
		},
		{
			ID: "SKL-100", Name: "动态模块加载",
			Severity: "Medium", Pattern: `__import__\s*\(|importlib\.import_module\s*\(`,
			Keywords: []string{"__import__(", "importlib.import_module"}, FileTypes: []string{".py"},
		},
	}
}
