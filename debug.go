package restflix

//var debugOperationMethod = "[POST]-accounts-google-directory-users"

//var debugOperationMethod = "[GET]-agents-exports"
var debugOperationMethod = "[POST]-agents-clickmaps-list"

// POST /accounts/google-directory/users -> users any => DONE
// POST /agents -> agents any => DONE
// GET /agents/auth/google => NOT_SUPPORTED
// GET /agents/auth/google/callback => NOT_SUPPORTED
// POST /agents/auth/login => NOT_SUPPORTED (nested function response)
// POST /agents/auth/logout => TODO: iris.Map{}
// POST /agents/clickmaps/element => TODO
// POST /agents/clickmaps/list => TODO
// GET /agents/exports => BUG (custom types)
// POST /agents/filters/cardinality => PARTIALLY (custom types)
// POST /accounts/password/set => BUG (nested selector)
// POST /agents/batch => BUG ([]error is parsed to `{"error": []}`)
// PUT /agents/change-password => BUG (nested selector)
// POST /agents/filters/cardinality => BUG (custom types, make(map[string]uint64))
// POST /agents/funnels => BUG (array of objects)
// PUT /agents/funnels => BUG (array of objects)
// POST /agents/funnels/compute => BUG (custom types)
// PUT /agents/websites/:website_id/recording-elements => BUG (type from variable statement (function return), nested selector)
