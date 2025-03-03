package service

import (
	"encoding/json"
	"fmt"
	"os"
)

type Result struct {
	OsiID     int    `json:"osiId"`
	OsiName   string `json:"osiName"`
	ActID     int    `json:"actID"`
	DocID     int    `json:"docID"`
	DocScanID int    `json:"docScanID"`
	DocBase64 string `json:"docBase64"`
	Err       string `json:"err"`
}

func writeToFile(res interface{}) error {
	b, err := json.Marshal(res)
	if err != nil {
		return err
	}
	f, err := os.OpenFile("../../result.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s,\n", string(b))); err != nil {
		return err
	}
	return nil
}

// func TestSignAllActs(t *testing.T) {
// 	cfg := config.New("../.")
// 	cfg.LoggingRequest = false
// 	cfg.LogLevel = "FATAL"
// 	logger := common.InitLogger(cfg)
// 	httpClient := common.InitHTTPClient(cfg, logger)
// 	crypto, err := kalkan.NewKalkanCrypto()
// 	if err != nil {
// 		logger.Fatalf("Error create crypto API: %v", err)
// 	}
// 	defer crypto.Close()

// 	containers := []signer.CAContainer{
// 		{Container: "../ca/GOSTKNCA_8d2269dfe2f64c1de4089dc4cebf3720dd1133c7.p12", Password: "A1234a", Alias: "OSI_STORE"},
// 		{Container: "../ca/GOSTKNCA_KAZAFN.p12", Password: "Aa1234", Alias: "KAZAFN_STORE"},
// 	}
// 	signer := signer.NewSigner(logger, crypto, containers)

// 	coreRepo := repository.NewCoreHTTPRepo(cfg, logger, httpClient)

// 	ctx := context.Background()
// 	osis, err := coreRepo.GetAllOSI(ctx)
// 	require.NoError(t, err)
// 	//fmt.Printf("%+v\n", osis)
// 	for _, osi := range osis {
// 		acts, err := coreRepo.GetSignedActs(ctx, osi.ID)
// 		require.NoError(t, err)

// 		for _, act := range acts {
// 			println(osi.ID, osi.Name, act.StateCode)
// 			if  act.StateCode == "SIGNED" {
// 				docs, err := coreRepo.GetActDocs(ctx, act.ID)
// 				require.NoError(t, err)
// 				for _, doc := range docs {
// 					res := &Result {
// 						OsiID: osi.ID,
// 						OsiName: osi.Name,
// 						ActID: act.ID,
// 						DocID: doc.ID,
// 						DocScanID: doc.Scan.ID,
// 					}
// 					scan, err := coreRepo.GetScan(ctx, doc.Scan.ID)
// 					assert.NoError(t, err)
// 					if err != nil {
// 						res.Err = err.Error()
// 						err = writeToFile(res)
// 						assert.NoError(t, err)
// 						continue
// 					}
// 					res.DocBase64 = scan
// 					//println(scan)

// 					sigDoc, err := signer.SignCMSBase64("OSI_STORE", scan)
// 					//assert.NoError(t, err)
// 					if err != nil {
// 						res.Err = err.Error()
// 						// req := &model.SignActRequest{
// 						// 	ID: act.ID,
// 						// 	Extension: "pdf",
// 						// 	DocBase64: scan,
// 						//   }
// 						//   err = coreRepo.SignAct(ctx, req)
// 						//   assert.NoError(t, err)
// 						//   if err != nil {
// 						// 	  res.Err = err.Error()
// 						//   }
// 						writeToFile(res)
// 					//	assert.NoError(t, err)
// 						continue
// 					}

// 					err = coreRepo.UnsignAct(ctx, act.ID)
// 					assert.NoError(t, err)
// 					if err != nil {
// 						res.Err = err.Error()
// 						err = writeToFile(res)
// 						assert.NoError(t, err)
// 						continue
// 					}

// 					req := &model.SignActRequest{
// 					  ID: act.ID,
// 					  Extension: "pdf",
// 					  DocBase64: sigDoc,
// 					}
// 					err = coreRepo.SignAct(ctx, req)
// 				//	assert.NoError(t, err)
// 					if err != nil {
// 						res.Err = err.Error()
// 						err = writeToFile(res)
// 						assert.NoError(t, err)
// 					}
//                     err = writeToFile(res)
// 					assert.NoError(t, err)
// 				}
// 			}

// 		}
// 	}
// }

// func TestDebtsReport(t *testing.T) {
// 	cfg := config.New("../.")
// 	cfg.LoggingRequest = true
// 	cfg.LogLevel = "FATAL"
// 	logger := common.InitLogger(cfg)
// 	httpClient := common.InitHTTPClient(cfg, logger)
// 	coreRepo := repository.NewCoreHTTPRepo(cfg, logger, httpClient)
// 	appService := NewAppService(cfg, logger, coreRepo, nil, nil, nil)
// 	ctx := auth.NewContext(context.Background(), "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJlZjE0NGJiNy1kZTE5LTQwNDAtOTM1Yy01ZDQwOWFlNmViNmQiLCJodHRwOi8vc2NoZW1hcy5taWNyb3NvZnQuY29tL3dzLzIwMDgvMDYvaWRlbnRpdHkvY2xhaW1zL3VzZXJkYXRhIjoiMTYyIiwiaHR0cDovL3NjaGVtYXMueG1sc29hcC5vcmcvd3MvMjAwNS8wNS9pZGVudGl0eS9jbGFpbXMvbmFtZWlkZW50aWZpZXIiOiI3NzEwMDAwODAwIiwiaHR0cDovL3NjaGVtYXMueG1sc29hcC5vcmcvd3MvMjAwNS8wNS9pZGVudGl0eS9jbGFpbXMvbmFtZSI6ItCS0LDRgdC40LvRjNC10LIg0JDRgNGC0LXQvCIsImh0dHA6Ly9zY2hlbWFzLm1pY3Jvc29mdC5jb20vd3MvMjAwOC8wNi9pZGVudGl0eS9jbGFpbXMvcm9sZSI6WyJDSEFJUk1BTiIsIkFCT05FTlQiLCJPUEVSQVRPUiIsIkFETUlOIl0sImV4cCI6MTY4MjA1ODAyOCwiaXNzIjoiT1NJLkNvcmUiLCJhdWQiOiJPU0kuQ29yZSJ9.-G93jrgoMy9Gq9MEAGINHOX5eYZHQacPJ2wusZX0V4U")

// 	r, err := appService.GetDebtsReport(ctx, &model.OSVReportRequest{
// 		ID:    227,
// 		Begin: time.Date(2023, time.April, 1, 0, 0, 0, 0, time.UTC),
// 		End:   time.Date(2023, time.April, 30, 0, 0, 0, 0, time.UTC),
// 	})
// 	require.NoError(t, err)

// 	f, err := os.Create("data.xlsx")
//     require.NoError(t, err)

//     defer f.Close()

// 	_, err = r.WriteTo(f)
// 	require.NoError(t, err)
// }

// func TestGetAbonentOSVReport(t *testing.T) {
// 		cfg := config.New("../.")
// 		cfg.LoggingRequest = true
// 		cfg.LogLevel = "FATAL"
// 		logger := common.InitLogger(cfg)
// 		httpClient := common.InitHTTPClient(cfg, logger)
// 		coreRepo := repository.NewCoreHTTPRepo(cfg, logger, httpClient)
// 		appService := NewAppService(cfg, logger, coreRepo, nil, nil, nil)
// 		ctx := auth.NewContext(context.Background(), "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiIwOWMzYWI5NS01ZDk0LTQ0MDctYjg3Mi0zMDg4YTI1OTUwNTUiLCJodHRwOi8vc2NoZW1hcy5taWNyb3NvZnQuY29tL3dzLzIwMDgvMDYvaWRlbnRpdHkvY2xhaW1zL3VzZXJkYXRhIjoiMTYyIiwiaHR0cDovL3NjaGVtYXMueG1sc29hcC5vcmcvd3MvMjAwNS8wNS9pZGVudGl0eS9jbGFpbXMvbmFtZWlkZW50aWZpZXIiOiI3NzEwMDAwODAwIiwiaHR0cDovL3NjaGVtYXMueG1sc29hcC5vcmcvd3MvMjAwNS8wNS9pZGVudGl0eS9jbGFpbXMvbmFtZSI6ItCS0LDRgdC40LvRjNC10LIg0JDRgNGC0LXQvCIsImh0dHA6Ly9zY2hlbWFzLm1pY3Jvc29mdC5jb20vd3MvMjAwOC8wNi9pZGVudGl0eS9jbGFpbXMvcm9sZSI6WyJDSEFJUk1BTiIsIkFCT05FTlQiLCJPUEVSQVRPUiIsIkFETUlOIl0sImV4cCI6MTY4NTg4MzQxNiwiaXNzIjoiT1NJLkNvcmUiLCJhdWQiOiJPU0kuQ29yZSJ9.BGLL0UwmDtDhZeuhgxpQLcOizzwM50mr5kZDVJuvRd8")

// 		r, err := appService.GetAbonentOSVReport(ctx, &model.AbonentOSVReportRequest{
// 			ID:    13,
// 			AbonentID: 4906,
// 			Flat: "2  Нежилое",
// 		})
// 		require.NoError(t, err)

// 		f, err := os.Create("data.xlsx")
// 	    require.NoError(t, err)

// 	    defer f.Close()

// 		_, err = r.WriteTo(f)
// 		require.NoError(t, err)
// 	}

// func TestGetAbonentReport(t *testing.T) {
// 		cfg := config.New("../.")
// 		cfg.LoggingRequest = true
// 		cfg.LogLevel = "FATAL"
// 		logger := common.InitLogger(cfg)
// 		httpClient := common.InitHTTPClient(cfg, logger)
// 		coreRepo := repository.NewCoreHTTPRepo(cfg, logger, httpClient)
// 		appService := NewAppService(cfg, logger, coreRepo, nil, nil, nil)
// 		ctx := auth.NewContext(context.Background(), "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJmYzhmNjc4ZS04OGUwLTQwMGEtYTI4NC0zZmNlOWNhY2FiNDAiLCJodHRwOi8vc2NoZW1hcy5taWNyb3NvZnQuY29tL3dzLzIwMDgvMDYvaWRlbnRpdHkvY2xhaW1zL3VzZXJkYXRhIjoiMTYyIiwiaHR0cDovL3NjaGVtYXMueG1sc29hcC5vcmcvd3MvMjAwNS8wNS9pZGVudGl0eS9jbGFpbXMvbmFtZWlkZW50aWZpZXIiOiI3NzEwMDAwODAwIiwiaHR0cDovL3NjaGVtYXMueG1sc29hcC5vcmcvd3MvMjAwNS8wNS9pZGVudGl0eS9jbGFpbXMvbmFtZSI6ItCS0LDRgdC40LvRjNC10LIg0JDRgNGC0LXQvCIsImh0dHA6Ly9zY2hlbWFzLm1pY3Jvc29mdC5jb20vd3MvMjAwOC8wNi9pZGVudGl0eS9jbGFpbXMvcm9sZSI6WyJDSEFJUk1BTiIsIkFCT05FTlQiLCJPUEVSQVRPUiIsIkFETUlOIl0sImV4cCI6MTY4NTg3NDYyNiwiaXNzIjoiT1NJLkNvcmUiLCJhdWQiOiJPU0kuQ29yZSJ9.ZCXVk_TL644uc5GxgU-pNgexOUdYwUUCpI4eSv88NSk")

// 		r, err := appService.GetOSIAbonentsReport(ctx, &model.OSIAbonentsRequest{
// 			ID:    24,
// 		})
// 		require.NoError(t, err)

// 		f, err := os.Create("data.xlsx")
// 	    require.NoError(t, err)

// 	    defer f.Close()

// 		_, err = r.WriteTo(f)
// 		require.NoError(t, err)
// 	}

// func TestGetDebtsReport(t *testing.T) {
// 		cfg := config.New("../.")
// 		cfg.LoggingRequest = true
// 		cfg.LogLevel = "FATAL"
// 		cfg.CoreURL = "http://10.1.1.25:5656"
// 		logger := common.InitLogger(cfg)
// 		httpClient := common.InitHTTPClient(cfg, logger)
// 		coreRepo := repository.NewCoreHTTPRepo(cfg, logger, httpClient)
// 		appService := NewAppService(cfg, logger, coreRepo, nil, nil, nil)
// 		ctx := auth.NewContext(context.Background(), "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiIwZGExMGU3Yi1lMWI5LTRkOTktODEwMi00OTZlOWFiZmIyMjgiLCJodHRwOi8vc2NoZW1hcy5taWNyb3NvZnQuY29tL3dzLzIwMDgvMDYvaWRlbnRpdHkvY2xhaW1zL3VzZXJkYXRhIjoiMTYyIiwiaHR0cDovL3NjaGVtYXMueG1sc29hcC5vcmcvd3MvMjAwNS8wNS9pZGVudGl0eS9jbGFpbXMvbmFtZWlkZW50aWZpZXIiOiI3NzEwMDAwODAwIiwiaHR0cDovL3NjaGVtYXMueG1sc29hcC5vcmcvd3MvMjAwNS8wNS9pZGVudGl0eS9jbGFpbXMvbmFtZSI6ItCS0LDRgdC40LvRjNC10LIg0JDRgNGC0LXQvCIsImh0dHA6Ly9zY2hlbWFzLm1pY3Jvc29mdC5jb20vd3MvMjAwOC8wNi9pZGVudGl0eS9jbGFpbXMvcm9sZSI6WyJBRE1JTiIsIk9QRVJBVE9SIiwiQ0hBSVJNQU4iLCJBQk9ORU5UIl0sImV4cCI6MTY5MTU4MTI2MCwiaXNzIjoiT1NJLkNvcmUiLCJhdWQiOiJPU0kuQ29yZSJ9.egakompG8N7t9itNicHOuSkDlTK0Dn9JmGYqlcgsJYo")

// 		r, err := appService.GetDebtsReport(ctx, &model.OSVReportRequest{
// 			ID:    71,
// 			Begin: time.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
// 		    End:   time.Date(2023, time.August, 31, 0, 0, 0, 0, time.UTC),
// 		    ForAbonent:  false,

// 		})
// 		require.NoError(t, err)

// 		f, err := os.Create("data.xlsx")
// 	    require.NoError(t, err)

// 	    defer f.Close()

// 		_, err = r.WriteTo(f)
// 		require.NoError(t, err)
// 	}
