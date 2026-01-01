package framework

import (
	"errors"
	"os"
	"testing"
)

func TestRunner_Run(t *testing.T) {
	// テスト用のコマンド定義
	called := false // 実行されたかのフラグ
	mockCmds := []Command {
		{
			Name: "test-cmd",
			ExecFunc: func(args []string) error {
				called = true
				if len(args) > 0 && args[0] == "error" {
					return errors.New("意図的なエラー")
				}
				return nil
			},
		},
	}

	// Runner作成
	runner := NewRunner(mockCmds)

	// 正常系
	origArgs := os.Args
	defer func() { os.Args = origArgs }() // プログラムが終わる時にos.Argsを元に戻しておく

	os.Args = []string{"mytool", "test-cmd", "arg1"} // 偽物のコマンド引数に変える
	if err := runner.Run(); err != nil {
		t.Errorf("エラー: %v", err)
	}
	if !called {
		t.Error("コマンドが実行されていません")
	}

	// 異常系: 知らないコマンド
	os.Args = []string{"mytool", "unknown-cmd"}
	if err := runner.Run(); err == nil {
		t.Error("知らないコマンドなのにエラーになりませんでした")
	}

	// 異常系: 引数が足りない
	os.Args = []string{"mytool"}
	if err := runner.Run(); err == nil { // :=で定義+代入
		t.Error("引数が足りないのにエラーになりませんでした")
	}
}