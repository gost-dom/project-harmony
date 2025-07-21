package corerepo

var DefaultMessageSource MessageSource
var DefaultDomainEventRepo DomainEventRepository

func init() {
	DB := &DefaultConnection
	DefaultDomainEventRepo = DomainEventRepository{DB}
	DefaultMessageSource = MessageSource{DefaultDomainEventRepo, DB}
}
