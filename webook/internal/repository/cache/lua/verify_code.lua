local key = KEYS[1]
--用户输入的code
local expectedCode = ARGV[1]
local code = redis.call("get",key)
local cntKey = key..":cnt"


--转成一个数字
local cnt = tonumber(redis.call("get",cntKey))
if cnt <= 0 then
    --说明，用户一直输错
    --或者已经用了
    return -1
elseif expectedCode==code then
    -- 说明对了
    --用完不能再用
    redis.call("set",cntKey,-1)
    return 0
else
    --可验证次数减一
    redis.call("decr",cntKey)
    --用户手一抖，输错了
    return -2
end