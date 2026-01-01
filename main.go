package main

import (
	"fmt"
	"os"

	"simple-cli/framework"
	"simple-cli/action"
)

// 実行部分
func main() {
	myCmds := []framework.Command {
		{
			Name: "resize",
			ExecFunc: action.ResizeImage,
		},
		{
			Name: "bulkResize",
			ExecFunc: action.BulkResizeImage,
		},
		{
			Name: "copy",
			ExecFunc: action.CopyFile,
		},
		{
			Name: "benchmark",
			ExecFunc: action.Benchmark,
		},
	}
	r := framework.NewRunner(myCmds)
	if err := r.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}