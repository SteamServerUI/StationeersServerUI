package discordbot

import (
	"StationeersServerUI/src/backupmgr"
	"StationeersServerUI/src/gamemgr"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func handleListCommand(content string) {
	fmt.Println("!list command received, fetching backup list...")

	// Extract the "top" number or "all" option from the command
	parts := strings.Split(content, ":")
	top := 5 // Default to 5
	var err error
	if len(parts) == 2 {
		if parts[1] == "all" {
			top = -1 // No limit
		} else {
			top, err = strconv.Atoi(parts[1])
			if err != nil || top < 1 {
				SendMessageToControlChannel("❌Invalid number provided. Use `!list:<number>` or `!list:all`.")
				return
			}
		}
	}

	// Use the Global Backup Manager to get the list of backups
	limit := top
	if top == -1 {
		limit = 0 // 0 means all backups in the new system
	}

	backups, err := backupmgr.GlobalBackupManager.ListBackups(limit)
	if err != nil {
		fmt.Println("Failed to fetch backup list:", err)
		SendMessageToControlChannel("❌Failed to fetch backup list.")
		return
	}

	// If no backups are found
	if len(backups) == 0 {
		SendMessageToControlChannel("No backups found.")
		return
	}

	// Format each backup group and send to channel
	count := 0
	for _, backup := range backups {
		// Format the backup information using ModTime from BackupGroup
		backupInfo := fmt.Sprintf("**Backup %d** - %s", backup.Index, backup.ModTime.Format("2006-01-02 15:04:05"))

		// Send message to control channel
		SendMessageToControlChannel(backupInfo)
		fmt.Println("Successfully sent backup info to Discord:", backupInfo)

		count++

		time.Sleep(500 * time.Millisecond)

		if top > 0 && count >= top {
			break
		}
	}
}

func handleRestoreCommand(content string) {
	parts := strings.Split(content, ":")
	if len(parts) != 2 {
		SendMessageToControlChannel("❌Invalid restore command. Use `!restore:<index>`.")
		SendMessageToStatusChannel("⚠️Restore command received, but not able to restore Server.")
		return
	}

	indexStr := parts[1]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		SendMessageToControlChannel("❌Invalid index provided for restore.")
		SendMessageToStatusChannel("⚠️Restore command received, but not able to restore Server.")
		return
	}

	// Stop the server before restoring
	gamemgr.InternalStopServer()

	// Use the Global Backup Manager to restore the backup
	err = backupmgr.GlobalBackupManager.RestoreBackup(index)
	if err != nil {
		SendMessageToControlChannel(fmt.Sprintf("❌Failed to restore backup at index %d: %v", index, err))
		SendMessageToStatusChannel("⚠️Restore command received, but not able to restore Server.")
		return
	}

	SendMessageToControlChannel(fmt.Sprintf("✅Backup %d restored successfully, Starting Server...", index))

	// Sleep 5 sec to give the server time to prepare
	time.Sleep(5 * time.Second)

	// Start the server after restoration
	gamemgr.InternalStartServer()
}
