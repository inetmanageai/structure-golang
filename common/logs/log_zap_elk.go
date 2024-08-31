package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	cf "structure-golang/config"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Implement ADAPTER
type appLogsElk struct {
	log    *zap.Logger
	client *elasticsearch.Client
	index  string
}

type LogsElasticModel struct {
	Level      string    `json:"level"`
	Time       time.Time `json:"time"`
	LoggerName string    `json:"logger_name"`
	Message    string    `json:"message"`
	Caller     struct {
		Defined  bool    `json:"defined"`
		PC       uintptr `json:"pc"`
		File     string  `json:"file"`
		Line     int     `json:"line"`
		Function string  `json:"function"`
	} `json:"caller"`
}

func NewAppLogsElk(client *elasticsearch.Client) AppLog {
	var log *zap.Logger

	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = ""

	var err error
	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	return appLogsElk{
		log:    log,
		client: client,
		index:  cf.Env.ElasticIndex,
	}
}

func (l appLogsElk) Info(msg string) {
	// get logs value
	ce := l.log.Check(0, msg)
	_ce, _ := json.Marshal(ce)

	// create elastic log
	l.indexElasticLog(_ce)

	// logs in terminal
	l.log.Info(msg)
}

func (l appLogsElk) Debug(msg string) {
	// get logs value
	ce := l.log.Check(-1, msg)
	_ce, _ := json.Marshal(ce)

	// create elastic log
	l.indexElasticLog(_ce)

	// logs in terminal
	l.log.Debug(msg)
}

func (l appLogsElk) Warning(msg string) {
	// get logs value
	ce := l.log.Check(1, msg)
	_ce, _ := json.Marshal(ce)

	// create elastic log
	l.indexElasticLog(_ce)

	// logs in terminal
	l.log.Warn(msg)
}

func (l appLogsElk) Error(msg interface{}) {
	switch v := msg.(type) {
	case error:
		// logs in terminal
		l.log.Error(v.Error())

		// get logs value
		ce := l.log.Check(2, v.Error())
		_ce, _ := json.Marshal(ce)

		// create elastic log
		l.indexElasticLog(_ce)
	case string:
		// logs in terminal
		l.log.Error(v)

		// get logs value
		ce := l.log.Check(2, v)
		_ce, _ := json.Marshal(ce)

		// create elastic log
		l.indexElasticLog(_ce)
	}
}

func (l appLogsElk) indexElasticLog(data []byte) error {
	// edit payload body elastic
	payload := LogsElasticModel{}
	json.Unmarshal(data, &payload)
	body, _ := json.Marshal(payload)

	// create index name
	year, week := time.Now().ISOWeek()
	indexName := l.index + fmt.Sprintf("-%vweek%v", year, week)

	// Send logs to Elastic
	go func() {
		l.client.Index(indexName, bytes.NewReader(body))
	}()

	return nil
}
