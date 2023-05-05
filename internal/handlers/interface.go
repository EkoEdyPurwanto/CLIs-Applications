package handlers

type WikisHandlers interface {
	AddWiki() error
	UpdateWiki() error
	DeleteWiki() error
	GetWiki() error
}
