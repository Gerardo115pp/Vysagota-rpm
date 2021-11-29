module vysagota-accounts

replace (
	github.com/Gerardo115pp/PatriotRouter => /home/el_maligno/.gwen/libs/go/http_router
	github.com/gwen/putils => /home/el_maligno/.gwen/libs/go/putils
	github.com/vysagota/libs => /home/el_maligno/SoftwareProjects/Vysagota/services/Vysagota-User-Librarys
)

go 1.17

require (
	github.com/Gerardo115pp/PatriotRouter v0.0.0-00010101000000-000000000000
	github.com/gwen/putils v0.0.0-00010101000000-000000000000
	github.com/vysagota/libs v0.0.0-00010101000000-000000000000
)
