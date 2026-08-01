package offices

import "github.com/gin-gonic/gin"

// RegisterPublicRoutes sengaja kosong -- gak ada endpoint offices yang
// boleh diakses tanpa login. Method ini tetap ada (bukan dihapus) biar
// kontraknya konsisten sama module lain (auth, dst), dipanggil dari
// internal/app/route.go walau isinya no-op.
func (h *Handler) RegisterPublicRoutes(g *gin.RouterGroup) {}

func (h *Handler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	offices := rg.Group("/offices")
	{
		offices.GET("", h.List)                  // ?level=cabang|sdss|ssr
		offices.GET("/tree", h.Tree)             // semua level, nested, buat select FE
		offices.GET("/:id/children", h.Children) // anak langsung dari :id
		offices.POST("", h.Create)
		offices.PUT("/:id", h.Update)
		offices.DELETE("/:id", h.Delete)
	}
}
