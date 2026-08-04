package stocks

type IncreaseStockDTO struct {
	ToolsID uint64 `json:"tools_id" binding:"required"`
	SlocID  uint64 `json:"sloc_id" binding:"required"`
	Qty     int    `json:"qty" binding:"required,gt=0"`
	Remarks string `json:"remarks"`
}

type DecreaseStockDTO struct {
	ToolsID uint64 `json:"tools_id" binding:"required"`
	SlocID  uint64 `json:"sloc_id" binding:"required"`
	Qty     int    `json:"qty" binding:"required,gt=0"`
}
