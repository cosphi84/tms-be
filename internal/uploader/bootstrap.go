package uploader

import "tms-be/cmd/app"

func UploaderBootstrapApp(app *app.AppConfig) {
	localStorage, err := NewLocalStorage()
	if err != nil {
		panic(err)
	}

	validators := NewUploaderValidators()
	UploaderRepo := NewUploaderRepository(app.DB)
	UploaderSvc := NewUploaderService(UploaderRepo, localStorage, validators)
	UploaderHandler := NewUploaderHandler(UploaderSvc)

	app.UploaderHandler = UploaderHandler
}
