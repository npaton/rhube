require "redis"

RedisRef = Redis.new(:port => 6379)
RedisRhube = Redis.new(:port => 6380)

# RedisRef.subscribe('toto', 'ruby-lang') do |on|
#   on.message do |channel, msg|
#     puts "##{channel} - #{msg}"
#   end
# end

# RedisRhube.subscribe('rubyonrails', 'ruby-lang') do |on|
#   on.message do |channel, msg|
#     puts "##{channel} - #{msg}"
#   end
# end

# puts RedisRef.publish "toto", 42
# puts RedisRhube.publish "toto", 42
# puts RedisRhube.set("toto", 42).inspect
# puts RedisRhube.get("toto").inspect
# exit
# puts RedisRhube.get("totoo")
# RedisRhube.set("totoo", "froujfd")
# puts RedisRhube.get("totoo")

# val = "a"*1048575
# val = "a"*1048575
val = "a"*100
times = 100000 #000

t = Time.now
times.times do |i|
	# RedisRef.set "test:#{i}", val
	RedisRef.sadd "test:set:a", "val#{i}"
	# raise if RedisRef.get("test:#{i}") != val
	# RedisRef.expire "test:#{i}", 2
end
puts "#{Time.now-t}s"

t = Time.now
times.times do |i|
	# RedisRhube.set "test:#{i}", val
	RedisRhube.sadd "test:set:a", "val#{i}"
	# raise if RedisRhube.get("test:#{i}") != val
	# RedisRhube.expire "test:#{i}", 2
end
puts "#{Time.now-t}s"


# command:["subscribe", "toto", "ruby-lang"]
# line:*3
# line:$9
# line:$4
# line::1
# line:*3
# line:$9
# line:$9
# line::2
# line:*3
# line:$7
# line:$4
# line:$2
# 
# ^Ccommand:[:unsubscribe]

# command:[:publish, "toto", 42]
# line::1
# 1
