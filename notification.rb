require 'twitter'

module Notification
  module Twitter
    module Config
      def self.token(campus)
        case campus
        when :imadegawa
          prefix = 'IMADEGAWA'
        when :kyotanabe
          prefix = 'KYOTANABE'
        else
          raise 'invalid campus'
        end

        token = {
          consumer_key: ENV["TWITTER_#{prefix}_CONSUMER_KEY"],
          consumer_secret: ENV["TWITTER_#{prefix}_CONSUMER_SECRET"],
          access_token: ENV["TWITTER_#{prefix}_ACCESS_TOKEN"],
          access_token_secret: ENV["TWITTER_#{prefix}_ACCESS_TOKEN_SECRET"]
        }

        raise "TWITTER_#{prefix}_* has nil" if token.value? nil

        token
      end
    end

    module Tweet
      @queue = :tweet

      def self.perform(campus, message)
        token = Notification::Twitter::Config.token(campus)
        twitter = Twitter::REST::Client.new(token)
        twitter.update message
        nil
      end
    end
  end
end
