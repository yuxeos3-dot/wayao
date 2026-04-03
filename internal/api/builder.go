package api

import "net/http"

func (app *App) HandleBuild(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	if app.BuildFunc == nil {
		jsonError(w, 500, "build engine not configured")
		return
	}
	if err := app.BuildFunc(id); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonOK(w, "build complete")
}

func (app *App) HandleDeploy(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	if app.DeployFunc == nil {
		jsonError(w, 500, "deploy engine not configured")
		return
	}
	if err := app.DeployFunc(id); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonOK(w, "deploy complete")
}

func (app *App) BuildSiteByID(id int64) error {
	if app.BuildFunc != nil {
		return app.BuildFunc(id)
	}
	return nil
}

func (app *App) DeploySiteByID(id int64) error {
	if app.DeployFunc != nil {
		return app.DeployFunc(id)
	}
	return nil
}
