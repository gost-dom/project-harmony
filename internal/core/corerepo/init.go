package corerepo

import "harmony/internal/couchdb"

var DefaultMessageSource MessageSource
var DefaultDomainEventRepo DomainEventRepository

func init() {
	DB := &couchdb.DefaultConnection
	DefaultDomainEventRepo = DomainEventRepository{DB}
	DefaultMessageSource = MessageSource{DefaultDomainEventRepo, DB}
}
