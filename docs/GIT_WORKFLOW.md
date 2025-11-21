# Git運用ルール

このドキュメントでは、Effisioプロジェクトにおける Git の運用ルールを定義します。

## 目次

1. [ブランチ戦略](#ブランチ戦略)
2. [コミットメッセージ規約](#コミットメッセージ規約)
3. [プルリクエスト運用](#プルリクエスト運用)
4. [コードレビュー](#コードレビュー)
5. [リリースフロー](#リリースフロー)
6. [緊急対応（Hotfix）](#緊急対応hotfix)

---

## ブランチ戦略

### Git Flow を採用

Effisioプロジェクトでは **Git Flow** をベースとしたブランチ戦略を採用します。

### ブランチの種類

#### 1. main ブランチ

- **用途**: 本番環境にデプロイされるコード
- **保護設定**: 直接コミット禁止、必ずPRを経由
- **マージ可能**: `release/*`, `hotfix/*` ブランチのみ
- **タグ**: すべてのリリースにバージョンタグを付与（例: `v1.0.0`）

```bash
# タグの命名規則
v<major>.<minor>.<patch>

# 例
v1.0.0  # メジャーリリース
v1.1.0  # マイナーアップデート
v1.1.1  # パッチ（バグ修正）
```

#### 2. develop ブランチ

- **用途**: 次期リリースの開発統合ブランチ
- **保護設定**: 直接コミット禁止、必ずPRを経由
- **マージ可能**: `feature/*`, `bugfix/*` ブランチ
- **デプロイ先**: ステージング環境（自動デプロイ）

#### 3. feature/* ブランチ

- **用途**: 新機能開発
- **命名規則**: `feature/<機能名>` または `feature/<issue番号>-<機能名>`
- **起点**: `develop` ブランチ
- **マージ先**: `develop` ブランチ
- **ライフサイクル**: マージ後に削除

```bash
# 命名例
feature/user-authentication
feature/123-user-authentication
feature/dashboard-ui
```

#### 4. bugfix/* ブランチ

- **用途**: 開発中のバグ修正
- **命名規則**: `bugfix/<バグ内容>` または `bugfix/<issue番号>-<バグ内容>`
- **起点**: `develop` ブランチ
- **マージ先**: `develop` ブランチ
- **ライフサイクル**: マージ後に削除

```bash
# 命名例
bugfix/login-validation
bugfix/456-login-validation
```

#### 5. release/* ブランチ

- **用途**: リリース準備（バージョン番号更新、最終テスト）
- **命名規則**: `release/v<version>`
- **起点**: `develop` ブランチ
- **マージ先**: `main` と `develop` の両方
- **ライフサイクル**: マージ後に削除

```bash
# 命名例
release/v1.0.0
release/v1.1.0
```

#### 6. hotfix/* ブランチ

- **用途**: 本番環境の緊急バグ修正
- **命名規則**: `hotfix/v<version>` または `hotfix/<バグ内容>`
- **起点**: `main` ブランチ
- **マージ先**: `main` と `develop` の両方
- **ライフサイクル**: マージ後に削除

```bash
# 命名例
hotfix/v1.0.1
hotfix/critical-security-fix
```

### ブランチ図

```
main      ─────●─────────────●─────────────●────> (本番)
               │             │             │
               │ (v1.0.0)    │ (v1.1.0)    │ (v1.0.1 hotfix)
               │             │             │
develop   ─┬───┴────┬────────┴────┬────────┴───┬──> (ステージング)
           │        │             │            │
feature/A  └─●──●───┘             │            │
feature/B          └───●──●───────┘            │
hotfix/1.0.1                                   └──●──●
```

---

## コミットメッセージ規約

### Conventional Commits を採用

すべてのコミットメッセージは **Conventional Commits** 形式に従います。

### 基本フォーマット

```
<type>(<scope>): <subject>

<body>

<footer>
```

#### 必須部分

```
<type>(<scope>): <subject>
```

- **type**: コミットの種類（必須）
- **scope**: 変更箇所（オプション）
- **subject**: 変更内容の要約（必須、50文字以内）

#### type の種類

| type | 説明 | 例 |
|------|------|-----|
| **feat** | 新機能追加 | `feat(auth): ユーザーログイン機能を追加` |
| **fix** | バグ修正 | `fix(users): メールバリデーションを修正` |
| **docs** | ドキュメント変更 | `docs(readme): セットアップ手順を更新` |
| **style** | コードフォーマット（動作に影響なし） | `style(backend): golangci-lint適用` |
| **refactor** | リファクタリング | `refactor(auth): JWT生成ロジックを整理` |
| **perf** | パフォーマンス改善 | `perf(db): クエリにインデックスを追加` |
| **test** | テスト追加・修正 | `test(auth): ログインテストを追加` |
| **build** | ビルドシステム・依存関係変更 | `build(deps): gin を v1.9.0 に更新` |
| **ci** | CI/CD設定変更 | `ci(actions): テストワークフローを追加` |
| **chore** | その他の変更 | `chore(gitignore): .envファイルを除外` |
| **revert** | コミット取り消し | `revert: "feat(auth): ログイン機能追加"` |

#### scope の例

- `auth`: 認証関連
- `users`: ユーザー管理
- `api`: API全般
- `db`: データベース
- `ui`: UI/フロントエンド
- `docs`: ドキュメント
- `ci`: CI/CD
- `deps`: 依存関係

### コミットメッセージの例

#### 良い例 ✅

```
feat(auth): JWT認証機能を実装

- JWTトークンの生成・検証機能を追加
- ミドルウェアでトークンをチェック
- リフレッシュトークン対応

Closes #123
```

```
fix(users): ユーザー削除時のNULLポインタエラーを修正

削除済みユーザーの参照時にNULLチェックを追加

Fixes #456
```

```
docs(setup): 開発環境セットアップ手順を追加

macOS、Windows、Linuxそれぞれの手順を記載
```

#### 悪い例 ❌

```
update code  # type がない、内容が不明確
```

```
feat: いろいろ修正した  # 内容が不明確、複数の変更が混在
```

```
fixed bug  # type が間違い（fix が正しい）、具体性がない
```

### コミットの粒度

#### ✅ 推奨

- **1コミット = 1つの論理的な変更**
- 独立してレビュー・テストできる単位
- リバート可能な単位

#### ❌ 避けるべき

- 複数の機能を1コミットにまとめる
- 無関係な変更を含める
- 大量のファイルを一度にコミット

---

## プルリクエスト運用

### PRの作成ルール

#### 1. PRのタイミング

- 機能が完成してからPRを作成（WIPの場合はドラフトPR）
- すべてのテストが通っていることを確認
- コンフリクトを解消してから作成

#### 2. PRのタイトル

コミットメッセージと同様に Conventional Commits 形式を使用：

```
feat(auth): ユーザーログイン機能を実装
fix(users): メールバリデーションを修正
docs(api): Swagger定義を更新
```

#### 3. PRの説明（テンプレート）

PRには以下の情報を含めてください：

```markdown
## 概要
この変更の目的を簡潔に説明

## 変更内容
- 変更点1
- 変更点2
- 変更点3

## 関連Issue
Closes #123
Relates to #456

## テスト方法
1. ステップ1
2. ステップ2
3. 期待される結果

## スクリーンショット（UIの場合）
変更前と変更後のスクリーンショット

## チェックリスト
- [ ] テストを追加・更新した
- [ ] ドキュメントを更新した
- [ ] すべてのテストが通る
- [ ] リントエラーがない
- [ ] コンフリクトがない
- [ ] レビュー依頼を送った
```

#### 4. レビュアーの指定

- **最低1名**のレビュアーを指定
- **技術的に関連性の高いメンバー**を選択
- バックエンド変更 → バックエンドエンジニア
- フロントエンド変更 → フロントエンドエンジニア

#### 5. ラベルの付与

以下のラベルを適切に付与：

| ラベル | 説明 |
|--------|------|
| `feature` | 新機能 |
| `bugfix` | バグ修正 |
| `enhancement` | 改善 |
| `documentation` | ドキュメント |
| `breaking change` | 破壊的変更 |
| `needs review` | レビュー待ち |
| `work in progress` | 作業中 |
| `ready to merge` | マージ可能 |

---

## コードレビュー

### レビュアーの責任

#### チェック項目

- [ ] コードが要件を満たしているか
- [ ] 設計が適切か（SOLID原則、DRY原則など）
- [ ] エラーハンドリングが適切か
- [ ] セキュリティ上の問題がないか
- [ ] パフォーマンスの問題がないか
- [ ] テストが十分か
- [ ] コーディング規約に従っているか
- [ ] ドキュメントが更新されているか
- [ ] コメントが適切か

#### レビューコメントの書き方

##### ✅ 良いコメント

```
提案: エラーハンドリングを追加した方が良いと思います。
この部分で err != nil のチェックが必要です。

理由:
ユーザー入力に不正な値が来た場合にpanicが発生する可能性があります。
```

```
質問: この実装を選択した理由を教えてください。
別のアプローチ（例: キャッシング）も検討されましたか?
```

```
nit: 変数名を `userData` から `user` に変更すると、より簡潔になります。
```

##### ❌ 悪いコメント

```
ダメです。  # 理由が不明確
```

```
これは間違っています。全部書き直してください。  # 具体性がない、攻撃的
```

#### コメントのプレフィックス

| プレフィックス | 意味 | 対応 |
|--------------|------|------|
| **MUST** | 必須の修正 | 必ず対応が必要 |
| **SHOULD** | 推奨の修正 | 対応が望ましい |
| **nit** | 些細な指摘 | 対応は任意 |
| **question** | 質問 | 回答のみで良い |
| **suggestion** | 提案 | 検討してほしい |

### レビューのタイムライン

- **24時間以内**に初回レビューを開始
- **2営業日以内**にレビューを完了
- 緊急の場合は即座にレビュー

### 承認基準

#### マージ可能条件

- [ ] **最低1名の承認**（Approve）を得ている
- [ ] すべての会話が解決済み（Resolved）
- [ ] CI/CDパイプラインが成功
- [ ] コンフリクトが解消されている
- [ ] レビューアの修正要求に対応済み

---

## リリースフロー

### 1. リリースブランチの作成

```bash
# developから最新を取得
git checkout develop
git pull origin develop

# releaseブランチを作成
git checkout -b release/v1.0.0
```

### 2. バージョン番号の更新

以下のファイルのバージョン番号を更新：

```bash
# バックエンド
backend/cmd/server/main.go  # Version変数

# フロントエンド
frontend/package.json       # version フィールド

# ドキュメント
CHANGELOG.md               # 変更履歴
```

### 3. CHANGELOGの更新

```markdown
# Changelog

## [1.0.0] - 2024-01-16

### Added
- ユーザー認証機能（JWT）
- ユーザー管理CRUD
- ダッシュボード

### Changed
- データベーススキーマを最適化

### Fixed
- ログインバリデーションのバグ修正

### Security
- OWASP Top 10対策を実装
```

### 4. 最終テスト

```bash
# バックエンドテスト
cd backend
go test ./... -v

# フロントエンドテスト
cd frontend
npm run test
npm run build
```

### 5. PRの作成

```
タイトル: release: v1.0.0

説明:
## リリース内容
- 新機能A
- バグ修正B
- 改善C

## テスト結果
- [ ] ユニットテスト: ✅
- [ ] 統合テスト: ✅
- [ ] E2Eテスト: ✅
- [ ] ステージング環境で動作確認: ✅
```

### 6. mainへマージ

```bash
# mainにマージ（GitHub PR経由）
# マージ後、mainからタグを作成

git checkout main
git pull origin main
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

### 7. developへもマージ

```bash
# releaseブランチの変更をdevelopにも反映
git checkout develop
git merge release/v1.0.0
git push origin develop

# releaseブランチを削除
git branch -d release/v1.0.0
git push origin --delete release/v1.0.0
```

---

## 緊急対応（Hotfix）

### 1. Hotfixブランチの作成

```bash
# mainから最新を取得
git checkout main
git pull origin main

# hotfixブランチを作成
git checkout -b hotfix/v1.0.1
```

### 2. バグ修正

緊急のバグを修正し、テストを追加。

### 3. バージョン番号の更新

パッチバージョンを上げる（例: 1.0.0 → 1.0.1）

### 4. mainとdevelopの両方にマージ

```bash
# mainにマージ
git checkout main
git merge hotfix/v1.0.1
git tag -a v1.0.1 -m "Hotfix: critical bug fix"
git push origin main
git push origin v1.0.1

# developにもマージ
git checkout develop
git merge hotfix/v1.0.1
git push origin develop

# hotfixブランチを削除
git branch -d hotfix/v1.0.1
git push origin --delete hotfix/v1.0.1
```

---

## 便利なGitコマンド

### ブランチ操作

```bash
# ブランチ一覧
git branch -a

# リモートの最新を取得
git fetch origin

# ブランチを削除
git branch -d <branch-name>

# リモートブランチを削除
git push origin --delete <branch-name>

# ブランチ名を変更
git branch -m <old-name> <new-name>
```

### コミット操作

```bash
# 最後のコミットを修正
git commit --amend

# コミットメッセージのみ修正
git commit --amend -m "新しいメッセージ"

# 過去のコミットをまとめる
git rebase -i HEAD~3

# コミットを取り消す（履歴を残す）
git revert <commit-hash>

# コミットを取り消す（履歴を消す）
git reset --hard <commit-hash>
```

### リモート操作

```bash
# リモートの状態を確認
git remote -v

# リモートから最新を取得してマージ
git pull origin develop

# リモートから最新を取得（マージなし）
git fetch origin

# ローカルの変更を強制的にプッシュ（注意！）
git push -f origin <branch-name>
```

### 便利なエイリアス

`.gitconfig` に追加：

```ini
[alias]
  st = status
  co = checkout
  br = branch
  ci = commit
  unstage = reset HEAD --
  last = log -1 HEAD
  lg = log --oneline --graph --all --decorate
  amend = commit --amend --no-edit
```

---

## まとめ

このGit運用ルールに従うことで：

✅ **一貫性**: 誰が見ても理解しやすい履歴
✅ **品質**: レビュープロセスによる品質向上
✅ **安全性**: main ブランチの保護
✅ **効率性**: 明確なワークフローによる開発効率向上

不明点があれば、チームのSlackチャンネルで質問してください。
