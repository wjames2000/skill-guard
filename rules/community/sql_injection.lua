-- 检测拼接 SQL 查询（比纯正则更精确）
-- LUA 规则可以跨行匹配、追踪变量，处理正则无法表达的复杂模式。

rule = {
    id = "LUA-002",
    name = "LUA 规则：SQL 拼接检测",
    severity = "High",
    description = "检测字符串拼接构建 SQL 查询，排除参数化查询",
    file_types = {".py", ".js", ".java"}
}

function rule.match(line, filename)
    -- 查找 SQL 关键字 + 变量拼接
    local sqlKeywords = {"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER"}
    for _, kw in ipairs(sqlKeywords) do
        if line:find(kw, 1, true) then
            -- 检查是否使用了字符串拼接
            if line:find("%+") or line:find("f\"") or line:find("f'") then
                -- 排除使用参数化查询的场景
                if not line:find("%%s") and not line:find("?")
                   and not line:find(":param") and not line:find("%$1")
                   and not line:find("%$2") and not line:find("%$3") then
                    return true, "检测到可能的 SQL 注入：使用字符串拼接构建 SQL"
                end
            end
        end
    end
    return false, ""
end

return rule
