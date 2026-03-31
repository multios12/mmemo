package diary

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type listModel struct {
	WritedMonths []string // 記載された年月(yyyy-mm形式)
	Lines        []lineModel
}

func readListFile(month string) (*listModel, error) {
	p := filepath.Join(diaryPath, month+".txt")

	l := new(listModel)
	if _, err := os.Stat(p); err != nil {
		l.Lines = []lineModel{}
		return l, nil
	}

	fp, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		d, err := newLineModel(scanner.Text())
		if err != nil {
			return nil, err
		}
		l.Lines = append(l.Lines, *d)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return l, nil
}

func (l *listModel) writeListFile(month string) error {
	sort.Slice(l.Lines, func(i, j int) bool { return l.Lines[i].Day > l.Lines[j].Day })

	dataFile := ""
	for _, l := range l.Lines {
		if dataFile != "" {
			dataFile += "\n"
		}
		dataFile += l.toString()
	}
	filename := filepath.Join(diaryPath, month+".txt")

	if len(l.Lines) == 0 {
		return os.Remove(filename)
	}
	return os.WriteFile(filename, []byte(dataFile), os.ModePerm)
}

func (l *listModel) updateLine(month string, target lineModel) error {
	for i, line := range l.Lines {
		if target.Day == line.Day {
			l.Lines[i] = target
			return l.writeListFile(month)
		}
	}
	l.Lines = append(l.Lines, target)
	return l.writeListFile(month)
}

// ----------------------------------------------------------------------------

// 月ファイルの行を表すモデル
type lineModel struct {
	Day      string   // 日付(yyyy-mm-dd形式)
	Outline  string   // 概要
	Tags     []string // タグ
	IsDetail bool     // 詳細
	HCount   int      // メモ数
}

func newLineModel(s string) (*lineModel, error) {
	if len(s) < 11 {
		return nil, fmt.Errorf("invalid diary line: %q", s)
	}
	l := new(lineModel)
	l.Day = s[0:10]
	l.IsDetail = s[10:11] == "*"
	s = s[11:]

	tags := strings.Split(s, "#")
	if len(tags) > 1 {
		l.Outline = tags[len(tags)-1]
		l.Tags = tags[:len(tags)-1]
	} else {
		l.Outline = s
		l.Tags = []string{}
	}
	l.HCount = 0
	return l, nil
}

func (l *lineModel) toString() string {
	outline := strings.ReplaceAll(l.Outline, "#", "")
	line := l.Day
	if l.IsDetail {
		line += "*"
	} else {
		line += " "
	}

	if len(l.Tags) > 0 {
		line += strings.Join(l.Tags, "#") + "#"
	}
	line += outline
	return line
}

// 月ファイルの行を表すモデル
type detailModel struct {
	Day     string   // 日付(yyyy-mm-dd形式)
	Outline string   // 概要
	Tags    []string // タグ
	Detail  string   // 詳細
}

func readDetailFile(l lineModel) *detailModel {
	d := new(detailModel)
	// 詳細ファイルの取得
	filename := strings.ReplaceAll(l.Day, "-", "") + ".txt"
	filename = path.Join(diaryPath, filename)
	d.Day = l.Day
	d.Outline = l.Outline
	d.Tags = l.Tags
	if _, err := os.Stat(filename); err == nil {
		b, err := os.ReadFile(filename)
		if err != nil {
			return d
		}
		body := string(b)
		lines := strings.Split(body, "\n")
		if len(lines) > 3 {
			body = strings.Join(lines[3:], "\n")
		}
		d.Detail = body
	}
	return d
}

// 詳細情報を書き込む
func (d *detailModel) writeDetailFile() error {
	month := d.Day[0:4] + d.Day[5:7]

	// monthファイルの更新
	m, err := readListFile(month)
	if err != nil {
		return err
	}
	target := lineModel{d.Day, d.Outline, d.Tags, len(d.Detail) > 0, 0}
	if err := m.updateLine(month, target); err != nil {
		return err
	}

	// 詳細ファイルの更新
	filename := strings.ReplaceAll(d.Day, "-", "") + ".txt"
	filename = path.Join(diaryPath, filename)
	if len(d.Detail) == 0 {
		if _, err := os.Stat(filename); err == nil {
			return os.Remove(filename)
		}
		return nil
	}

	data := d.Day + "\n" + d.Outline + "\n" + strings.Join(d.Tags, "#") + "\n" + d.Detail
	return os.WriteFile(filename, []byte(data), fs.ModePerm)
}
