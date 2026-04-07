package counter

import (
	"embed"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// The go:embed directive tells the Go compiler to read the 'migrations' folder
// at compile time and package all SQL schemas natively inside the binary!
//
//go:embed migrations/*.sql
var MigrationsFS embed.FS

// Module wraps all required dependencies, pulling our freshly factored Repositories and Logic 
// underneath a single elegant package umbrella.
type Module struct {
	Repo      *Repository
	Broadcast func([]byte)
}

// NewModule constructs the strict context boundaries
func NewModule(db *pgxpool.Pool, broadcast func([]byte)) *Module {
	return &Module{
		Repo:      &Repository{DB: db},
		Broadcast: broadcast,
	}
}

// RegisterRoutes links the separated handler logic to the parent execution tree
func (m *Module) RegisterRoutes(r chi.Router) {
	r.Get("/api/counter", m.handleGet)
	r.Post("/api/counter/increment", m.handleIncrement)
}
