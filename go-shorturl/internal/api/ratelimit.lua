-- File: middleware/ratelimit.lua
-- KEYS[1] = The unique key (e.g., rate_limit:127.0.0.1)
-- ARGV[1] = Burst Capacity (e.g., 10 tokens)
-- ARGV[2] = Refill Rate (e.g., 1 token per second)
-- ARGV[3] = Current Unix Time (Seconds)
-- ARGV[4] = Tokens requested (usually 1)

local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

-- Get current state (tokens and last_refill_time)
local bucket = redis.call("hmget", key, "tokens", "last_refill")
local tokens = tonumber(bucket[1])
local last_refill = tonumber(bucket[2])

-- Initialize if key doesn't exist
if not tokens then
    tokens = capacity
    last_refill = now
end

-- Calculate refill
-- (now - last_refill) * rate = new tokens to add
local delta = math.max(0, now - last_refill)
local filled_tokens = math.min(capacity, tokens + (delta * rate))

-- Check if allowed
local allowed = 0
if filled_tokens >= requested then
    filled_tokens = filled_tokens - requested
    allowed = 1
    
    -- Save new state
    redis.call("hmset", key, "tokens", filled_tokens, "last_refill", now)
    -- Expire key after 60s of inactivity to save RAM
    redis.call("expire", key, 60)
end

return allowed