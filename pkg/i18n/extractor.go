/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - Go 开发框架 (BlueKing - Go Framework) available.
 * Copyright (C) 2017 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 *	https://opensource.org/licenses/MIT
 *
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package i18n

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"bk.tencent.com/blueapps-go/pkg/config"
)

// 匹配 Golang 模板中使用 i18n 的正则表达式
var tmplI18nRegex = regexp.MustCompile(".*?i18n.*?\"(.*?)\" \\.lang")

// ExtractMessages 提取国际化消息
func ExtractMessages() error {
	extractor := &msgExtractor{
		baseDir:      config.BaseDir,
		msgFilepath:  MsgFilepath(),
		msgLocations: map[string][]string{},
	}
	return extractor.Exec()
}

// 国际化消息提取器
type msgExtractor struct {
	baseDir     string
	msgFilepath string
	// 消息位置映射：{消息 ID: 消息位置列表}
	// 消息位置示例：templates/web/401.html:27
	msgLocations map[string][]string
}

// Exec 执行国际化消息提取流程
func (e *msgExtractor) Exec() error {
	// 1. 从源代码 & 模板中提取国际化消息
	extractedMsgs, err := e.extract()
	if err != nil {
		return errors.Wrapf(err, "failed to extract messages")
	}
	// 2. 从已有的消息文件中读取存量配置
	existingMsgs, err := e.read()
	if err != nil {
		return errors.Wrapf(err, "failed to read messages from file: %s", e.msgFilepath)
	}
	// 3. 合并存量配置与提取配置
	mergedMsgs := e.merge(extractedMsgs, existingMsgs)
	// 4. 将合并后的配置写入到消息文件中
	if err = e.write(mergedMsgs); err != nil {
		return errors.Wrapf(err, "failed to write messages to file: %s", e.msgFilepath)
	}
	return nil
}

// 从源代码 & 模板中提取国际化消息
func (e *msgExtractor) extract() ([]msg, error) {
	var messages []msg
	walkErr := filepath.Walk(e.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		// 只处理 Go 源代码 & 模板文件
		if ext == ".go" {
			// Go 源代码通过解析语法树提取
			msgs, aErr := e.astAnalyze(path)
			if aErr != nil {
				return errors.Wrapf(aErr, "failed to ast match file: %s", path)
			}
			messages = append(messages, msgs...)
		} else if slices.Contains([]string{".html", ".yaml", ".tpl"}, ext) {
			// 模板文件通过正则表达式提取（仅限单行）
			msgs, rErr := e.regexMatch(path)
			if rErr != nil {
				return errors.Wrapf(rErr, "failed to regex match file: %s", path)
			}
			messages = append(messages, msgs...)
		}
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}
	return messages, nil
}

// 记录国际化消息的位置
func (e *msgExtractor) storeLocation(msgID, filepath string, lineNum int) {
	if e.msgLocations == nil {
		e.msgLocations = make(map[string][]string)
	}
	location := fmt.Sprintf("%s:%d", strings.TrimPrefix(filepath, e.baseDir+"/"), lineNum)
	e.msgLocations[msgID] = append(e.msgLocations[msgID], location)
}

// 通过语法树解析文件内容，提取国际化消息（只处理 Go 源代码）
func (e *msgExtractor) astAnalyze(filepath string) ([]msg, error) {
	fSet := token.NewFileSet()
	// 解析文件
	file, err := parser.ParseFile(fSet, filepath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	// 遍历语法树，提取国际化消息
	var messages []msg
	ast.Inspect(file, func(node ast.Node) bool {
		// 检查是否为调用表达式类型
		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		// 检查调用表达式类型 Fun 是否为选择器表达式
		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		// 检查 i18n. 调用（i18n package 中直调 T 不会被提取）
		ident, ok := selectorExpr.X.(*ast.Ident)
		if !ok || ident.Name != "i18n" {
			return true
		}
		// 检查 T 调用（i18n.T)
		if selectorExpr.Sel.Name != "T" {
			return true
		}
		// 至少需要包含两个参数（ctx，msgID）
		args := callExpr.Args
		if len(args) < 2 {
			return true
		}
		// 第二个参数若为字符串类型，则为 msgID
		arg, ok := args[1].(*ast.BasicLit)
		if !ok {
			return true
		}
		// 去除前后的引号
		msgID := strings.Trim(arg.Value, "\"")
		// 提取国际化消息内容
		messages = append(messages, newMsgWithPlaceholders(msgID))
		// 记录消息位置
		e.storeLocation(msgID, filepath, fSet.Position(arg.Pos()).Line)
		return true
	})
	return messages, nil
}

// 正则匹配文件内容，提取国际化消息（只匹配模板的）
func (e *msgExtractor) regexMatch(filepath string) ([]msg, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var messages []msg
	// 逐行扫描文件内容，匹配国际化消息
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	curLineNum := 0
	for scanner.Scan() {
		curLineNum++
		text := scanner.Text()
		// 跳过包含 I18nRegex 的行
		if strings.Contains(text, "I18nRegex") {
			continue
		}
		// 对每一行的每个匹配都要记录
		for _, matches := range tmplI18nRegex.FindAllStringSubmatch(text, -1) {
			if len(matches) > 1 {
				msgID := matches[1]
				messages = append(messages, newMsgWithPlaceholders(msgID))
				e.storeLocation(msgID, filepath, curLineNum)
			}
		}
	}
	return messages, nil
}

// 从已有的消息文件中读取存量配置
func (e *msgExtractor) read() ([]msg, error) {
	if _, err := os.Stat(e.msgFilepath); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "msg file not found: %s", e.msgFilepath)
	}

	data, err := os.ReadFile(e.msgFilepath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read msg file: %s", e.msgFilepath)
	}

	var messages []msg
	if err = yaml.Unmarshal(data, &messages); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal msg file: %s", e.msgFilepath)
	}
	return messages, nil
}

// 合并存量配置与提取配置
func (e *msgExtractor) merge(extracted []msg, existing []msg) []msg {
	msgMap := map[string]msg{}
	// 优先填入提取出的消息
	for _, m := range extracted {
		msgMap[m.ID] = m
	}
	// 合并已存在的消息（忽略已经删除的）
	for _, m := range existing {
		if _, ok := msgMap[m.ID]; ok {
			msgMap[m.ID] = m
		}
	}
	msgs := lo.Values(msgMap)
	// 填充 Locations（文件遍历时有序添加，因此无需重新排序）
	for idx, m := range msgs {
		msgs[idx].Locations = e.msgLocations[m.ID]
	}
	// 按照 ID 排序
	slices.SortFunc(msgs, func(m1, m2 msg) int {
		return strings.Compare(m1.ID, m2.ID)
	})
	return msgs
}

// 将合并后的配置写入到消息文件中
func (e *msgExtractor) write(msgs []msg) error {
	// 由于需要以注释的形式记录国际化消息位置，因此采用模板而非 yaml.Unmarshal
	tmpl := template.Must(template.New("msgFile").Funcs(sprig.FuncMap()).Parse(msgFileTpl))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, msgs); err != nil {
		return errors.Wrapf(err, "failed to execute msgFile template")
	}

	// 写入到消息文件
	if err := os.WriteFile(e.msgFilepath, buf.Bytes(), 0o644); err != nil {
		return errors.Wrapf(err, "failed to write messages to file: %s", e.msgFilepath)
	}
	return nil
}
