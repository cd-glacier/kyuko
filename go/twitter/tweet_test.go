package twitter

import "testing"

func testUpdate(t *testing.T) {
	err := Update("test")
	if err != nil {
		t.Fatalf("tweetに失敗しました\nerr:%s", err)
	}
}
