package diagnostics

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/config"
	"github.com/SteamServerUI/StationeersServerUI/v6/src/logger"
)

var ErrNoRecording = errors.New("no stopped support recording is available")

var recorded struct {
	sync.Mutex
	dir    string
	report Report
}

func StopRecording() (logger.SupportRecording, error) {
	recording, err := logger.StopSupportRecording()
	if errors.Is(err, logger.ErrNoSupportRecording) {
		return recording, err
	}
	recorded.Lock()
	recorded.dir, recorded.report = recording.Dir, Run()
	recorded.Unlock()
	return recording, err
}

func WritePackage(w io.Writer) error {
	archive := zip.NewWriter(w)
	if err := writeReport(archive, Run()); err != nil {
		return err
	}
	return archive.Close()
}

func WriteRecordingPackage(w io.Writer) error {
	recording := logger.SupportRecordingStatus()
	recorded.Lock()
	report, matches := recorded.report, recorded.dir != "" && recorded.dir == recording.Dir
	recorded.Unlock()
	if recording.Active || !matches {
		return ErrNoRecording
	}

	archive := zip.NewWriter(w)
	if err := writeReport(archive, report); err != nil {
		return err
	}
	redact := secretRedactor()
	for _, name := range []string{"ssui.log", "gameserver.log"} {
		if err := writeRedactedLog(archive, filepath.Join(recording.Dir, name), redact); err != nil {
			return err
		}
	}
	return archive.Close()
}

func writeReport(archive *zip.Writer, report Report) error {
	file, err := archive.Create("diagnosis.json")
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return err
	}
	file, err = archive.Create("report.html")
	if err != nil {
		return err
	}
	return reportTemplate.Execute(file, report)
}

func writeRedactedLog(archive *zip.Writer, path string, redact *strings.Replacer) error {
	source, err := os.Open(path)
	if err != nil {
		return err
	}
	defer source.Close()
	if info, err := source.Stat(); err != nil || info.Size() == 0 {
		return err
	}
	file, err := archive.Create(filepath.Base(path))
	if err != nil {
		return err
	}
	reader := bufio.NewReader(source)
	for {
		line, err := reader.ReadString('\n')
		if _, writeErr := redact.WriteString(file, line); writeErr != nil {
			return writeErr
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func secretRedactor() *strings.Replacer {
	var pairs []string
	for _, secret := range []string{config.GetServerPassword(), config.GetServerAuthSecret(), config.GetAdminPassword(), config.GetDiscordToken(), config.GetJwtKey()} {
		if len(secret) >= 4 {
			pairs = append(pairs, secret, "[REDACTED]")
		}
	}
	return strings.NewReplacer(pairs...)
}

var reportTemplate = template.Must(template.New("report").Funcs(template.FuncMap{
	"json": func(value any) string {
		data, _ := json.MarshalIndent(value, "", "  ")
		return string(data)
	},
}).Parse(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>SSUI support report</title>
<style>
body{max-width:900px;margin:0 auto;padding:24px;font:15px/1.5 system-ui,sans-serif;background:#11151b;color:#e9edf3}
.muted{color:#9ba7b5}.card{margin:10px 0;padding:12px 16px;border:1px solid #303946;border-left-width:4px;border-radius:8px;background:#1a2029}
.fail{border-left-color:#e05d5d}.warn{border-left-color:#e0a73a}.unknown{border-left-color:#8492a6}.healthy{border-left-color:#3fae6a}
code,pre{font:13px/1.45 ui-monospace,Consolas,monospace}pre{white-space:pre-wrap;color:#cbd5e1}summary{cursor:pointer}
</style></head><body>
<h1>SSUI support report</h1>
<p class="muted">{{.System.Version}} ({{.System.Branch}}) · {{.System.OS}}/{{.System.Architecture}} · generated {{.GeneratedAt.Format "2006-01-02 15:04:05 MST"}}</p>
<h2>Overall: {{.Overall}}</h2>
{{range .Findings}}<div class="card {{.Status}}"><code>{{.Code}}</code> {{.Summary}}{{if .Fix}}<br><span class="muted">{{.Fix}}</span>{{end}}</div>
{{else}}<div class="card healthy">Nothing needs attention.</div>{{end}}
<details><summary>Collected data</summary><pre>{{json .}}</pre></details>
</body></html>
`))
