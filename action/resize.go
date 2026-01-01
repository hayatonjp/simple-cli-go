package action

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"flag"

	"github.com/nfnt/resize"
)

const OutputDir = "output" // 出力先フォルダ名

func BulkResizeImage(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("使い方: resize <幅> <ファイル1> <ファイル2> ...")
	}
	fs := flag.NewFlagSet(args[1], flag.ExitOnError)
	widthPtr := fs.Int("w", 0, "リサイズ後の幅(px)") // -wという名前でデフォルト値, 説明文
	fs.Parse(args[1:]) // resize以降の引数を読み取り
	if *widthPtr == 0 {
		return fmt.Errorf("幅を指定してください(例: -w 300)")
	}

	files := fs.Args()
	if len(files) == 0 {
		return fmt.Errorf("ファイル名を指定してください")
	}
	// ファイル名を取得
	width := *widthPtr
	widthStr := strconv.Itoa(width)
	fmt.Printf("[開始]幅: %d px, 対象ファイル数: %d 件\n", width, len(files))

	// outputフォルダの作成
	if err := os.MkdirAll(OutputDir, 0755); err != nil {
		return fmt.Errorf("出力フォルダの生成に失敗: %w", err)
	}

	var wg sync.WaitGroup
	limit := make(chan struct{}, 5) // 同時に5ファイルまで

	for _, file := range files {
		// すでに"_resized"がついているファイルは処理対象から外す
		if strings.Contains(file, "_resized") {
			fmt.Printf("[スキップ]生成済みファイルです: %s\n", file)
			continue
		}

		wg.Add(1) // +1

		// 並行処理
		go func(f string) {
			defer wg.Done() // -1

			limit <- struct{}{}
			defer func() { <- limit }()

			dummyArgs := []string{"dummy", f, widthStr}
			if err := ResizeImage(dummyArgs); err != nil {
				fmt.Printf("[error]%s: %v\n", f, err) // 並行処理中のエラーは画面に表示するだけ
			}
		}(file)
	}

	wg.Wait()

	fmt.Println("画像リサイズ処理が完了しました")
	return nil
}

func ResizeImage(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("ファイル名が指定されていません");
	}
	filename := args[1]
	if len(args) < 3 {
		return fmt.Errorf("リサイズしたい幅(px)が指定されていません")
	}
	widthStr := args[2]
	fmt.Printf("[画像処理]%sをリサイズします\n", filename)

	// string->int変換
	width, err := strconv.Atoi(widthStr);
	if err != nil {
		return fmt.Errorf("幅は数字で指定してください: %s", widthStr)
	}

	// ファイルを開く
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("ファイルが開けません: %w", err)
	}
	/**
	 * ファイルを開きすぎてエラーになる（ファイルディスクリプタの枯渇）問題への対処
	 * ここで書いておかないと、以降の処理でエラー起きた時に都度定義する必要があるため、書いておく
	 */
	defer file.Close()

	// 画像を読み込む
	img, format, err := image.Decode(file) // image.Decodeで拡張子を自動判別して画像データに変換
	if err != nil {
		return fmt.Errorf("画像として読み込めませんでした: %w", err)
	}
	fmt.Printf("画像を読み込みました(形式: %s)\n", format)

	// リサイズ実行
	fmt.Printf("幅 %d pxにリサイズ中...\n", width)
	newImg := resize.Resize(uint(width), 0, img, resize.Lanczos3)
	filenameOnly := filepath.Base(filename)

	// 保存するファイル名の作成
	ext := filepath.Ext(filenameOnly) // ".jpg"
	base := strings.TrimSuffix(filenameOnly, ext) // "photo"
	newFilename := getUniqueFileName(OutputDir, base, ext)

	// 書き込み用ファイルの作成
	outFile, err := os.Create(newFilename)
	if err != nil {
		return fmt.Errorf("保存用ファイルが作れません: %w", err)
	}
	defer outFile.Close()

	// 画像形式に合わせて保存(エンコード)
	switch format {
	case "jpeg":
		err = jpeg.Encode(outFile, newImg, &jpeg.Options{ Quality: 80 })
	case "png":
		err = png.Encode(outFile, newImg)
	default:
		return fmt.Errorf("未対応の保存形式です: %s", format)
	}

	if err != nil {
		return fmt.Errorf("保存に失敗しました: %w", err)
	}

	fmt.Printf("保存しました: %s\n", newFilename)
	return nil
}

// 同名のファイル名がある場合、_resized(1), _resized(2)...と番号を振って返す
func getUniqueFileName(dir, base, ext string) string {
	name := filepath.Join(dir, base + "_resized" + ext)
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return name
	}

	for i := 1; ; i++ {
		fname := fmt.Sprintf("%s_resized(%d)%s",base,i,ext)
		fullPath := filepath.Join(dir, fname)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fullPath
		}
	}
}