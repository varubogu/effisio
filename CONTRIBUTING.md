# コントリビューションガイド

Effisioプロジェクトへの貢献に興味を持っていただき、ありがとうございます！

## 目次

- [行動規範](#行動規範)
- [開始方法](#開始方法)
- [開発フロー](#開発フロー)
- [コーディング規約](#コーディング規約)
- [コミットメッセージ](#コミットメッセージ)
- [プルリクエスト](#プルリクエスト)
- [レビュープロセス](#レビュープロセス)
- [質問やサポート](#質問やサポート)

## 行動規範

### 基本原則

- **尊重**: 全ての参加者を尊重し、建設的なフィードバックを提供する
- **協力**: チームワークを大切にし、知識を共有する
- **品質**: 高品質なコードとドキュメントを維持する
- **透明性**: 問題や懸念事項をオープンに共有する

### 禁止事項

- ハラスメントや差別的な言動
- 他者への攻撃的なコメント
- プライベート情報の無断公開
- プロフェッショナルでない振る舞い

## 開始方法

### 1. 環境セットアップ

```bash
# リポジトリをクローン
git clone https://github.com/your-org/effisio.git
cd effisio

# 初回セットアップ実行
make setup

# 開発環境起動
make dev
```

詳細は [docs/DEVELOPMENT_SETUP.md](docs/DEVELOPMENT_SETUP.md) を参照してください。

### 2. ドキュメントを読む

貢献を始める前に、以下のドキュメントを確認してください：

- [開発環境セットアップ](docs/DEVELOPMENT_SETUP.md)
- [Gitワークフロー](docs/GIT_WORKFLOW.md)
- [Goコーディングガイドライン](docs/CODING_GUIDELINES_GO.md)
- [TypeScriptコーディングガイドライン](docs/CODING_GUIDELINES_TYPESCRIPT.md)

## 開発フロー

### 1. Issueを作成または選択

- 新しい機能や修正を実装する前に、Issueを作成または既存のIssueを選択
- Issueには明確なタイトルと説明を記載
- 必要に応じてラベルを追加

### 2. ブランチを作成

```bash
# 最新のdevelopブランチを取得
git checkout develop
git pull origin develop

# 新しいブランチを作成
git checkout -b feature/issue-123-add-user-auth
```

ブランチ命名規則：
- `feature/issue-{番号}-{簡潔な説明}` - 新機能
- `bugfix/issue-{番号}-{簡潔な説明}` - バグ修正
- `hotfix/issue-{番号}-{簡潔な説明}` - 緊急修正
- `docs/issue-{番号}-{簡潔な説明}` - ドキュメント更新

### 3. 開発作業

#### コードを書く

- コーディング規約に従う
- 適切なテストを追加
- エッジケースを考慮

#### テストを実行

```bash
# バックエンド
cd backend
make test
make lint

# フロントエンド
cd frontend
npm test
npm run lint
```

#### ローカルで動作確認

```bash
# 開発環境で確認
make dev

# ブラウザで確認
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

### 4. コミット

```bash
# 変更をステージング
git add .

# コミット（Conventional Commitsに従う）
git commit -m "feat(auth): JWT認証機能を実装

- JWTトークンの生成・検証機能を追加
- ミドルウェアでトークンをチェック
- ログイン・ログアウトエンドポイントを実装

Closes #123"
```

### 5. プルリクエストを作成

```bash
# リモートにプッシュ
git push origin feature/issue-123-add-user-auth
```

GitHubでプルリクエストを作成し、テンプレートに従って記入してください。

## コーディング規約

### Go

詳細は [docs/CODING_GUIDELINES_GO.md](docs/CODING_GUIDELINES_GO.md) を参照。

**主要な規約：**
- `gofmt` でフォーマット
- エラーハンドリングを適切に行う
- テストカバレッジ70%以上を目指す
- godocコメントを記述

**例：**
```go
// CreateUser は新しいユーザーを作成します
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }

    user, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return user, nil
}
```

### TypeScript/React

詳細は [docs/CODING_GUIDELINES_TYPESCRIPT.md](docs/CODING_GUIDELINES_TYPESCRIPT.md) を参照。

**主要な規約：**
- `any` を使わない（strict mode）
- Functional Componentを使用
- Hooksを適切に使用
- Props型を明示的に定義

**例：**
```typescript
interface UserProfileProps {
  userId: number;
  onUpdate?: (user: User) => void;
}

export function UserProfile({ userId, onUpdate }: UserProfileProps) {
  const { data: user, isLoading } = useUser(userId);

  if (isLoading) return <LoadingSpinner />;
  if (!user) return <NotFound />;

  return (
    <div className="user-profile">
      <h1>{user.username}</h1>
      {/* ... */}
    </div>
  );
}
```

## コミットメッセージ

Conventional Commits形式を使用します。

### フォーマット

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type

- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメントのみの変更
- `style`: コードの意味に影響しない変更（空白、フォーマット等）
- `refactor`: バグ修正や機能追加ではないコード変更
- `perf`: パフォーマンス改善
- `test`: テストの追加・修正
- `chore`: ビルドプロセスやツールの変更

### 例

```
feat(auth): JWT認証機能を実装

- JWTトークンの生成・検証機能を追加
- ミドルウェアでトークンをチェック
- ログイン・ログアウトエンドポイントを実装

Closes #123
```

詳細は [docs/GIT_WORKFLOW.md](docs/GIT_WORKFLOW.md) を参照。

## プルリクエスト

### PRを作成する前のチェックリスト

- [ ] 全てのテストがパスする
- [ ] リンターエラーがない
- [ ] 適切なテストを追加した
- [ ] ドキュメントを更新した（必要な場合）
- [ ] CHANGELOG.mdを更新した
- [ ] コミットメッセージが規約に従っている

### PR説明の書き方

テンプレート（.github/PULL_REQUEST_TEMPLATE.md）に従って記入：

1. **概要**: 何を変更したか
2. **変更内容**: 具体的な変更点
3. **テスト**: どのようにテストしたか
4. **スクリーンショット**: UI変更の場合
5. **関連Issue**: `Closes #123`

### レビュー依頼

- 適切なレビュアーを指定
- 必要に応じてラベルを追加
- CI/CDが全てパスしていることを確認

## レビュープロセス

### レビュアーの責任

- コードの品質を確認
- セキュリティ上の問題がないか確認
- テストが適切か確認
- 建設的なフィードバックを提供

### レビュー時のチェックポイント

#### コード品質
- [ ] コーディング規約に従っているか
- [ ] 適切なエラーハンドリングがされているか
- [ ] 命名が分かりやすいか
- [ ] 重複コードがないか

#### セキュリティ
- [ ] SQLインジェクション対策がされているか
- [ ] XSS対策がされているか
- [ ] 認証・認可が適切か
- [ ] 機密情報がハードコードされていないか

#### テスト
- [ ] 十分なテストカバレッジがあるか
- [ ] エッジケースをテストしているか
- [ ] テストが失敗しないか

#### ドキュメント
- [ ] コードコメントが適切か
- [ ] README等の更新が必要か
- [ ] API仕様が更新されているか

### フィードバックの反映

- レビューコメントに対して返信または修正
- 修正をコミット＆プッシュ
- 再度レビューを依頼

### マージ

- 最低1名の承認が必要
- 全てのCIチェックがパスしている必要がある
- Squash and Mergeを使用（履歴を綺麗に保つ）

## 質問やサポート

### 困ったときは

1. **ドキュメントを確認**: docsフォルダ内のドキュメントを確認
2. **既存のIssueを検索**: 同じ問題がないか確認
3. **Issueを作成**: 解決しない場合は新しいIssueを作成
4. **チームに相談**: チャットやミーティングで相談

### Issue作成時の注意

- 明確なタイトルをつける
- 再現手順を記載する
- 期待される動作と実際の動作を説明
- 環境情報を含める（OS、バージョン等）

## ライセンス

このプロジェクトに貢献することで、あなたの貢献がプロジェクトと同じライセンス（MITライセンス）の下で配布されることに同意したものとみなされます。

---

貢献していただき、ありがとうございます！ 🎉
