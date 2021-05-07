// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"stem-cell/meta"
	"stem-cell/util"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var (
	createCmd = &cobra.Command{
		Use:   "new",
		Short: "在当前目录下创建一个基于takumi的空项目",
		Long: `For example:

stem-cell new [项目名] [项目存放路径]

# 随机生成[40000,50000]之间的服务监听端口，如与现有服务冲突请自行修改。`,
		Run:  createTakumi,
		Args: cobra.MinimumNArgs(2),
	}

	arg              = meta.Params{}
	silent           bool
	randomListenPort int
	placedPath       string
	namePattern      = regexp.MustCompile("^[0-9a-zA-Z_./-]+$")
	groupPattern     = regexp.MustCompile("^[0-9a-zA-Z_.-]+$")

	//go:embed tmpl
	embededTmpl embed.FS
)

func init() {
	createCmd.Flags().StringVarP(&arg.Org, "organization", "o", "amtcloud.cn", "公司名")
	createCmd.Flags().StringVarP(&arg.ProjectGroup, "group", "g", "", "业务分组")
	createCmd.Flags().StringVarP(&arg.ProjectName, "name", "n", "", "服务名")
	createCmd.Flags().StringVarP(&arg.Desc, "desc", "d", "to be or not to be", "服务简介")
	createCmd.Flags().IntVarP(&arg.Port, "port", "p", 0, "服务端口号")
	createCmd.Flags().BoolVarP(&silent, "silent", "s", false, "gen code silently")

	RootCmd.AddCommand(createCmd)

	rand.Seed(time.Now().Unix())
	portOffset := rand.Intn(10000)
	randomListenPort = 40000 + portOffset
}

func createTakumi(cmd *cobra.Command, args []string) {
	arg.ProjectName = args[0]
	placedPath = args[1]
	if !silent {
		readArgFromConsole()
	}

	arg.CamelProjectName = util.Camel(arg.ProjectName)
	if arg.Port == 0 {
		arg.Port = randomListenPort
	}

	placedPath, err := filepath.Abs(placedPath)
	if err != nil {
		panic(err)
	}

	placedPath = path.Join(placedPath, arg.ProjectGroup, arg.ProjectName)
	fmt.Printf("as your wish, gen project: \n%#v\n", arg)

	embedroot, err := fs.Sub(embededTmpl, "tmpl")
	if err != nil {
		panic(err)
	}

	// 遍历模板目录
	err = fs.WalkDir(embedroot, ".", func(fpath string, info fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", fpath, err)
			return err
		}

		targetPath := path.Join(placedPath, fpath)
		if info.IsDir() {
			//fmt.Printf("mkdir: %+v \n", info.Name())
			err = util.Makedir(err, targetPath)
			if err != nil {
				return err
			}
			return nil
		}

		//生成代码
		if err = render(fpath, targetPath); err != nil {
			panic(err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("")
	exec.Command("sh", "-c", fmt.Sprintf(`ln -s ./conf_dev.yml %s/conf.yml`, placedPath)).Output()
	out, _ := exec.Command("sh", "-c", fmt.Sprintf(`tree %s`, placedPath)).Output()
	fmt.Printf("generate new project[%s] at:\n", arg.ProjectName)
	fmt.Println(string(out))
}

func readArgFromConsole() {
	reader := bufio.NewReader(os.Stdin)
	// 在交互模式下读取工作参数
	scan := func(hint string, defval string, accept func(string) bool) string {
		for {
			fmt.Printf("请输入%s, 默认：%s >", hint, defval)
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1) // convert CRLF to LF
			if text == "" {
				if defval == "" {
					continue
				}
				text = defval
			}
			if !accept(text) {
				fmt.Println("输入错误!")
				continue
			}
			return text
		}
	}

	arg.Org = scan("公司名称", arg.Org, func(t string) bool { return groupPattern.MatchString(t) })
	arg.ProjectGroup = scan("业务分组名称", arg.ProjectGroup, func(t string) bool { return groupPattern.MatchString(t) })
	arg.ProjectName = scan("项目名", arg.ProjectName, func(t string) bool { return groupPattern.MatchString(t) })
	arg.Desc = scan("简介", arg.Desc, func(t string) bool { return t != "" })
	tPort := scan("端口", strconv.Itoa(randomListenPort), func(t string) bool { _, b := util.CheckInt(t); return b })
	var err error
	arg.Port, err = strconv.Atoi(tPort)
	if err != nil {
		panic(err)
	}
}

func render(tmplPath string, targetPath string) error {
	embedTmplPath := path.Join("tmpl", tmplPath)
	tmpl, err := template.ParseFS(embededTmpl, embedTmplPath)
	if err != nil {
		fmt.Println("template.New failed:", err)
		return err
	}

	switch filepath.Base(targetPath) {
	case "gitignore":
		fallthrough
	case "drone.yml":
		dir, file := filepath.Split(targetPath)
		targetPath = path.Join(dir, "."+file)
	default:
		if filepath.Ext(targetPath) == ".tmpl" {
			targetPath = targetPath[:len(targetPath)-5]
		}
	}

	f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("os.OpenFile faild: ", err)
		return err
	}
	tmpl.Execute(f, arg)
	f.Close()

	return nil
}
