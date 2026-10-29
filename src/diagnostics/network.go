package diagnostics

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/config"
)

const dnsTimeout = 3 * time.Second

type DNSResult struct {
	Host          string `json:"host"`
	Purpose       string `json:"purpose"`
	OK            bool   `json:"ok"`
	Addresses     int    `json:"addresses,omitempty"`
	LatencyMillis int64  `json:"latencyMs"`
	Error         string `json:"error,omitempty"`
}

func checkDNS() []DNSResult {
	hosts := []DNSResult{
		{Host: "api.github.com", Purpose: "SSUI updates"},
		{Host: "steamcdn-a.akamaihd.net", Purpose: "SteamCMD download"},
		{Host: "jxsn.dev", Purpose: "Spacecat and support connection"},
	}
	if config.GetIsDiscordEnabled() {
		hosts = append(hosts, DNSResult{Host: "discord.com", Purpose: "Discord integration"})
	}
	var wg sync.WaitGroup
	for i := range hosts {
		wg.Add(1)
		go func(result *DNSResult) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), dnsTimeout)
			defer cancel()
			started := time.Now()
			addresses, err := net.DefaultResolver.LookupHost(ctx, result.Host)
			result.LatencyMillis = time.Since(started).Milliseconds()
			result.OK, result.Addresses = err == nil && len(addresses) > 0, len(addresses)
			if err != nil {
				result.Error = err.Error()
			}
		}(&hosts[i])
	}
	wg.Wait()
	return hosts
}

func (report *Report) checkDNS() {
	failed := []string{}
	for _, result := range report.DNS {
		if !result.OK {
			failed = append(failed, fmt.Sprintf("%s (%s)", result.Host, result.Purpose))
		}
	}
	switch {
	case len(failed) == 0:
	case len(failed) == len(report.DNS):
		report.add(codeDNS, StatusFail, "SSUI cannot resolve any internet host names.", "Check the DNS server of this machine or container and its internet access.")
	default:
		report.add(codeDNS, StatusWarn, "SSUI cannot resolve "+strings.Join(failed, ", ")+".", "Check the DNS server and firewall rules for these hosts.")
	}
}
