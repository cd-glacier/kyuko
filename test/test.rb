require_relative "../app.rb"
require_relative "../tweet.rb"
require_relative "../pass.rb"

describe 'app' do

  context "change_youbi_int_test" do
    before do
      @nolec = NoLectures.new("mon", 1)
    end

    it '1 => 月' do
      expect(@nolec.change_youbi_int(1)).to eq "月"
    end

    it '6 => 土' do
      expect(@nolec.change_youbi_int(6)).to eq "土"
    end

    it 'Mon => 1' do
      expect(@nolec.change_youbi_int("Mon")).to eq 1
    end

    it 'Sat => 6' do
      expect(@nolec.change_youbi_int("Sat")).to eq 6
    end
  end


  context 'tomorrow_test' do
    before do
      @nolec = NoLectures.new("mon", 1)
    end

    it 'if mon, tue' do 
      expect(@nolec.tomorrow(1)).to eq 2
    end

    it 'if sat, mon' do 
      expect(@nolec.tomorrow(6)).to eq 1
    end
  end

  context 'xml to text' do
    before do
      @nolec = NoLectures.new("mon", 1)
    end

    it "<hoge>foo</hoge> to foo" do
      expect(@nolec.xml_to_text("<hoge>foo</hoge>")).to eq "foo"
    end
  end

  context "crawl test" do
    before do
      @nolec = NoLectures.new("mon", 1)
    end

    it "休講がないとき" do
      #kyuko_0.html参照して
      @nolec.set_url("./kyuko_0.html")
      expect(@nolec.crawl_today).to eq 0
    end
  end

end

describe "tweet" do
  before do
    @tw = Tweet.new(T_CONSUMER_KEY, T_CONSUMER_SECRET, T_ACCESS_TOKEN, T_ACCESS_TOKEN_SECRET, 2)
    @nolec = NoLectures.new("Mon", 2)
  end

  context "set_tomrrow_test" do
    it 'set_tomorrow_test' do
      now = Time.now
      youbi = now.strftime("%a")
      expect(@tw.set_tomorrow).to eq @nolec.change_youbi_int(youbi)
    end 
  end

  context "create_contents_test" do
    it "create_content_test" do

    end
  end

end






















