package controllers

import "net/http"

type Adapter func(http.Handler) http.Handler
