local key = KEYS[1]
local cntKey = key .. ":cnt"
-- 你准备的存储的验证码
local val = ARGV[1]
local expireTime = 60

local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- key 存在，但是没有过期时间
    return -2
elseif ttl == -2 or ttl < 50 then
    -- 可以发验证码
    redis.call("set", key, val)
    -- 60 秒
    redis.call("expire", key, expireTime)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, expireTime)
    return 0
else
    -- 发送太频繁
    return -1
end