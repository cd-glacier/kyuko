# encoding: utf-8
#require "/projects/kyuko/pass.rb"
require "./pass.rb"
require "twitter"
require "clockwork"
require "date"
require "./app.rb"

class Tweet 
  @@consumer_key = ""
  @@consemer_secret = ""
  @@access_token = ""
  @@access_token_secret = ""
  @@client
  @@nolec = NoLectures.new(0, 0)   
  @@date = DateTime.now
  @@youbi = @@nolec.change_youbi_int(@@date.strftime("%a"))
  @@contents = []

  def initialize(c_key, c_secret, a_token, a_token_secret, place)
    @@client = Twitter::REST::Client.new(
      consumer_key:        c_key,
      consumer_secret:     c_secret,
      access_token:        a_token,
      access_token_secret: a_token_secret,
    )

    @@nolec = NoLectures.new(@@youbi, place)   
    @@nolec.set_today(@@youbi)     
  end

  def set_tomorrow()
    #今何時か調べて、21よりあとなら明日の情報	
    hour = @@date.strftime("%H").to_i
    if hour >= 21 then
      @@nolec.tomorrow()
    end
  end

  def create_contents()
    @@nolec.crawl_today()
    youbi_name = @@nolec.change_youbi_int(@@youbi)        
    nolec = @@nolec.show_nolec

    content = "#{youbi_name}曜日の休講情報\n#{@@date.strftime("%H時%M分")}時点\n" 

    i = 0
    nolec[@@youbi].each do |ttable|
      unless i == 0 then 
        unless ttable.nil? then
          ttable.each do |sub_info|
            nangen = i
            sub_name = sub_info[:sub_name] 
            lecturer = sub_info[:lecturer]
            reason = sub_info[:reason]

            content << "#{nangen}限目:#{sub_name} 講師(#{lecturer})\n"

            unless content[100].nil? then
              @@contents << content
              content = "#{youbi_name}曜日の休講情報\n#{@@date.strftime("%H時%M分")}時点\n" 
            end
          end
        else
          #content = "#{youbi_name}曜日の休講はありません"
        end
      end
      i = i + 1
    end
    @@contents << content
  end

  def update_tweet()
    @@contents.each do |content|
      puts content
      #client.update(content) 
    end
  end


end


include Clockwork

every(2.hours, "work") do

  #今日の曜日をset	
  tw_tanabe = Tweet.new(CONSUMER_KEY, CONSUMER_SECRET, ACCESS_TOKEN, ACCESS_TOKEN_SECRET, 2)
  #今何時か調べて、21よりあとなら明日の情報	
  tw_tanabe.set_tomorrow()
  tw_tanabe.create_contents()
  tw_tanabe.update_tweet() 

  puts "end"
end

















