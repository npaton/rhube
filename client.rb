require "redis"

RedisReference = Redis.new(:port => 6379)
RedisRhube = Redis.new(:port => 6380)

puts RedisReference.info
puts RedisRhube.info
