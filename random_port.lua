local ports = {"1024", "1025"}

request = function()
    local port = ports[math.random(#ports)]
    local key = "key_" .. math.random(1000000)
    local path = "/Scores/" .. key
    wrk.headers["Host"] = "127.0.0.1:" .. port
    return wrk.format(nil, path)
end