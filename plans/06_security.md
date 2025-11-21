# セキュリティ計画

## 認証・認可

### パスワード管理
- **ハッシュ関数**: bcrypt (cost: 12)
- **パスワード要件**:
  - 最小8文字
  - 大文字・小文字・数字・記号を含む
- **パスワード有効期限**: 90日（オプション）

### JWT トークン
- **署名アルゴリズム**: HS256 または RS256
- **有効期限**:
  - Access Token: 1時間
  - Refresh Token: 30日
- **格納場所**: HttpOnly Cookie または Authorization ヘッダ

### セッション管理
- Stateless 認証を基本とする
- Redis でセッション情報をキャッシュ（オプション）

---

## HTTPS/TLS

- **最小バージョン**: TLS 1.2
- **推奨**: TLS 1.3
- **証明書**: Let's Encrypt (自動更新)
- **HSTS**: max-age=31536000; includeSubDomains

---

## CORS (Cross-Origin Resource Sharing)

```go
// 許可するオリジン
AllowedOrigins: []string{
  "https://internal-system.com",
  "https://app.internal-system.com",
}

// 許可するメソッド
AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}

// 許可するヘッダ
AllowedHeaders: []string{"Content-Type", "Authorization"}

// キャッシュ時間
MaxAge: 86400
```

---

## OWASP Top 10 対策

### 1. Injection (インジェクション)
- **SQL Injection**: パラメータ化クエリ、ORM 使用
- **Command Injection**: シェルコマンド実行を避ける
- **対策**: 入力検証、プリペアドステートメント

### 2. Broken Authentication
- **対策**:
  - 強力なパスワード要件
  - MFA (Multi-Factor Authentication) の実装
  - ブルートフォース攻撃対策（レート制限）

### 3. Sensitive Data Exposure
- **対策**:
  - HTTPS の強制
  - 機密情報を環境変数で管理
  - ログにはセンシティブ情報を記録しない

### 4. XML External Entities (XXE)
- **対策**: XML パーサーで外部エンティティを無効化

### 5. Access Control (不適切なアクセス制御)
- **対策**:
  - ロールベースアクセス制御 (RBAC)
  - 各エンドポイントで権限チェック
  - リソースレベルの権限確認

### 6. Security Misconfiguration (不適切なセキュリティ設定)
- **対策**:
  - デフォルト認証情報の削除
  - 不要なサービスの無効化
  - セキュリティヘッダーの設定

### 7. Cross-Site Scripting (XSS)
- **対策**:
  - HTML エスケープ
  - Content Security Policy (CSP) の設定
  - HttpOnly Cookie の使用

### 8. Insecure Deserialization
- **対策**:
  - JSON パーサーの信頼性確保
  - 外部ソースからのデータは検証

### 9. Using Components with Known Vulnerabilities
- **対策**:
  - 定期的な依存関係更新
  - `go list -m all` で脆弱性チェック
  - CI/CD で自動チェック

### 10. Insufficient Logging & Monitoring
- **対策**:
  - ログに認証試行を記録
  - 異常なアクティビティの検出
  - アラート設定

---

## 入力検証

### バリデーションルール

```go
// メール検証
Email: "required,email,max=255"

// ユーザー名
Username: "required,min=3,max=50,alphanumeric"

// パスワード
Password: "required,min=8,max=128"

// URL
URL: "url"

// 数値
Age: "numeric,min=0,max=150"
```

---

## レート制限

### 実装戦略

```
// IP ベース
- 100 リクエスト / 分

// ユーザーベース
- 1000 リクエスト / 分

// 特定エンドポイント
- /auth/login: 10 試行 / 15分
- /users: 100 リクエスト / 分
```

---

## CSRF (Cross-Site Request Forgery) 対策

- **トークンベース**: CSRF トークンの使用
- **同一サイト Cookie**: SameSite=Strict/Lax
- **二重送信 Cookie**: GET と POST で値が一致するか確認

---

## ログとモニタリング

### ログ記録項目

```
- ユーザー認証 (成功/失敗)
- データベース操作 (INSERT/UPDATE/DELETE)
- エラー内容
- リクエスト ID
- ユーザー ID
- IP アドレス
- タイムスタンプ
```

### 監視対象

- 複数の認証失敗
- 異常なアクセスパターン
- エラーレート上昇
- 応答時間の増加
- リソース使用率

---

## 依存関係管理

### 定期的なチェック

```bash
# 脆弱性スキャン
go list -m all | go-audit

# または nancy を使用
nancy sleuth

# 最新バージョン確認
go list -u -m all
```

---

## 本番環境デプロイ前チェックリスト

- [ ] HTTPS 有効化
- [ ] CORS 設定確認
- [ ] パスワード要件設定
- [ ] JWT キー設定 (環境変数)
- [ ] ロギング設定
- [ ] レート制限有効化
- [ ] セキュリティヘッダー設定
- [ ] 依存関係の脆弱性チェック
- [ ] テストケース作成・実行
- [ ] セキュリティ監査実施
