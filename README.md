# notion-to-md

NotionページをMarkdownに変換するCLIツール

## 概要

Notion APIを使ってページのブロックを再帰的に取得し、ネスト構造を保持したMarkdown形式に変換します。

## 特徴

- ✅ Notion URLまたはBlock IDから直接変換
- ✅ ネストされたリストに対応（最大10階層）
- ✅ ページネーション対応（100件以上のブロック）
- ✅ 豊富なブロックタイプサポート
- ✅ テキストアノテーション（太字、イタリック、コード、取り消し線、リンク）

## インストール

```bash
go install github.com/syou6162/notion-to-md@latest
```

または:

```bash
git clone https://github.com/syou6162/notion-to-md
cd notion-to-md
go build
```

## 使い方

### 環境変数の設定

Notion Integration Tokenを設定します:

```bash
export NOTION_TOKEN="secret_xxxxxxxxxxxxx"
```

### 実行

```bash
# Notion URLを指定
notion-to-md https://www.notion.so/workspace/Page-title-<page-id>

# Block IDを直接指定
notion-to-md <block-id>

# ファイルに保存
notion-to-md <block-id> > output.md
```

## サポートしているブロックタイプ

### 見出し
- `heading_1` → `# タイトル`
- `heading_2` → `## タイトル`
- `heading_3` → `### タイトル`

### テキストブロック
- `paragraph` → 通常のテキスト
- `bulleted_list_item` → `- リストアイテム` (ネスト対応)
- `numbered_list_item` → `1. リストアイテム` (ネスト対応)
- `toggle` → `- トグル`
- `quote` → `> 引用`
- `callout` → `> コールアウト`

### その他
- `code` → ````language\nコード\n````
- `divider` → `---`

## サポートしているアノテーション

- **太字** → `**text**`
- *イタリック* → `*text*`
- `コード` → `` `text` ``
- ~~取り消し線~~ → `~~text~~`
- リンク → `[text](url)`
- 複数アノテーションの組み合わせ

## ネストされたリストの例

Notion:
```
- 親アイテム
  - 子アイテム1
    - 孫アイテム
  - 子アイテム2
```

変換後:
```markdown
- 親アイテム
  - 子アイテム1
    - 孫アイテム
  - 子アイテム2
```

## 制限事項

- 再帰深さ: 最大10階層
- ページネーション: 100ブロック/ページ（自動対応）
- 一部のブロックタイプ（画像、テーブル等）は未対応

## Notion Integration Tokenの取得方法

1. [Notion Integrations](https://www.notion.so/my-integrations) にアクセス
2. "New integration" をクリック
3. Integration名を入力して作成
4. "Internal Integration Token" をコピー
5. 対象のNotionページで、"Share" → Integrationを追加

## ライセンス

MIT License
