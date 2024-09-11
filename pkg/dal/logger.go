package dal

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

func NewLogger(name string, database string) *PostgresZap {
	return &PostgresZap{
		dbName: database,
		l:      *zap.S().Named(name).With("database", database),
	}
}

type PostgresZap struct {
	IgnorerecordNotFoundError bool
	SlowThreshold             time.Duration
	dbName                    string
	l                         zap.SugaredLogger
}

func (p *PostgresZap) LogMode(level logger.LogLevel) logger.Interface {
	// not gorms responsibility here we're already managing this
	return p
}

/*
{"level":"error","tz":"2024-09-03T23:36:46.534-0700","logger":"emri","caller":"dal/logger.go:25",
"msg":"gorm message",
"pid":3376936,
"fmt":"failed to initialize database, got error %v",
"data":"/mnt/d/development/emri/pkg/dal/dal.go:30%!(EXTRA *pgconn.ConnectError=failed to connect to `host=10.0.0.210 user=admin@admin.com database=emrides`: failed SASL auth (FATAL: password authentication failed for user \"admin@admin.com\" (SQLSTATE 28P01)))","stacktrace":"github.com/asciifaceman/emri/pkg/dal.(*PostgresZap).Error\n\t/mnt/d/development/emri/pkg/dal/logger.go:25\ngorm.io/gorm.Open\n\t/home/ccorbett/go/1.21.5/pkg/mod/gorm.io/gorm@v1.25.11/gorm.go:214\ngithub.com/asciifaceman/emri/pkg/dal.(*PG).Connect\n\t/mnt/d/development/emri/pkg/dal/dal.go:30\ngithub.com/asciifaceman/emri/pkg/dal.New\n\t/mnt/d/development/emri/pkg/dal/dal.go:17\ngithub.com/asciifaceman/emri/cmd.glob..func1\n\t/mnt/d/development/emri/cmd/db.go:23\ngithub.com/spf13/cobra.(*Command).execute\n\t/home/ccorbett/go/1.21.5/pkg/mod/github.com/spf13/cobra@v1.8.1/command.go:989\ngithub.com/spf13/cobra.(*Command).ExecuteC\n\t/home/ccorbett/go/1.21.5/pkg/mod/github.com/spf13/cobra@v1.8.1/command.go:1117\ngithub.com/spf13/cobra.(*Command).Execute\n\t/home/ccorbett/go/1.21.5/pkg/mod/github.com/spf13/cobra@v1.8.1/command.go:1041\ngithub.com/asciifaceman/emri/cmd.Execute\n\t/mnt/d/development/emri/cmd/root.go:35\nmain.main\n\t/mnt/d/development/emri/main.go:10\nruntime.main\n\t/home/ccorbett/.goenv/versions/1.21.5/src/runtime/proc.go:267"}
*/

func (p *PostgresZap) Error(_ctx context.Context, format string, data ...interface{}) {
	//zap.S().Errorw("gorm message", "fmt", format, "data", fmt.Sprintf("%v", append([]interface{}{utils.FileWithLineNum()}, data...)...))
	var msgs []string
	errtype := "UNDEFINED"
	for _, piece := range data {
		msgs = append(msgs, fmt.Sprintf("%v", piece))
	}
	assembled := strings.Join(msgs, "  ")
	if strings.Contains(assembled, "failed to connect to") {
		errtype = "FAILED CONNECTION"
	}

	p.l.Errorw("gorm error",
		"source", utils.FileWithLineNum(),
		"type", errtype,
		"err", assembled,
	)

}

func (p *PostgresZap) Info(_ctx context.Context, format string, data ...interface{}) {
	p.l.Infow("gorm message",
		"status", "UNHANDLED",
		"fmt", format,
		"source", utils.FileWithLineNum(),
		"data", fmt.Sprintf("%v", data...),
	)
}

func (p *PostgresZap) Warn(_ctx context.Context, format string, data ...interface{}) {
	p.l.Warnw("gorm message",
		"status", "UNHANDLED",
		"source", utils.FileWithLineNum(),
		"data", fmt.Sprintf("%v", data...),
	)
}

func (p *PostgresZap) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	delta := time.Since(begin)
	switch {
	case err != nil && (!errors.Is(err, logger.ErrRecordNotFound) || !p.IgnorerecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			p.l.Errorw("gorm trace error",
				"err", err,
				"elapsed", float64(delta.Nanoseconds())/1e6,
				"sql", sql,
			)
		} else {
			p.l.Errorw("gorm trace error",
				"err", err,
				"elapsed", float64(delta.Nanoseconds())/1e6,
				"sql", sql,
				"rows", rows,
			)
		}
	case delta > p.SlowThreshold && p.SlowThreshold != 0:
		sql, rows := fc()
		if rows == -1 {
			p.l.Warnw("gorm slow sql",
				"threshold", p.SlowThreshold,
				"elapsed_ns", float64(delta.Nanoseconds())/1e6,
				"sql", sql,
			)
		} else {
			p.l.Warnw("gorm slow sql",
				"threshold", p.SlowThreshold,
				"elapsed_ns", float64(delta.Nanoseconds())/1e6,
				"sql", sql,
				"rows", rows,
			)
		}
	default:
		sql, rows := fc()
		if rows == -1 {
			p.l.Debugw("gorm query",
				"elapsed_ns", float64(delta.Nanoseconds())/1e6,
				"sql", sql,
			)
		} else {
			p.l.Debugw("gorm query",
				"elapsed_ns", float64(delta.Nanoseconds())/1e6,
				"sql", sql,
				"rows", rows,
			)
		}
	}
}
