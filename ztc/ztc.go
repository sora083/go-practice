package main

/// Zero Time Cacheの簡易実装
/// 重い処理 (Update()関数の後半部分で行われる処理)を無駄に実行しないように制御するためのパターン
// Ref: https://gist.github.com/tukeJonny/db9aa5fdd066e095904aae0880f913fd

import (
	"log"
	"strings"
	"sync"
	"time"
)

var (
	isDone      = false
	DoneChannel = make(chan bool)
	keywords    []string

	wg = &sync.WaitGroup{}

	KeywordLock       sync.Mutex // キーワード追加処理用
	ZeroTimeCacheLock sync.Mutex // Zero Time Cache用
	LastUpdated       time.Time  // 最終更新時刻
)

func AddNewKeyword(keyword string) {
	wg.Add(1)
	for { // 同時呼び出し
		if isDone {
			break
		}
	}

	KeywordLock.Lock()
	log.Printf("BEFORE append %s\n", keyword)
	keywords = append(keywords, keyword) // ここで新しいキーワードを追加
	log.Printf("AFTER append %s\n", keyword)
	KeywordLock.Unlock()

	Update() // 更新処理を呼び出す
	wg.Done()
}

func Update() {
	beforeLockTime := time.Now()

	ZeroTimeCacheLock.Lock()
	defer ZeroTimeCacheLock.Unlock()

	if LastUpdated.After(beforeLockTime) { //  最後に更新した時刻(厳密には、最後に更新を開始した時刻)は、ロック取得前の時刻(キーワード更新後の時刻)よりも後であるか
		log.Printf("%v After %v\n", LastUpdated, beforeLockTime)
		return //  更新の必要はない
	}
	LastUpdated = time.Now() // 最終更新時刻を更新

	log.Printf("[UPDATE] %s\n", strings.Join(keywords, ","))
}

// 同時実行のための関数
func SimultaneouslyExecutor() {
	msg := <-DoneChannel // チャネルから値を受信
	log.Printf("GET message %v\n", msg)
	isDone = true
}

func main() {
	// A,B,C,D,E,F(順不同)が出力されるのが１回であれば、成功(バッチ的に処理できている)

	go AddNewKeyword("A")
	go AddNewKeyword("B")
	go AddNewKeyword("C")
	go AddNewKeyword("D")
	go AddNewKeyword("E")
	go AddNewKeyword("F")

	go SimultaneouslyExecutor()

	// Notify
	DoneChannel <- true // チャネルへ値を送信

	// Wait
	wg.Wait()

	log.Println("DONE")
}
