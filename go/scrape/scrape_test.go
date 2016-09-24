package scrape

import "testing"

func TestSetUrl(t *testing.T) {
	if url, err := SetUrl(1, 1); url != "http://duet.doshisha.ac.jp/info/KK1000.jsp?katei=1&youbi=1&kouchi=1" {
		t.Fatalf("urlの生成がうまくできていないようです\n err: %s", err)
	}

	if url, err := SetUrl(2, 5); url != "http://duet.doshisha.ac.jp/info/KK1000.jsp?katei=1&youbi=5&kouchi=2" || err != nil {
		t.Fatalf("urlの生成がうまくできていないようです\n err: %s", err)
	}

	if url, err := SetUrl(3, 1); err == nil {
		t.Fatalf("存在しない校地のurlが生成されています\n created url: %s", url)
	}

	if url, err := SetUrl(1, 7); err == nil {
		t.Fatalf("日曜日のurlは必要ありません\n created url: %s", url)
	}

}

//まだできてない
func TestScrape(t *testing.T) {
	noTestHTML := ""
	if r, err := Scrape(noTestHTML); r.Place == 0 || err != nil {
		t.Fatal("休講がない時のテストデータでエラー\n")
	}

	testHTML := ""
	if r, err := Scrape(testHTML); r.Place == 0 || err != nil {
		t.Fatal("HTMLから休講情報が取得できないようです\n")
	}
}
