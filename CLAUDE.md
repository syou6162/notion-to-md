# CLAUDE.md

このファイルは、Claude Codeなどの開発支援AIが参照するためのプロジェクト情報です。

## プロジェクト概要

Notion APIを使ってページのブロックを再帰的に取得し、Markdownに変換するCLIツール。
Notion URLまたはBlock IDを引数として受け取り、標準出力にMarkdownを出力する。

## アーキテクチャ

### ファイル構成

```
.
├── main.go      # エントリーポイント（CLI引数処理、URL/ID解析）
├── fetcher.go   # Notion API呼び出し（再帰的取得、ページネーション処理）
├── converter.go # Markdown変換ロジック（ブロック→Markdown、アノテーション処理）
├── go.mod       # Goモジュール定義
└── README.md    # ユーザー向けドキュメント
```

### 依存ライブラリ

- `github.com/jomei/notionapi` v1.13.3 - Notion API Go SDK

### 主要な型定義

#### BlockWithIndent
ブロックとインデントレベルを保持する構造体

```go
type BlockWithIndent struct {
    Block  notionapi.Block
    Indent int
}
```

### 主要な関数

#### `main()`
CLIエントリーポイント。
1. コマンドライン引数からNotion URLまたはBlock IDを取得
2. `extractBlockID`でBlock IDを抽出
3. 環境変数`NOTION_TOKEN`からトークンを取得
4. Notion APIクライアント初期化
5. `fetchAllBlocks`でブロック取得
6. `convert`でMarkdown変換
7. 標準出力に出力

#### `extractBlockID(input string) (notionapi.BlockID, error)`
Notion URLまたはBlock ID文字列からBlock IDを抽出する。
- URL形式: `https://www.notion.so/workspace/Title-{32-hex}`から32桁の16進数を抽出し、UUID形式（ハイフン付き）に変換
- ID形式: そのまま返す

#### `fetchAllBlocks(ctx, client, blockID) ([]BlockWithIndent, error)`
ブロックIDを受け取り、そのブロックの子要素を再帰的に取得する。
内部的に`fetchAllBlocksRecursive`を呼び出す（初期深度: 0）。

#### `fetchAllBlocksRecursive(ctx, client, blockID, depth) ([]BlockWithIndent, error)`
再帰的にブロックを取得する内部関数。
- 最大深度: 10階層（超過時はエラー）
- `fetchBlockChildren`で子ブロックを取得
- 各ブロックの`has_children`フラグをチェックして子要素を再帰的に取得
- 各階層のインデントレベル（depth）を保持
- デバッグログを出力（os.Stderr）

#### `fetchBlockChildren(ctx, client, blockID) ([]notionapi.Block, error)`
ページネーション対応のブロック子要素取得関数。
- Notion APIは100ブロック/ページで制限
- `has_more`フラグをチェックして全ページを取得
- `next_cursor`を使って次ページを取得
- デバッグログを出力（os.Stderr）

#### `convert(blocks []BlockWithIndent) string`
BlockWithIndent配列を受け取り、Markdown文字列を返すメイン変換関数。
各ブロックタイプ（notionapi.Block）に応じて適切なMarkdownフォーマットに変換する。
- インデント: 2スペース/レベル（`strings.Repeat("  ", bwi.Indent)`）
- 型アサーションで具体的なブロック型に変換してから処理

#### `formatRichText(richTexts []notionapi.RichText) string`
notionapi.RichText配列を受け取り、アノテーション（太字、イタリック、コード、取り消し線）とリンクを処理してMarkdown文字列を返す。

## サポートしているブロックタイプ

現在実装されているブロックタイプ:

- `heading_1` → `# タイトル`
- `heading_2` → `## タイトル`
- `heading_3` → `### タイトル`
- `paragraph` → 通常のテキスト
- `bulleted_list_item` → `- リスト`（ネスト対応）
- `numbered_list_item` → `1. リスト`（ネスト対応）
- `code` → ````language\ncode\n````
- `toggle` → `- トグル`（ネスト対応）
- `quote` → `> 引用`
- `callout` → `> コールアウト`
- `divider` → `---`

## アノテーション処理の順序

`formatRichText`関数では、以下の順序でアノテーションを適用：

1. Code (`` `text` ``)
2. Bold (`**text**`)
3. Italic (`*text*`)
4. Strikethrough (`~~text~~`)
5. Link (`[text](url)`)

この順序により、`**bold**`と`*italic*`が同時に適用された場合、`***bold italic***`として正しくレンダリングされる。

## 使用方法

### 環境変数

- `NOTION_TOKEN`: Notion Integration Token（必須）
  - 形式: `ntn_xxxxxxxxxxxxx` または `secret_xxxxxxxxxxxxx`

### コマンドライン

```bash
# Notion URLを指定
notion-to-md https://www.notion.so/workspace/Page-title-<page-id>

# Block IDを直接指定
notion-to-md <block-id>

# ファイルに保存
notion-to-md <block-id> > output.md
```

## 制限事項

- 再帰深さ: 最大10階層
- ページネーション: 100ブロック/ページ（自動対応）
- 大きなページ（400ブロック以上）は処理に5分以上かかる場合がある
- 一部のブロックタイプ（画像、テーブル等）は未対応

## 今後の拡張可能性

現在サポートしていないが、将来追加可能な機能：

### ブロックタイプ
- `table` - テーブル
- `image` - 画像
- `table_of_contents` - 目次
- `column_list` / `column` - カラムレイアウト

### 機能
- ストリーミング出力（現在は全取得後に一括出力）
- 複数ページの一括変換
- エラーハンドリングの強化
- 変換オプション（Markdownスタイルのカスタマイズなど）
- キャッシュ機構（API呼び出し削減）

## テスト

### テストファイル構成

- **converter_test.go**: 変換ロジックのテスト（21テストケース）
  - ブロックタイプ別変換（heading, paragraph, list, code, toggle, quote, callout, divider）
  - アノテーション処理（bold, italic, code, strikethrough, 複数アノテーション）
  - リンク、ネストリスト、複数ブロック、空段落、複数RichText要素
- **parser_test.go**: URL/Block ID解析のテスト（6テストケース）
  - Notion URL解析（https/http）
  - Block ID直接指定（ハイフン有無）
  - エラーケース
- **fetcher_test.go**: API取得ロジックのテスト（13テストケース）
  - ページネーション処理（単一ページ、複数ページ、空）
  - 再帰的取得（子なし、子あり、深いネスト、最大深度超過）
  - ネストされた箇条書き（bulleted_list_item、numbered_list_item）
  - 混在ブロックタイプ
  - エラーハンドリング

### テスト駆動開発

- BlockFetcherインターフェースを使った依存性注入により、外部APIに依存しないユニットテストを実現
- モック実装（mockBlockFetcher）を使用してNotion APIをシミュレート

### テスト実行

```bash
# 全テスト実行
go test -v

# 特定のテストのみ実行
go test -v -run TestConvert
go test -v -run TestFetch
go test -v -run TestExtract
```

## パフォーマンス特性

実測値（実際のNotionページでの測定）:
- 小規模ページ（30ブロック程度）: 数秒
- 中規模ページ（100ブロック程度）: 1-2分
- 大規模ページ（400ブロック程度）: 5分前後

処理時間の大部分はNotion API呼び出しによるもの。
ネストされたブロック（has_children=true）は追加のAPI呼び出しが必要なため、ネストが深いページほど時間がかかる。
