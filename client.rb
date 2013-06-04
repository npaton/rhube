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


t = Time.now
1000.times do
	RedisRef.get "toto"
	RedisRef.set "toto", "42\n42"
end
puts "#{Time.now-t}s"

t = Time.now
1000.times do
	# RedisRhube.info
	RedisRhube.get "toto"
	RedisRhube.set "toto", "42\n42"
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
