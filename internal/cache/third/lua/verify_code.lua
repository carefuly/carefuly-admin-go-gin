local key = KEYS[1]
local cntKey = key .. ":cnt"
local blockKey = key .. ":block" -- 用户被限制的key
-- 用户输入的验证码
local expectedCode = ARGV[1]
local expireTime = 600

-- 检查用户是否被限制
if redis.call("exists", blockKey) == 1 then
    return -3  -- 返回-3表示用户被限制
end

local cnt = tonumber(redis.call("get", cntKey))
local code = redis.call("get", key)

-- 检查傻逼用户不发送验证码直接登录
if cnt == nil then
    return -4  -- 返回-4没发送验证码
end

if cnt <= 0 then
    -- 验证次数耗尽，设置限制标记
    redis.call("set", blockKey, 1)
    redis.call("expire", blockKey, expireTime)
    -- 验证次数耗尽了
    return -1
end

if code == expectedCode then
    redis.call("del", key)
    redis.call("del", cntKey)
    -- redis.call("set", cntKey, 0)
    return 0
else
    redis.call("decr", cntKey)
    -- 不相等，用户输错了
    return -2
end