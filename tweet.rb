# encoding: utf-8
require "/projects/kyuko/pass.rb"
require "twitter"
require "clockwork"
require "date"
require "/projects/kyuko/app.rb"

client = Twitter::REST::Client.new(
	consumer_key:        CONSUMER_KEY,
	consumer_secret:     CONSUMER_SECRET,
	access_token:        ACCESS_TOKEN,
	access_token_secret: ACCESS_TOKEN_SECRET,
)

include Clockwork

every(2.hours, "work") do
	tanabe = 2
	youbi = 1
	nolec = NoLectures.new(youbi, tanabe)   

	#今日の曜日をset	
	date = DateTime.now
	youbi = nolec.change_youbi_int(date.strftime("%a"))
	nolec.set_today(youbi)     

	#今何時か調べて、21よりあとなら明日の情報	
	hour = date.strftime("%H").to_i
	if hour > 20 then
		nolec.tomorrow()
	end

	nolec.crawl_today()
	youbi_name = nolec.change_youbi_int(youbi)        
	nolec = nolec.show_nolec

	contents = []
	content = "#{youbi_name}曜日の休講情報\n#{date.strftime("%H時%M分")}時点\n" 

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
					
					unless content[100].nil? then
						contents << content
						content = "#{youbi_name}曜日の休講情報\n#{date.strftime("%H時%M分")}時点\n" 
					end
				end
			else
				#content = "#{youbi_name}曜日の休講はありません"
			end
		end
		i = i + 1
	end
	contents << content

	contents.each do |content|
		#puts content
	 	client.update(content) 
	end
	puts "end"
end

















