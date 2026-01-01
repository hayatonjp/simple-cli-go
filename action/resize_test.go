package action

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func createDummyImage(filename string) error {
	img := image.NewRGBA(image.Rect(0,0,10,10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x,y,color.RGBA{255,0,0,255})
		}
	}

	f, err := os.Create(filename);
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func TestResizeImage(t *testing.T) {
	// テスト用の画像を作成
	inputName := "test_input.png"
	outputName := "test_input_resized.png"

	defer os.Remove(inputName)
	defer os.Remove(outputName)

	if err := createDummyImage(inputName); err != nil {
		t.Fatalf("テスト画像の作成に失敗: %v", err)
	}

	// ResizeImage実行
	args := []string{"main.go", inputName, "5"}
	err := ResizeImage(args)

	// 検証
	if err != nil {
		t.Errorf("リサイズ処理でエラー: %v", err)
	}

	// 出力できているか
	if _, err := os.Stat(outputName); os.IsNotExist(err) {
		t.Errorf("リサイズ後のファイル%sが作られてません", outputName)
	}
}