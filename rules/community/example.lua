-- 示例 LUA 规则：检测硬编码密码赋值
-- 这是 LUA 脚本化规则的示例，可以处理正则无法表达的复杂逻辑。
-- 全局 'rule' 表定义了规则的元数据，rule.match() 函数执行检测。

rule = {
    id = "LUA-001",
    name = "LUA 脚本规则示例 - 硬编码密码检测",
    severity = "Medium",
    description = "检测变量名与密码相关且值为硬编码字符串",
    file_types = {".py", ".sh", ".js", ".yaml"}
}

function rule.match(line, filename)
    -- 检查是否是密码赋值
    if line:match("[Pp][Aa][Ss][Ss][Ww][Oo][Rr][Dd]%s*=[^=]") then
        -- 排除环境变量引用
        if not line:match("os%.getenv") and not line:match("process%.env")
           and not line:match("environ%.get") and not line:match("config%.") then
            return true, "检测到硬编码密码赋值"
        end
    end
    return false, ""
end

return rule
