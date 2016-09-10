# encoding: utf-8
require '/projects/kyuko/pass.rb'
# require "./pass.rb"
require 'twitter'
require 'clockwork'
require 'date'
require '/projects/kyuko/app.rb'
# require "./app.rb"

class Tweet
  def initialize(c_key, c_secret, a_token, a_token_secret, place)
    @consumer_key = ''
    @consemer_secret = ''
    @access_token = ''
    @access_token_secret = ''

    @client = Twitter::REST::Client.new(
      consumer_key:        c_key,
      consumer_secret:     c_secret,
      access_token:        a_token,
      access_token_secret: a_token_secret
    )

    @nolec = NoLectures.new(@youbi, place)
    @date = DateTime.now
    @youbi = @nolec.change_youbi_int(@date.strftime('%a'))
    @nolec.set_today(@youbi)
    @contents = []
  end

  def set_time
    @date = DateTime.now
   end

  def set_tomorrow
    # 今何時か調べて、21よりあとなら明日の情報
    set_time
    hour = @date.strftime('%H').to_i
    @nolec.tomorrow(@nolec.show_today) if hour >= 21
  end

  def create_contents
    @nolec.crawl_today
    youbi_name = @nolec.change_youbi_int(@nolec.show_today)
    nolec = @nolec.show_nolec

    content = "#{youbi_name}曜日の休講情報\n#{@date.strftime("%H時%M分")}時点\n"

    i = 0
    nil_counter = 0
    nolec[@nolec.show_today].each do |ttable|
      unless i == 0
        if ttable.nil?
          nil_counter += 1
        else
          ttable.each do |sub_info|
            nangen = i
            sub_name = sub_info[:sub_name]
            lecturer = sub_info[:lecturer]
            reason = sub_info[:reason]

            content << "#{nangen}限目:#{sub_name} 講師(#{lecturer})\n"

            unless content[90].nil? then
              @contents << content
              content = "#{youbi_name}曜日の休講情報\n#{@date.strftime("%H時%M分")}時点\n"
            end
          end
        end
      end
      i += 1
    end

    content = "#{youbi_name}曜日の休講はありません" if nil_counter == 7

    @contents << content
    content = nil
  end

  def update_tweet
    @contents.each do |content|
      puts content
      @client.update(content)
    end
  end
end

include Clockwork

every(4.hours, "work") do

	puts "田辺"
	#田辺
  #今日の曜日をset	

  tw_tanabe = Tweet.new(T_CONSUMER_KEY, T_CONSUMER_SECRET, T_ACCESS_TOKEN, T_ACCESS_TOKEN_SECRET, 2)
  # 今何時か調べて、21よりあとなら明日の情報
  tw_tanabe.set_tomorrow
  tw_tanabe.create_contents
  tw_tanabe.update_tweet

  puts "今出川"
  # 今出川
  # 今日の曜日をset
  tw_imade = Tweet.new(I_CONSUMER_KEY, I_CONSUMER_SECRET, I_ACCESS_TOKEN, I_ACCESS_TOKEN_SECRET, 1)
  # 今何時か調べて、21よりあとなら明日の情報
  tw_imade.set_tomorrow
  tw_imade.create_contents
  tw_imade.update_tweet

  puts 'end'
end
