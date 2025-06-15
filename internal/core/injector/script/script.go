package script

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/flew1x/grpc-chaos-proxy/internal/apperr"
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"github.com/flew1x/grpc-chaos-proxy/internal/entity"
)

const tmpScriptPrefix = "/tmp/grpc-chaos-script-"

type Injector struct {
	Language   string
	Source     string
	Entrypoint string
	TimeoutMS  int
	Env        map[string]string
	Args       []string
}

func NewScriptInjector(cfg any) (engine.Injector, error) {
	sa, ok := cfg.(*config.ScriptAction)
	if !ok || sa == nil {
		return nil, apperr.ErrInvalidConfig
	}

	return &Injector{
		Language:   sa.Language,
		Source:     sa.Source,
		Entrypoint: sa.Entrypoint,
		TimeoutMS:  sa.TimeoutMS,
		Env:        sa.Env,
		Args:       sa.Args,
	}, nil
}

// Apply executes the script with the provided configuration
func (s *Injector) Apply(f *engine.Frame) error {
	if s.Language != entity.ScriptLanguageBash.String() && s.Language != entity.ScriptLanguageSh.String() {
		return fmt.Errorf("only 'sh' or 'bash' scripts supported")
	}

	ctx := f.Ctx

	if s.TimeoutMS > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(s.TimeoutMS)*time.Millisecond)
		defer cancel()
	}

	tmpFile, err := writeTempScript(s.Source)
	if err != nil {
		return fmt.Errorf("failed to write temp script: %w", err)
	}

	defer os.Remove(tmpFile)

	cmd := exec.CommandContext(ctx, s.Language, tmpFile)
	if len(s.Args) > 0 {
		cmd.Args = append(cmd.Args, s.Args...)
	}

	if len(s.Env) > 0 {
		env := os.Environ()
		for k, v := range s.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}

		cmd.Env = env
	}

	output, err := cmd.CombinedOutput()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("script timeout: %w", ctx.Err())
	}

	if err != nil {
		return fmt.Errorf("script error: %w, output: %s", err, string(output))
	}

	parseScriptOutput(string(output), f)

	return nil
}

// writeTempScript creates a temporary script file with the provided source code
func writeTempScript(source string) (string, error) {
	tmpFile := fmt.Sprintf("%s%d.sh", tmpScriptPrefix, time.Now().UnixNano())
	file, err := os.Create(tmpFile)
	if err != nil {
		return "", err
	}

	defer file.Close()

	if _, err := file.WriteString(source); err != nil {
		return "", err
	}

	if err := file.Chmod(0700); err != nil {
		return "", err
	}

	return tmpFile, nil
}

// parseScriptOutput processes the output of the script execution
func parseScriptOutput(output string, f *engine.Frame) {
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "X-CHAOS-ERROR:") {
			msg := strings.TrimSpace(strings.TrimPrefix(line, "X-CHAOS-ERROR:"))
			f.MD.Set("x-chaos-error", msg)
		}

		if strings.HasPrefix(line, "X-CHAOS-HEADER:") && f.MD != nil {
			kv := strings.TrimSpace(strings.TrimPrefix(line, "X-CHAOS-HEADER:"))

			parts := strings.SplitN(kv, "=", 2)
			if len(parts) == 2 {
				f.MD.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}
}

func init() {
	engine.Register(entity.ScriptType, NewScriptInjector)
}
