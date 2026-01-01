## 使い方
### ベンチマーク実行(順次処理と並行処理での比較)
```go run main.go benchmark <画像フォルダのパス>```
### ファイルコピー(n回分)
```go run main.go copy -n <回数> <ファイル名>```
### リサイズ(単発)
```go run main.go resize <リサイズ後の幅> <ファイル名>```
### リサイズ(一括)
```go run main.go BulkResize -w <リサイズ後の幅> <ファイル名1> <ファイル名2> ... <ファイル名n>```
