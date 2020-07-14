# todolist

## TODO
* NodeName は、文字列をもとに生成できる: `func Parse(s string)`
* GetByID の owner は pointer にしよう
* 全体的に owner ns.User を使うようにする
* ns.User のコンストラクタ関数（ファクトリ関数）を作る
* interface sqlx.Queryer を使って sqlx.DB と sqlx.Tx を透過的に扱う
* 未保存のエンティティは見分けがつく様にしたい 別のエンティティにするべき？

## DOING
* エンティティに CreatedAt と UpdatedAt は必要

## DONE
* テストを簡単にするため、 `postgres.MustGetConn()` が欲しい
