package query

// BookReadModel - 書籍データのReadModel (CQRSクエリサイド)
//
// これはドメインエンティティではありません。読み取り操作に
// 最適化されたシンプルな構造体です。
type BookReadModel struct {
	ID          string         // 書籍ID
	ISBN        string         // 書籍ISBN
	Title       string         // 書籍タイトル
	Author      string         // 著者
	IsAvailable bool           // loans から導出
	CurrentLoan *LoanReadModel // 貸出中の場合、現在の貸出情報
}

// LoanReadModel - 貸出データのReadModel
type LoanReadModel struct {
	LoanID       string // 貸出ID
	UserID       string // 借りたユーザーID
	BorrowedDate string // 書籍が借りられた日時
	DueDate      string // 返却期限
}
