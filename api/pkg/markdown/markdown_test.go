package markdown

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestToHTML_ConvertsSampleMarkdown(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	body, err := os.ReadFile(filepath.Join(root, "markdown.md"))
	if err != nil {
		t.Fatalf("read markdown.md: %v", err)
	}

	got := ToHTML(string(body))

	expected := strings.Join([]string{
		"<h1>見出し１</h1>",
		"<h2>見出し２</h2>",
		"<h3>見出し３</h3>",
		"<p>ノーマル</p>",
		"<ol>",
		"<li>数字付きリスト</li>",
		"<li>数字付きリスト２</li>",
		"</ol>",
		"<ul>",
		"<li>リスト１</li>",
		"<li>リスト２</li>",
		"</ul>",
		"<pre><code>コード\n</code></pre>",
		"<blockquote>",
		"<p>引用単体行</p>",
		"</blockquote>",
		"<blockquote>",
		"<p>引用複数行<br>引用２行目</p>",
		"</blockquote>",
		"<p>文字装飾「<strong>太字</strong>」「<em>斜線</em>」「<u>下線</u>」「<del>取り消し線</del>」「<a href=\"https://google.co.jp\">リンク</a>」</p>",
		"<p><img src=\"./diary/2026-04-03/001.png\" alt=\"イメージ\"></p>",
		"<hr>",
		"<ul>",
		"<li>親リスト１",
		"<ul>",
		"<li>子リスト１</li>",
		"<li>子リスト２</li>",
		"</ul></li>",
		"<li>親リスト２</li>",
		"</ul>",
		"<table>",
		"<thead>",
		"<tr><th>名前</th><th>説明</th></tr>",
		"</thead>",
		"<tbody>",
		"<tr><td>サンプル1</td><td>テーブル１行目</td></tr>",
		"<tr><td>サンプル2</td><td>テーブル２行目</td></tr>",
		"</tbody>",
		"</table>",
	}, "\n")

	if got != expected {
		t.Fatalf("unexpected html\nexpected:\n%s\n\ngot:\n%s", expected, got)
	}
}

func TestToHTML_SupportsTablesNestedListsAndDecorations(t *testing.T) {
	input := strings.Join([]string{
		"- 親1",
		"  - 子1",
		"  - 子2",
		"- 親2",
		"",
		"| 名前 | 説明 |",
		"| --- | --- |",
		"| a | `code` |",
		"| b | [link](https://example.com) |",
		"",
		"---",
		"",
		"__下線__ と ~~取り消し~~",
	}, "\n")

	got := ToHTML(input)
	expected := strings.Join([]string{
		"<ul>",
		"<li>親1",
		"<ul>",
		"<li>子1</li>",
		"<li>子2</li>",
		"</ul></li>",
		"<li>親2</li>",
		"</ul>",
		"<table>",
		"<thead>",
		"<tr><th>名前</th><th>説明</th></tr>",
		"</thead>",
		"<tbody>",
		"<tr><td>a</td><td><code>code</code></td></tr>",
		"<tr><td>b</td><td><a href=\"https://example.com\">link</a></td></tr>",
		"</tbody>",
		"</table>",
		"<hr>",
		"<p><u>下線</u> と <del>取り消し</del></p>",
	}, "\n")

	if got != expected {
		t.Fatalf("unexpected html\nexpected:\n%s\n\ngot:\n%s", expected, got)
	}
}

func TestToHTML_HidesCarryOverMarkerLine(t *testing.T) {
	input := strings.Join([]string{
		"## お店",
		"テストデータ",
		"",
		"----ここまで前回内容で置換",
		"## お話",
		"----",
	}, "\n")

	got := ToHTML(input)
	expected := strings.Join([]string{
		"<h2>お店</h2>",
		"<p>テストデータ</p>",
		"<hr>",
		"<h2>お話</h2>",
		"<hr>",
	}, "\n")

	if got != expected {
		t.Fatalf("unexpected html\nexpected:\n%s\n\ngot:\n%s", expected, got)
	}
}
