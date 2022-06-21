package app

import "toto-config-api/internal/skuconfig/app/query"

type Application struct {
	Queries  Queries
	Commands Commands
}

type Commands struct {
}

type Queries struct {
	SKUForConfig query.SKUForConfigHandler
}
