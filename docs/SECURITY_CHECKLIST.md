# セキュリティチェックリスト

## デプロイ前チェックリスト

### 認証・認可
- [ ] JWT シークレットキーが本番環境用に変更されている
- [ ] パスワードハッシュ化（bcrypt cost 12以上）
- [ ] セッションタイムアウトが適切に設定されている
- [ ] RBAC が正しく実装されている
- [ ] 認証トークンが HttpOnly Cookie または Authorization ヘッダで送信される

### HTTPS/TLS
- [ ] HTTPS が有効化されている
- [ ] TLS 1.2 以上を使用
- [ ] HSTS ヘッダーが設定されている
- [ ] SSL証明書が有効（Let's Encrypt など）

### CORS
- [ ] 許可するオリジンが本番環境のみに制限されている
- [ ] 許可するメソッド・ヘッダーが必要最小限
- [ ] Access-Control-Allow-Credentials が適切に設定

### 入力検証
- [ ] すべてのユーザー入力がバリデーションされている
- [ ] SQLインジェクション対策（パラメータ化クエリ）
- [ ] XSS対策（HTMLエスケープ）
- [ ] CSRF対策（トークン検証）
- [ ] ファイルアップロードの検証（サイズ・タイプ・拡張子）

### セキュリティヘッダー
- [ ] Content-Security-Policy
- [ ] X-Content-Type-Options: nosniff
- [ ] X-Frame-Options: DENY
- [ ] X-XSS-Protection: 1; mode=block
- [ ] Referrer-Policy: no-referrer

### データ保護
- [ ] 機密情報がログに出力されていない
- [ ] パスワードが平文で保存されていない
- [ ] APIレスポンスに不要な情報が含まれていない
- [ ] エラーメッセージが詳細すぎない

### 依存関係
- [ ] すべての依存関係が最新バージョン
- [ ] 既知の脆弱性がない（nancy, npm audit）
- [ ] 不要なパッケージが削除されている

### レート制限
- [ ] API エンドポイントにレート制限が設定されている
- [ ] ブルートフォース攻撃対策（ログイン試行制限）

### 監査ログ
- [ ] 重要な操作がログに記録されている
- [ ] ログが改ざん防止されている
- [ ] ログ保持期間が適切に設定されている

## OWASP Top 10 チェックリスト

1. [ ] **Injection** - パラメータ化クエリ使用
2. [ ] **Broken Authentication** - 強力な認証機構
3. [ ] **Sensitive Data Exposure** - データ暗号化
4. [ ] **XML External Entities (XXE)** - XMLパーサー設定
5. [ ] **Broken Access Control** - RBAC実装
6. [ ] **Security Misconfiguration** - デフォルト設定変更
7. [ ] **Cross-Site Scripting (XSS)** - HTMLエスケープ
8. [ ] **Insecure Deserialization** - 入力検証
9. [ ] **Using Components with Known Vulnerabilities** - 依存関係更新
10. [ ] **Insufficient Logging & Monitoring** - ログ・監視設定

## 脆弱性スキャン

### バックエンド
```bash
cd backend
# Go依存関係スキャン
go list -json -m all | nancy sleuth
# または
golangci-lint run --enable=gosec
```

### フロントエンド
```bash
cd frontend
# npm 依存関係スキャン
npm audit
npm audit fix
```

## 定期的なチェック

- [ ] 月次: 依存関係の脆弱性スキャン
- [ ] 四半期: セキュリティ監査
- [ ] 年次: ペネトレーションテスト

