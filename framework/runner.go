package framework 

import (
	"fmt"
	"os"
)

// ユーザーにどうやってコマンドを定義させるかを決める
type CommandFunc func(args []string) error
type Command struct {
	Name string
	ExecFunc CommandFunc
}

// 登録されたコマンドを管理し、実行時に適切なものを選んで動かす仕組み
type Runner struct {
	cmds []Command
}
func NewRunner(inputList []Command) *Runner {
	return &Runner {
		cmds: inputList,
	}
}

// Runnerという構造にRunという機能を追加する
func (r *Runner) Run() error {
	args := os.Args
	if len(args) < 2 {
		return fmt.Errorf("コマンドを指定してください")
	}
	subCommand := args[1]
	for _, cmd := range r.cmds {
		if cmd.Name == subCommand {
			return cmd.ExecFunc(args[1:])
		}
	}
	return fmt.Errorf("そのコマンドは知りません: %s", subCommand)
}