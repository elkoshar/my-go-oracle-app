package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"oracle.com/oracle/my-go-oracle-app/infra/database/sql"
)

type HealthChecker struct {
	Master sql.MasterDB
	Slave  sql.SlaveDB
}

func (h *HealthChecker) health(ctx context.Context) map[string]interface{} {
	OK := "OK"
	FAILED := "FAILED"

	applicationStatus := OK
	mysqlStatus := OK
	mysqlSlaveStatus := OK
	redisStatus := OK

	err := h.Master.Ping()
	if err != nil {
		mysqlStatus = FAILED
	}

	err = h.Slave.Ping()
	if err != nil {
		mysqlSlaveStatus = FAILED
	}

	resp := map[string]interface{}{
		"name": os.Args[0],
		"status": map[string]string{
			"application": applicationStatus,
			"mysql":       mysqlStatus,
			"mysqlSlave":  mysqlSlaveStatus,
			"redis":       redisStatus,
		},
	}

	return resp
}

func (h *HealthChecker) HealthChi(w http.ResponseWriter, r *http.Request) {
	resp := h.health(r.Context())
	data, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *HealthChecker) HealthEcho(c echo.Context) error {
	resp := h.health(c.Request().Context())
	return c.JSON(200, resp)
}
