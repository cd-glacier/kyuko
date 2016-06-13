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

client.update("test") 
















