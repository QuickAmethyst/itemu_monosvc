package model

type Uom struct {
	ID          int64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Decimal     int32  `json:"decimal"`
}

type UomsInput struct {
	Paging *PagingInput `json:"paging"`
}

type UomsResult struct {
	Data   []*Uom  `json:"data"`
	Paging *Paging `json:"paging"`
}

type WriteUomInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Decimal     int32  `json:"decimal"`
}
