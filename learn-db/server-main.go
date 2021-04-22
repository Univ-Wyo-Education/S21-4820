package main

// Copyright (C) Philip Schlump 2016-2020.

// https://fabianlee.org/2017/05/21/golang-running-a-go-binary-as-a-systemd-service-on-ubuntu-16-04/

// xyzzy88 - un/pw/register - add a quick table.

// xyzzy501101 ./server-main --tool encrypt-unencrypted-accounts

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
	template "github.com/pschlump/textTemplate"
	"gitlab.com/pschlump/PureImaginationServer/CliResponseWriter"
	pgsave "gitlab.com/pschlump/PureImaginationServer/PgSave"
	"gitlab.com/pschlump/PureImaginationServer/ReadConfig"
	"gitlab.com/pschlump/PureImaginationServer/apache_logger"
	"gitlab.com/pschlump/PureImaginationServer/auth_check"
	"gitlab.com/pschlump/PureImaginationServer/awss3"
	"gitlab.com/pschlump/PureImaginationServer/cors"
	"gitlab.com/pschlump/PureImaginationServer/email"
	"gitlab.com/pschlump/PureImaginationServer/ethProc"
	"gitlab.com/pschlump/PureImaginationServer/get"
	"gitlab.com/pschlump/PureImaginationServer/pstate"
	"gitlab.com/pschlump/PureImaginationServer/ymux"
)

// ----------------------------------------------------------------------------------
// Notes:
//   Graceful Shutdown: https://stackoverflow.com/questions/39320025/how-to-stop-http-listenandserve
//   Email with HTML and Text: https://stackoverflow.com/questions/9950098/how-to-send-an-email-using-go-with-an-html-formatted-body
//   Image Paint on Image : https://blog.golang.org/go-imagedraw-package
// ----------------------------------------------------------------------------------

var Cfg = flag.String("cfg", "cfg.json", "config file for this call")
var Cli = flag.String("cli", "", "Run as a CLI command intead of a server")
var HostPort = flag.String("hostport", ":8003", "Host/Port to listen on") // q8s
var DbFlagParam = flag.String("db_flag", "", "Additional Debug Flags")
var TLS_crt = flag.String("tls_crt", "", "TLS Signed Publick Key")
var TLS_key = flag.String("tls_key", "", "TLS Signed Private Key")
var Version = flag.Bool("version", false, "Report version of code and exit")
var ChkTables = flag.Bool("chk-tables", false, "Chack table structre and exit")
var Comment = flag.String("comment", "", "Unused comment for ps.")
var CdTo = flag.String("CdTo", "", "Change Directory To.")

type QRGeneration struct {
	// Data used in generating / processing the QR codes
	LnLocal             string `json:"ln_local" default:"./www/"`
	MinQRBeforeGenerate int    `json:"min_qr_before_generate" default:"150"`
	BaseQRUrl           string `json:"base_qr_url" default:"http://q8s.co/"`
	BaseServerUrl       string `json:"base_server_url" default:"http://www.q8s.com"` // urlBase := "http://www.q8s.com"		// xyzzy441 - tempalte {{.scheme}}
	NameOnImage         string `json:"name_on_image" default:"QR Today"`
	ImageTag            string `json:"image_tag" default:"./ex-60x60e1.png"`
	ArialFont           string `json:"arial_font" default:"/Library/Fonts/Arial.ttf"`
}

type GlobalConfigData struct {
	ymux.BaseConfigType

	ethProc.EthConfigData

	pstate.PstateConfigType

	QRGeneration

	ExpiredConnectionThreshold int `json:"expired_connection_threshold" defauilt:"50000"`
}

var n_ticks = 0
var gCfg GlobalConfigData
var GitCommit string
var ch chan string
var timeout chan string
var logFilePtr *os.File
var DB *sql.DB
var DbOn map[string]bool
var wg sync.WaitGroup
var logger *log.Logger
var shutdownWaitTime = time.Duration(1)
var isTLS bool
var curDir string
var LastGen string = ""
var awsSession *session.Session
var httpServerList []*http.Server

func init() {
	isTLS = false
	ch = make(chan string, 1)
	timeout = make(chan string, 2)
	DbOn = make(map[string]bool)
	logger = log.New(os.Stdout, "", 0)
	template.SetNoValue("", true)
}

func main() {

	flag.Parse() // Parse CLI arguments to this, --cfg <name>.json

	fns := flag.Args()
	if *Cli != "" {
		ymux.SetCliOpts(Cli, fns)
	} else if len(fns) != 0 {
		fmt.Printf("Extra arguments are not supported [%s]\n", fns)
		os.Exit(1)
	}

	if *CdTo != "" {
		err := os.Chdir(*CdTo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Uable to chagne to %s directory, error:%s\n", *CdTo, err)
			os.Exit(1)
		}
	}

	if *Version {
		fmt.Printf("Version (Git Commit): %s\n", GitCommit)
		os.Exit(0)
	}

	if Cfg == nil {
		fmt.Printf("--cfg is a required parameter\n")
		os.Exit(1)
	}

	// ------------------------------------------------------------------------------
	// Read in Configuration
	// ------------------------------------------------------------------------------
	err := ReadConfig.ReadFile(*Cfg, &gCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read confguration: %s error %s\n", *Cfg, err)
		os.Exit(1)
	}
	ReadConnectinSetup("./db_cfg.json")
	if DbOn["List.Conn"] {
		fmt.Printf("List Conn: ->%s<-\n", godebug.SVarI(G_ConnPool))
	}

	// fmt.Printf("Config: %s\n", godebug.SVarI(G_ConnPool))
	// os.Exit(1)
	fmt.Printf("cfg.json/config: %s\n", godebug.SVarI(gCfg))

	// ------------------------------------------------------------------------------
	// Logging File
	// ------------------------------------------------------------------------------
	logFilePtr = ymux.LoggingSetup(&(gCfg.BaseConfigType))

	// ------------------------------------------------------------------------------
	// Debug Flag Processing
	// ------------------------------------------------------------------------------
	ymux.DebugFlagProcess(DbFlagParam, DbOn, &(gCfg.BaseConfigType))
	ymux.SetDbFlag(DbOn)
	CliResponseWriter.SetDbFlag(DbOn)
	get.SetupGet(DbOn)

	// fmt.Printf("%s\n", godebug.SVarI(DbOn))

	// ------------------------------------------------------------------------------
	// TLS parameter check / setup
	// ------------------------------------------------------------------------------
	ymux.TLSSetup(TLS_crt, TLS_key, &(gCfg.BaseConfigType))

	// ------------------------------------------------------------------------------
	// Get the current directory.
	// ------------------------------------------------------------------------------
	curDir, err = os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get current working directory: %s\n", err)
		os.Exit(1)
	}

	// ------------------------------------------------------------------------------
	// Connect to PG
	// ------------------------------------------------------------------------------
	DB = ymux.ConnectToPG(&(gCfg.BaseConfigType))
	ymux.Ping()

	// ------------------------------------------------------------------------------
	// Test database creation, usage.
	// ------------------------------------------------------------------------------
	// C1, e0 := GetConn("nameBob")
	// C2, e1 := GetConn("nameJane")

	// fmt.Printf("e0=%s e1=%s\n", e0, e1)

	// func SQLInsertDB( DB *sql.DB, stmt string, data ...interface{}) (err error) {
	// e0 = ymux.SQLInsertDB(C1, "create table bob ( isC1 text )")
	// e1 = ymux.SQLInsertDB(C2, "create table bob ( isC2 text )")
	// fmt.Printf("e0=%s e1=%s\n", e0, e1)
	// fmt.Printf("e0=%s e1=%s\n", e0, e1)
	// os.Exit(1)

	ymux.SetupCrud(logFilePtr, DbOn)
	auth_check.Setup(&(gCfg.BaseConfigType), DbOn, logFilePtr, gCfg.QRGeneration.BaseServerUrl, GitCommit) // func Setup(gcfg *ymux.BaseConfigType, dbf map[string]bool, f *os.File, base string, istls bool, gitCommit string) {
	// auth_check.SetupCall()
	pstate.Setup(&(gCfg.BaseConfigType), DbOn, logFilePtr, GitCommit, &(gCfg.PstateConfigType))
	email.SetupEmail(&(gCfg.BaseConfigType), DbOn, logFilePtr)

	if *Cli == "" || DbOn["pg"] {
		fmt.Fprintf(os.Stderr, "%sConnected to PG%s\n", MiscLib.ColorGreen, MiscLib.ColorReset)
	}

	// ------------------------------------------------------------------------------
	// Check that the DDL / Table / Sequences are in PG and ready to use.
	// ------------------------------------------------------------------------------
	good, bad := ValidateTablesDDL(false)

	if *Cli == "" || DbOn["pg"] {
		fmt.Fprintf(os.Stderr, "%sPG Validated%s\n", MiscLib.ColorGreen, MiscLib.ColorReset)
	}
	if *ChkTables {
		fmt.Fprintf(os.Stderr, "%sPG Tables Passed: %s%s\n", MiscLib.ColorGreen, good, MiscLib.ColorReset)
		if len(bad) > 0 {
			fmt.Fprintf(os.Stderr, "%sPG Tables Failed: %s%s\n", MiscLib.ColorGreen, bad, MiscLib.ColorReset)
		} else {
			fmt.Fprintf(os.Stderr, "%sPG All Tables passed:%s\n", MiscLib.ColorGreen, MiscLib.ColorReset)
		}
		// --------------------------------------------------------------------------------------------------
		// --------------------------------------------------------------------------------------------------
		// This is the place to do a migration (upgrade) of tables, stored procs etc.
		// xyzzy
		// TODO
		// xyzzy
		// --------------------------------------------------------------------------------------------------
		// --------------------------------------------------------------------------------------------------
		fmt.Fprintf(os.Stderr, "%sPG Validated - tables checked early exit%s\n", MiscLib.ColorGreen, MiscLib.ColorReset)
		os.Exit(0)
	}

	// ------------------------------------------------------------------------------
	// Set config in d.b. based on ./cfg.json file
	// ------------------------------------------------------------------------------

	// -------------------------------------------- TODO xyzzy - Temp Comment Out -------------------------------------------------------
	// -------------------------------------------- TODO xyzzy - Temp Comment Out -------------------------------------------------------
	// -------------------------------------------- TODO xyzzy - Temp Comment Out -------------------------------------------------------
	// -------------------------------------------- TODO xyzzy - Temp Comment Out -------------------------------------------------------
	// -------------------------------------------- TODO xyzzy - Temp Comment Out -------------------------------------------------------
	if false {
		fmt.Printf("AT: %s gCfg = --->>>\n%s<<<---\n", godebug.LF(), godebug.SVarI(gCfg))
		pgsave.PgSaveStruct(&gCfg)
	}

	// ------------------------------------------------------------------------------
	// Announce if Encrypted
	// ------------------------------------------------------------------------------
	// Cfg.LogFileEncryptionKey string `json:"log_file_encryption_key" default:"$ENV$LOG_ENCRYPTION_KEY"`
	// Cfg.DB_Enc_Key string `json:"DB_Enc_Key" default:"$ENV$DB_ENC_KEY"`
	// Cfg.DB_IV_Data string `json:"DB_IV_Data" default:"$ENV$DB_ENC_IV_DATA"`

	if gCfg.LogFileEncryptionKey != "" {
		fmt.Printf("\n%sLog File Encryption is enabled.%s\n", MiscLib.ColorBlueOnWhite, MiscLib.ColorReset)
	}
	if gCfg.DB_Enc_Key != "" {
		fmt.Printf("\n%sDatabase User Info Encryption is enabled.%s\n", MiscLib.ColorBlueOnWhite, MiscLib.ColorReset)
		if gCfg.DB_IV_Data == "" {

			fmt.Printf("%sError - missing IV in : User Info Encryption is enabled.%s\n", MiscLib.ColorRed, MiscLib.ColorReset)
		}
		fmt.Printf("\n")
	}

	// ------------------------------------------------------------------------------
	// Setup connection to S3
	// ------------------------------------------------------------------------------
	if DbOn["push-to-aws"] {
		awsSession = awss3.SetupS3(gCfg.S3_bucket, gCfg.S3_region, DbOn)
	}

	// ------------------------------------------------------------------------------
	// Setup HTTP End Points
	// ------------------------------------------------------------------------------
	mux := ymux.NewServeMux(logFilePtr, &(gCfg.BaseConfigType))
	mux.BaseServerUrl = gCfg.QRGeneration.BaseServerUrl // {{.scheme}}

	if false {
		// mux.Handle("/api/v1/", http.HandlerFunc(cors.CreateHandleCORS(cors.AnySite))).Method('OPTIONS').DocTag("<h2>/api/v1/statusxyzzy1")
		cors.SetupMux("/api/v1", mux)
	}

	mux.Handle("/api/v1/status", http.HandlerFunc(auth_check.HandleStatus)).DocTag("<h2>/api/v1/status")
	mux.Handle("/api/status", http.HandlerFunc(auth_check.HandleStatus)).DocTag("<h2>/api/status")
	mux.Handle("/status", http.HandlerFunc(auth_check.HandleStatus)).Method("GET", "POST").DocTag("<h2>/status")
	mux.Handle("/api/v1/config", http.HandlerFunc(HandleConfig)).Method("GET").AuthRequired().DocTag("<h2>/api/v1/config")
	mux.Handle("/api/v1/kick", http.HandlerFunc(HandleKick)).Method("GET", "POST").AuthRequired().DocTag("<h2>/api/v1/kick")
	mux.Handle("/api/v1/exit-server", http.HandlerFunc(HandleExitServer)).Method("GET", "POST").AuthRequired().DocTag("<h2>/api/v1/exit-server")

	mux.Handle("/api/v1/run-sql-in-db", http.HandlerFunc(HandleRunSQLInDatabase)).DocTag("<h2>/api/v1/statusxyzzy1").AuthRequired()
	mux.Handle("/api/v1/user-desc-table", http.HandlerFunc(HandleDescTable)).DocTag("<h2>/api/v1/statusxyzzy1").AuthRequired()

	// For the list of end-points (URI Paths) see ./handle.go
	ymux.HandleTables(mux, &(gCfg.BaseConfigType))

	mux.Handle("/desc.html", http.HandlerFunc(mux.Desc())).NoDoc() // Documentation for API (and test run code) v2

	// ----------------------------------------------------------------------------------------------------------------------------------------------
	// user authentication
	// ----------------------------------------------------------------------------------------------------------------------------------------------
	auth_check.SetupMux(mux)

	// ----------------------------------------------------------------------------------------------------------------------------------------------
	// pstate setup mux
	// ----------------------------------------------------------------------------------------------------------------------------------------------
	pstate.SetupMux(mux)

	// ----------------------------------------------------------------------------------------------------------------------------------------------
	// Serve Files (now with auth)
	// Files are not cached, 1st set requires login auth in the ./www/auth-req directory.
	// ----------------------------------------------------------------------------------------------------------------------------------------------
	StripAndServeFiles := func(prefix string) http.HandlerFunc {
		if prefix == "" {
			fmt.Fprintf(os.Stderr, "Fatal Error - prefix can not be empty ->%s<- at:%s", prefix, godebug.LF())
			os.Exit(1)
		}
		Wd, err := os.Getwd()
		if err != nil {
			Wd = os.Getenv("PWD")
		}
		Dir := path.Join(Wd, gCfg.StaticPath, "auth-req") + "/"
		godebug.DbPf(DbOn["auth-req-file"], "Path To Server From ->%s<-\n", Dir, err)
		fs := http.FileServer(http.Dir(Dir))
		return http.HandlerFunc(func(www http.ResponseWriter, req *http.Request) {
			if pth := strings.TrimPrefix(req.URL.Path, prefix); len(pth) < len(req.URL.Path) {
				godebug.DbPf(DbOn["auth-req-file"], "AT: %s\n", godebug.LF())
				reqCpy := new(http.Request)
				*reqCpy = *req
				reqCpy.URL = new(url.URL)
				*reqCpy.URL = *req.URL
				reqCpy.URL.Path = pth
				www.Header().Set("Cache-Control", "public, max-age=1")
				fs.ServeHTTP(www, reqCpy)
			} else {
				godebug.DbPf(DbOn["auth-req-file"], "AT: %s\n", godebug.LF())
				http.NotFound(www, req)
			}
		})
	}

	mux.Handle("/auth-req/", http.HandlerFunc(StripAndServeFiles("/auth-req/"))).NoDoc().NoValidate().AuthRequired()

	{
		Dir := gCfg.StaticPath
		fs := http.FileServer(http.Dir(Dir))
		fx := func(www http.ResponseWriter, req *http.Request) {
			www.Header().Set("Cache-Control", "public, max-age=1")
			fs.ServeHTTP(www, req)
		}
		mux.Handle("/", http.HandlerFunc(fx)).NoDoc().NoValidate()
	}

	mux.Compile()

	// ------------------------------------------------------------------------------
	// Main Processing
	// ------------------------------------------------------------------------------
	go TimedDispatch()

	// ticker on channel - send once a minute
	go func(n int) {
		for {
			time.Sleep(time.Duration(n) * time.Second)
			timeout <- "timeout"
			n_ticks++
		}
	}(gCfg.TickerSeconds)

	//	go func() {
	//		time.Sleep(1 * time.Second)
	//		timeout <- "timeout"
	//		n_ticks++
	//	}()

	// ------------------------------------------------------------------------------
	// Setup signal capture
	// ------------------------------------------------------------------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// ------------------------------------------------------------------------------
	// Live Monitor Setup
	// ------------------------------------------------------------------------------
	ymux.LiveMonSetup("Auth-Demo-Service", DbOn, logFilePtr, &(gCfg.BaseConfigType))

	// ------------------------------------------------------------------------------
	// ------------------------------------------------------------------------------
	if *Cli != "" {
		www := CliResponseWriter.NewCliResonseWriter() // www := http.ResponseWriter
		/*
		   type url.URL struct {
		   	Scheme     string
		   	Opaque     string    // encoded opaque data
		   	User       *Userinfo // username and password information
		   	Host       string    // host or host:port
		   	Path       string    // path (relative paths may omit leading slash)
		   	RawPath    string    // encoded path hint (see EscapedPath method)
		   	ForceQuery bool      // append a query ('?') even if RawQuery is empty
		   	RawQuery   string    // encoded query values, without '?'
		   	Fragment   string    // fragment for references, without '#'
		   }
		   // From: https://golang.org/src/net/url/url.go?s=9736:10252#L353 :363
		*/
		qryParam := ymux.GenQryFromCli()
		if DbOn["cli"] {
			fmt.Printf("qry_params= ->%s<- at:%s\n", qryParam, godebug.LF())
		}
		u := url.URL{
			User:     nil,
			Host:     "127.0.0.1:80",
			Path:     *Cli,
			RawQuery: qryParam,
		}
		req := &http.Request{ // https://golang.org/src/net/http/request.go:113
			Method:     "GET",
			URL:        &u, // *url.URL
			Proto:      "HTTP/1.0",
			ProtoMajor: 1,
			ProtoMinor: 0,
			Header:     make(http.Header),
			// Body io.ReadCloser // :182 -- not used, GET request - no body.
			// Form url.Values -- Populate with values from CLI
			Host:       "127.0.0.1:80",
			RequestURI: *Cli + "?" + qryParam, // "RequestURI": "/api/v1/status?id=dump-request",
		}
		switch *Cli {
		case "/api/v1/status":
			auth_check.HandleStatus(www, req)
		case "/api/v1/kick":
			HandleKick(www, req)
		case "/api/v1/exit-server":
			fmt.Printf("Exit server\n")

		case "/desc.html":
			x := mux.Desc()
			x(www, req)
		case "/image/lclr.gif":
			x := mux.GenerateAnalytics()
			x(www, req)

		default:
			fn := "./www/" + *Cli
			s, err := ioutil.ReadFile(fn)
			if err != nil {
				fmt.Printf("Status: 404, Request [%s]\n", *Cli)
			} else {
				fmt.Printf("Status: 200\n")
				fmt.Printf("%s\n", s)
			}
		}
		www.Flush()
		if DbOn["Cli.Where"] {
			www.DumpWhere()
		}
		os.Exit(0)
	}

	// Split host:port into a set of values based on comma.  This is so you can listen to a list
	// of addresses.  If you want to setup TLS (https, wss) then you need
	// *HostPort == "https://name:port,https://www.name1.com:port" etc.
	hpList := strings.Split(*HostPort, ",")
	for _, hp := range hpList {

		if strings.HasPrefix(hp, "https://") {
			isTLS = true
			hp = hp[len("https://"):]
		} else if strings.HasPrefix(hp, "wss://") {
			isTLS = true
			hp = hp[len("wss://"):]
		}

		var httpServer *http.Server
		fmt.Printf("%sSetup Server: hp=[%s]%s\n", MiscLib.ColorGreen, hp, MiscLib.ColorReset)

		// ------------------------------------------------------------------------------
		// Setup / Run the HTTP Server.
		// ------------------------------------------------------------------------------
		// authHandler := auth_check.NewAuthCheckerHandler(mux, mux, logFilePtr, &(gCfg.BaseConfigType))
		// Note: https://ieftimov.com/post/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/
		pstateHandler := pstate.NewPageStateHandler(mux, mux, logFilePtr, &(gCfg.BaseConfigType))
		authHandler := auth_check.NewAuthCheckerHandler(pstateHandler, mux, logFilePtr, &(gCfg.BaseConfigType))
		loggingHandler := apache_logger.NewApacheLoggingHandler(authHandler, logFilePtr, nil, &apache_logger.DBACleanupFunc)
		if isTLS {
			cfg := &tls.Config{
				MinVersion:               tls.VersionTLS12,
				CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
				PreferServerCipherSuites: true,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				},
			}
			httpServer = &http.Server{
				Addr:              hp,
				Handler:           loggingHandler,
				ReadTimeout:       5 * time.Second,
				WriteTimeout:      5 * time.Second,
				IdleTimeout:       90 * time.Second,
				ReadHeaderTimeout: 10 * time.Second,
				TLSConfig:         cfg,
				TLSNextProto:      make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
			}
		} else {
			httpServer = &http.Server{
				Addr:              hp,
				Handler:           loggingHandler,
				ReadTimeout:       5 * time.Second,
				WriteTimeout:      5 * time.Second,
				IdleTimeout:       90 * time.Second,
				ReadHeaderTimeout: 10 * time.Second,
			}
		}

		go func(httpServer *http.Server, hp string) {
			httpServerList = append(httpServerList, httpServer)
			wg.Add(1)
			defer wg.Done()
			if isTLS {
				fmt.Fprintf(os.Stderr, "%sListening on https://%s%s\n", MiscLib.ColorGreen, hp, MiscLib.ColorReset)
				if err := httpServer.ListenAndServeTLS(*TLS_crt, *TLS_key); err != nil {
					logger.Fatal(err)
				}
			} else {
				fmt.Fprintf(os.Stderr, "%sListening on http://%s%s\n", MiscLib.ColorGreen, hp, MiscLib.ColorReset)
				if err := httpServer.ListenAndServe(); err != nil {
					logger.Fatal(err)
				}
			}
		}(httpServer, hp)

	}

	// ------------------------------------------------------------------------------
	// Catch signals from [Control-C]
	// ------------------------------------------------------------------------------
	select {
	case <-stop:
		fmt.Fprintf(os.Stderr, "\nShutting down the server... Received OS Signal...\n")
		fmt.Fprintf(logFilePtr, "\nShutting down the server... Received OS Signal...\n")
		for _, httpServer := range httpServerList {
			ctx, cancel := context.WithTimeout(context.Background(), shutdownWaitTime*time.Second)
			defer cancel()
			err := httpServer.Shutdown(ctx)
			if err != nil {
				fmt.Printf("Error on shutdown: [%s]\n", err)
			}
		}
	}

	// ------------------------------------------------------------------------------
	// Wait for HTTP server to exit.
	// ------------------------------------------------------------------------------
	wg.Wait()
}

// TimedDispatch waits for a "kick" or a timeout and calls QrGenerate forever.
func TimedDispatch() {
	for {
		select {
		case <-ch:
			fmt.Printf("%sChanel Activated\n%s", MiscLib.ColorCyan, MiscLib.ColorReset)
			CloseExpiredConnections(gCfg.ExpiredConnectionThreshold, n_ticks)
			// QrGenerate()

		case <-timeout:
			ts := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("%sTimeout at %s\n%s", MiscLib.ColorYellow, ts, MiscLib.ColorReset)
			CloseExpiredConnections(gCfg.ExpiredConnectionThreshold, n_ticks)
			// QrGenerate()
		}
	}
}

func GetCurTick() int {
	return n_ticks
}

/* vim: set noai ts=4 sw=4: */
