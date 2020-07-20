# todolist

## DOING

## TODO
* test helper `NewRamdomNodeName()`
* postgres への Insert時に CreatedAt と UpdatedAt は払い出させる？
* 全体的に owner ns.User を使うようにする
* ns.User のコンストラクタ関数（ファクトリ関数）を作る
* interface sqlx.Queryer を使って sqlx.DB と sqlx.Tx を透過的に扱う
* 未保存のエンティティは見分けがつく様にしたい 別のエンティティにするべき？

## DONE
* ノードの追加処理
* GetByID の owner は pointer にしよう
* NodeName は、文字列をもとに生成できる: `func Parse(s string)`
* エンティティに CreatedAt と UpdatedAt は必要
* テストを簡単にするため、 `postgres.MustGetConn()` が欲しい
