mock:
	mockgen -source=internal/usecases/userrepo/userrepo.go -destination=internal/usecases/userrepo/userrepo_mock.go -package=userrepo UserStore