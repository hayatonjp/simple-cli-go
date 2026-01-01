package action

import (
	"os"
	"strings"
	"strconv"
	"path/filepath"
	"fmt"
	"sync"
	"time"
)

func Benchmark(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("使い方: benchmark <画像フォルダのパス>")
	}
	targetDir := args[1]
	const width = 100
	widthStr := strconv.Itoa(width)
	files, err := collectImageFiles(targetDir)
	if err != nil {
		return fmt.Errorf("画像の取得に失敗: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("指定されたフォルダに画像ファイルが見つかりません")
	}

	fmt.Println("===== ベンチマーク開始 =====")

	defer func() {
		fmt.Println("\n[cleanup]outputフォルダを削除しています...")
		os.RemoveAll(OutputDir)
	}()

	// --- 実験1: 順次処理
	fmt.Println("\n[1]順次処理スタート")
	startSeq := time.Now()

	for _, f := range files {
		args := []string{"main.go", f, widthStr}
		_ = ResizeImage(args)
	}

	durationSeq := time.Since(startSeq)
	fmt.Printf(">> time: %v\n", durationSeq)


	// --- 実験2: 並行処理(5つ同時に実行)
	fmt.Println("\n[2]並行処理スタート")
	startPara := time.Now()

	var wg sync.WaitGroup
	limit := make(chan struct{}, 5)

	for _, f := range files {
		wg.Add(1)
		go func(target string) {
			defer wg.Done()
			limit <- struct{}{}
			defer func() { <- limit }()

			args := []string{"main.go", target, widthStr}
			_ = ResizeImage(args)
		}(f)
	}
	wg.Wait()

	durationPara := time.Since(startPara)
	fmt.Printf(">> time: %v\n", durationPara)

	// --- 結果 ---
	fmt.Println("\n=== 結果 ===")
	fmt.Printf("画像枚数: %d枚\n", len(files))
	fmt.Printf("順次処理: %v\n", durationSeq)
	fmt.Printf("並行処理: %v\n", durationPara)

	if durationPara > 0 { 
		ratio := float64(durationSeq) / float64(durationPara)
		fmt.Printf("speedup: 約 %.2f 倍\n", ratio)
	}
	return nil
}

// 指定フォルダ内のjpg,jpeg,pngを全て取得
func collectImageFiles(dir string) ([]string, error) {
	var files []string

	// フォルダの中身を読む
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		// ディレクトリは無視
		if entry.IsDir() {
			continue
		}

		// 拡張子チェック
		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))

		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			fullPath := filepath.Join(dir, name)
			files = append(files, fullPath)
		}
	}

	return files, nil
}