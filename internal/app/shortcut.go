package app

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

const (
	gsettingsPath        = "/usr/bin/gsettings"
	shortcutPath         = "/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/wgupdown-gui-toggle/"
	shortcutName         = "Toggle VPN"
	shortcutBinding      = "<Primary><Alt>v"
	shortcutDescription  = "Shortcut: Ctrl+Alt+V"
	customKeybindingsKey = "custom-keybindings"
)

var gsettingsArrayEntryPattern = regexp.MustCompile(`'([^']+)'`)

func ensureToggleShortcut(ctx context.Context, ifaceName string) error {
	if _, err := exec.LookPath("gsettings"); err != nil {
		return fmt.Errorf("gsettings not found: %w", err)
	}

	shortcutCommand := fmt.Sprintf("/usr/bin/sudo -n /usr/local/bin/wgupdown toggle %s", ifaceName)

	paths, err := readCustomShortcutPaths(ctx)
	if err != nil {
		return err
	}

	paths = appendIfMissing(paths, shortcutPath)

	if err := runGSettings(ctx,
		"set",
		"org.gnome.settings-daemon.plugins.media-keys",
		customKeybindingsKey,
		formatGSettingsArray(paths),
	); err != nil {
		return err
	}

	schemaWithPath := "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:" + shortcutPath

	for _, args := range [][]string{
		{"set", schemaWithPath, "name", shortcutName},
		{"set", schemaWithPath, "command", shortcutCommand},
		{"set", schemaWithPath, "binding", shortcutBinding},
	} {
		if err := runGSettings(ctx, args...); err != nil {
			return err
		}
	}

	return nil
}

func readCustomShortcutPaths(ctx context.Context) ([]string, error) {
	out, err := exec.CommandContext(
		ctx,
		gsettingsPath,
		"get",
		"org.gnome.settings-daemon.plugins.media-keys",
		customKeybindingsKey,
	).Output()
	if err != nil {
		return nil, fmt.Errorf("gsettings get %s failed: %w", customKeybindingsKey, err)
	}

	matches := gsettingsArrayEntryPattern.FindAllStringSubmatch(string(out), -1)
	paths := make([]string, 0, len(matches))
	for _, match := range matches {
		paths = append(paths, match[1])
	}

	return paths, nil
}

func formatGSettingsArray(values []string) string {
	if len(values) == 0 {
		return "[]"
	}

	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, fmt.Sprintf("'%s'", value))
	}

	return "[" + strings.Join(quoted, ", ") + "]"
}

func appendIfMissing(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}

	return append(values, value)
}

func runGSettings(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, gsettingsPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("gsettings %s failed: %w: %s", strings.Join(args, " "), err, stderr.String())
		}

		return fmt.Errorf("gsettings %s failed: %w", strings.Join(args, " "), err)
	}

	return nil
}
