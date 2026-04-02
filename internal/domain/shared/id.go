package shared

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// ID は全エンティティで共通のID型
// なぜ型エイリアス？シンプルさと型安全性のバランス
type ID string

// NewID はULIDを生成して新しいIDを返す
// なぜULID？ソート可能で一意性が保証される
func NewID() ID {
	ms := ulid.Timestamp(time.Now())
	id, err := ulid.New(ms, rand.Reader)
	if err != nil {
		panic(err)
	}
	return ID(id.String())
}

// Value は文字列としての値を返す
// データベース/API使用のためにプリミティブを返す
func (id ID) Value() string {
	return string(id)
}

// Equals は値オブジェクトの等価性を比較
func (id ID) Equals(other ID) bool {
	return id == other
}

// String はfmt.Stringerインターフェースを実装
func (id ID) String() string {
	return string(id)
}
