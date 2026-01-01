package action

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CopyFile(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("使い方: copy -n <回数> <ファイル名>")
	}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	countPtr := fs.Int("n", 1, "コピーする回数")
	fs.Parse(args[1:])

	remainingArgs := fs.Args()
	if len(remainingArgs) == 0 {
		return fmt.Errorf("コピー元のファイルを指定して下さい")
	}
	srcFilename := remainingArgs[0]
	count := *countPtr

	fmt.Printf("[開始]%sを%d回コピーします\n", srcFilename, count)

	// コピー元ファイルを開く
	srcFile, err := os.Open(srcFilename)
	if err != nil {
		return fmt.Errorf("ファイルが開けません: %w", err)
	}
	defer srcFile.Close()

	ext := filepath.Ext(srcFilename)
	base := strings.TrimSuffix(srcFilename, ext)

	for i := 1; i <= count; i++ {
		// 保存ファイル名
		dstFilename := fmt.Sprintf("%s(%d)%s",base,i,ext)

		// コピー先ファイル作成
		dstFile, err := os.Create(dstFilename)
		if err != nil {
			fmt.Printf("[error]作成失敗 (%s): %v\n", dstFilename, err)
			continue
		}

		// コピー元の読み込み位置を先頭に戻す。これをしないと空ファイルができる
		srcFile.Seek(0,0)

		// コピー
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			fmt.Printf("[error]failed copy %s: %v\n", dstFilename, err)
			dstFile.Close()
			continue
		}

		dstFile.Close() // 書き込み完了したら閉じる(deferだとループ終わるまで閉じられないため)
		fmt.Printf("作成: %s\n", dstFilename)
	}

	fmt.Println("コピー完了")
	return nil
}