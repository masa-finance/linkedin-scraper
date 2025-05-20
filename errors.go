package linkedinscraper

import "errors"

var ErrAuthMissing = errors.New("linkedinscraper: authentication credentials (li_at, csrf_token) are missing")
