# notion-json-to-md

Notion APIから取得したブロックデータ（JSON形式）をMarkdownに変換するCLIツール

## 概要

Notion APIの `/blocks/{block_id}/children` エンドポイントから取得したJSONデータを、読みやすいMarkdown形式に変換します。

## インストール

```bash
go build -o notion-json-to-md
```

## 使い方

標準入力からJSONを受け取り、標準出力にMarkdownを出力します。

```bash
# 基本的な使い方
cat blocks.json | ./notion-json-to-md > output.md

# ファイルリダイレクト
./notion-json-to-md < input.json > output.md
```

## サポートしているブロックタイプ

### 見出し
- `heading_1` → `# タイトル`
- `heading_2` → `## タイトル`
- `heading_3` → `### タイトル`

### テキストブロック
- `paragraph` → 通常のテキスト
- `bulleted_list_item` → `- リストアイテム`
- `numbered_list_item` → `1. リストアイテム`

### コードブロック
- `code` → ````language\nコード\n````

## サポートしているアノテーション

- **太字** → `**text**`
- *イタリック* → `*text*`
- `コード` → `` `text` ``
- ~~取り消し線~~ → `~~text~~`
- リンク → `[text](url)`
- 複数アノテーションの組み合わせ

## 入力形式

Notion APIの `/blocks/{block_id}/children` エンドポイントのレスポンス形式：

```json
{
  "object": "list",
  "results": [
    {
      "type": "heading_1",
      "heading_1": {
        "rich_text": [
          {
            "plain_text": "タイトル",
            "annotations": {
              "bold": false,
              "italic": false,
              "code": false,
              "strikethrough": false
            }
          }
        ]
      }
    }
  ]
}
```

## 開発

### テスト実行

```bash
go test -v
```

## ライセンス

MIT License
