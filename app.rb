# encoding: utf-8
# frozen_string_literal: true
require 'nokogiri'
require 'open-uri'
require 'kconv'
# require '/projects/kyuko/pass.rb'
require_relative './pass.rb'

class NoLectures
  def initialize(today, place)
    # 1 => mon, 2 => tue,,,
    @today = today
    # 1 => imadegawa, 2 => tanabe
    @place = place
    @mon = Array.new(8, nil)
    @tue = Array.new(8, nil)
    @wed = Array.new(8, nil)
    @thu = Array.new(8, nil)
    @fri = Array.new(8, nil)
    @sat = Array.new(8, nil)
    @sun = Array.new(8, nil)
    # @no_lec = {:mon => @mon, :tue => @tue, :wed => @wed, :thu => @thu, :fri => @fri, :sat => @sat}
    @no_lec = [nil, @mon, @tue, @wed, @thu, @fri, @sat, @sun]

    @url = "http://duet.doshisha.ac.jp/info/KK1000.jsp?katei=1&youbi=#{@today}&kouchi=#{@place}"
  end

  def set_url(url)
    @url = url
  end

  def set_today(youbi)
    @today = youbi
    @url = "http://duet.doshisha.ac.jp/info/KK1000.jsp?katei=1&youbi=#{@today}&kouchi=#{@place}"
  end

  def show_nolec
    @no_lec
  end

  def show_today
    @today
   end

  def change_youbi_int(arg)
    i_to_youbi = { 1 => '月', 2 => '火', 3 => '水', 4 => '木', 5 => '金', 6 => '土', 7 => '日' }
    youbi_to_i = { 'Mon' => 1, 'Tue' => 2, 'Wed' => 3, 'Thu' => 4, 'Fri' => 5, 'Sat' => 6, 'Sun' => 7 }
    if arg.is_a?(Integer)
      return i_to_youbi[arg]
    elsif arg.is_a?(String)
      return youbi_to_i[arg]
    else
      puts 'error'
    end
  end

  def tomorrow(today)
    @today = today
    @today += 1
    @today = 1 if @today == 7
    set_url("http://duet.doshisha.ac.jp/info/KK1000.jsp?katei=1&youbi=#{@today}&kouchi=#{@place}")
    @today
  end

  def xml_to_text(xml)
    xml = xml.split('>')[1]
    xml = xml.split('<')[0]
    xml = '' if xml.nil?
    xml
  end

  def crawl_today
    # charset = open(@url).charset
    nangen = 0

    # doc = Nokogiri::HTML.parse(open(@url), nil, charset)
    doc = Nokogiri::HTML.parse(open(@url), nil, 'shift_jis')

    return 0 if doc.css('.styleE').inner_text.include?('該当する休講はありません')

    subjects = doc.css('.style1').each do |node|
      if node.children.inner_text.include?('講時')
        @array = []
        nangen = node.children.css('th').inner_text.delete('講時').to_i

        sub_name = xml_to_text(node.children.css('td').to_s.split("\n")[0])
        lecturer = xml_to_text(node.children.css('td').to_s.split("\n")[1])
        lecturer = lecturer.delete(' ')
        reason = xml_to_text(node.children.css('td').to_s.split("\n")[2])
        reason = reason.split('&')[0]

        @array << { sub_name: sub_name.toutf8, lecturer: lecturer.toutf8, reason: reason.toutf8 }

      else
        sub_name = xml_to_text(node.children.css('td').to_s.split("\n")[0])
        lecturer = xml_to_text(node.children.css('td').to_s.split("\n")[1])
        lecturer = lecturer.delete(' ')
        reason = xml_to_text(node.children.css('td').to_s.split("\n")[2])
        reason = reason.split('&')[0]

        @array << { sub_name: sub_name.toutf8, lecturer: lecturer.toutf8, reason: reason.toutf8 }
      end
      @no_lec[@today][nangen] = @array

      return @no_lec[@today][nangen].length
    end
  end

  def crawl_week
    6.times do |_i|
      crawl_today
      tomorrow(@today)
    end
  end
end
