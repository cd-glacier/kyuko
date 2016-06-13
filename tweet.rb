require "./pass.rb"
require "twitter"
require "clockwork"
require "./app.rb"

client = Twitter::REST::Client.new(
  consumer_key:        CONSUMER_KEY,
  consumer_secret:     CONSUMER_SECRET,
  access_token:        ACCESS_TOKEN,
  access_token_secret: ACCESS_TOKEN_SECRET,
)

def limit_140(arg)
  return arg.scan(/.{1,137}/m)[0]
end


include Clockwork

every(1.hour, "work") do
  tanabe = 2
  youbi = 2
  nolec = NoLectures.new(youbi, tanabe)   
  nolec.crawl_today()
  youbi_name = nolec.change_youbi_int(youbi)        
  nolec = nolec.show_nolec

  content = "#{youbi_name}曜日の休講情報\n" 

  i = 0
  nolec[youbi].each do |ttable|
    unless i == 0 then 
      unless ttable.nil? then
        ttable.each do |sub_info|
          nangen = i
          sub_name = sub_info[:sub_name] 
          lecturer = sub_info[:lecturer]
          reason = sub_info[:reason]

          content << "#{nangen}限目:#{sub_name} 講師(#{lecturer})\n"

        end
      end
    end
    i = i + 1
  end
  content = limit_140(content)
  content << "..."
  client.update(content)  
end

















