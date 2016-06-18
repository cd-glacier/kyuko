# encoding: utf-8
require 'twitter'
require 'clockwork'
require 'date'
require '/projects/kyuko/app.rb'
# require "./app.rb"

class Tweet
  def initialize(place)
    @client = Twitter::REST::Client.new(
      consumer_key:        ENV['TWITTER_CONSUMER_KEY'],
      consumer_secret:     ENV['TWITTER_CONSUMER_SECRET'],
      access_token:        ENV['TWITTER_ACCESS_TOKEN'],
      access_token_secret: ENV['TWITTER_ACCESS_TOKEN_SECRET']
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

            unless content[100].nil?
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

every(2.hours, 'work') do
  puts "田辺"
  # 田辺
  # 今日の曜日をset
  tw_tanabe = Tweet.new(2)
  # 今何時か調べて、21よりあとなら明日の情報
  tw_tanabe.set_tomorrow
  tw_tanabe.create_contents
  tw_tanabe.update_tweet

  puts "今出川"
  # 今出川
  # 今日の曜日をset
  tw_imade = Tweet.new(1)
  # 今何時か調べて、21よりあとなら明日の情報
  tw_imade.set_tomorrow
  tw_imade.create_contents
  tw_imade.update_tweet

  puts 'end'
end
