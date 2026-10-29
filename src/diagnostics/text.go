package diagnostics

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func WriteText(w io.Writer, report Report) error {
	var text strings.Builder
	system, env := report.System, report.Environment
	fmt.Fprintf(&text, "SSUI diagnosis, %s\n\n", report.GeneratedAt.Format("2006-01-02 15:04:05 MST"))
	fmt.Fprintf(&text, "SSUI %s (%s) on %s/%s", system.Version, system.Branch, system.OS, system.Architecture)
	if system.OSRelease != "" {
		fmt.Fprintf(&text, ", %s", system.OSRelease)
	}
	if env.Kernel != "" {
		fmt.Fprintf(&text, ", kernel %s", env.Kernel)
	}
	fmt.Fprintf(&text, "\nContainer: %s", yesNo(system.Container))
	if len(env.ContainerSignals) > 0 {
		fmt.Fprintf(&text, " (%s)", strings.Join(env.ContainerSignals, ", "))
	}
	if env.UID >= 0 {
		fmt.Fprintf(&text, ", UID %d, GID %d", env.UID, env.GID)
	}
	fmt.Fprintf(&text, ", %d CPUs", env.CPUs)
	if env.MemoryMB > 0 {
		fmt.Fprintf(&text, ", %d MB RAM", env.MemoryMB)
	}
	fmt.Fprintf(&text, "\n\nOverall: %s\n", report.Overall)

	if len(report.Findings) > 0 {
		text.WriteString("\nFindings:\n")
		for _, finding := range report.Findings {
			fmt.Fprintf(&text, "  %-7s %s: %s\n", finding.Status, finding.Code, finding.Summary)
			if finding.Fix != "" {
				fmt.Fprintf(&text, "          Fix: %s\n", finding.Fix)
			}
		}
	}

	text.WriteString("\nFolders:\n")
	for _, volume := range env.Volumes {
		fmt.Fprintf(&text, "  %-8s %s, %s", volume.Name, volume.Status, volume.Path)
		if volume.Filesystem != "" {
			fmt.Fprintf(&text, ", %s on %s", volume.Filesystem, volume.MountPoint)
		}
		if volume.ReadOnly {
			text.WriteString(", read-only")
		}
		if volume.TotalMB > 0 {
			fmt.Fprintf(&text, ", %.1f GB free", float64(volume.FreeMB)/1024)
		}
		if volume.OwnerUID != nil {
			fmt.Fprintf(&text, ", owner %d:%d %s", *volume.OwnerUID, *volume.OwnerGID, volume.Mode)
		}
		text.WriteString("\n")
	}
	if report.Access != nil {
		fmt.Fprintf(&text, "\nFile access: %d paths checked, %d problems", report.Access.Checked, report.Access.Problems)
		if report.Access.Truncated {
			text.WriteString(", stopped early")
		}
		text.WriteString("\n")
		for _, issue := range report.Access.Issues {
			fmt.Fprintf(&text, "  %s\n", issue)
		}
	}

	if len(report.DNS) > 0 {
		text.WriteString("\nDNS:\n")
		for _, result := range report.DNS {
			if result.OK {
				fmt.Fprintf(&text, "  %-24s ok, %d ms (%s)\n", result.Host, result.LatencyMillis, result.Purpose)
			} else {
				fmt.Fprintf(&text, "  %-24s failed: %s (%s)\n", result.Host, result.Error, result.Purpose)
			}
		}
	}
	spacecat := report.Spacecat
	fmt.Fprintf(&text, "\nSpacecat: %s", spacecat.State)
	if spacecat.Port > 0 {
		fmt.Fprintf(&text, ", port %d", spacecat.Port)
	}
	for _, probe := range spacecat.Probes {
		if probe.Received {
			fmt.Fprintf(&text, ", %d bytes ok (%d ms)", probe.Size, probe.LatencyMillis)
		} else {
			fmt.Fprintf(&text, ", %d bytes lost", probe.Size)
		}
	}
	if spacecat.Message != "" {
		fmt.Fprintf(&text, "\n  %s", spacecat.Message)
	}
	if spacecat.Error != "" {
		fmt.Fprintf(&text, "\n  Error: %s", spacecat.Error)
	}
	text.WriteString("\n")

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintf(&text, "\nDocumentation: https://ssui.dev\n\nFull data:\n%s\n", data)
	_, err = io.WriteString(w, Redact(text.String()))
	return err
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
