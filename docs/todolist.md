# todolist

## TODO
* エンティティに CreatedAt と UpdatedAt は必要
* GetByID の owner は pointer にしよう
* 全体的に owner ns.User を使うようにする
* ns.User のコンストラクタ関数（ファクトリ関数）を作る
* interface sqlx.Queryer を使って sqlx.DB と sqlx.Tx を透過的に扱う

## DONE
* テストを簡単にするため、 `postgres.MustGetConn()` が欲しい
