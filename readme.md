![Go](https://img.shields.io/badge/Go-1.22.1-blue)
![Platform](https://img.shields.io/badge/Platform-Windows-lightgrey)
![Archeology](https://img.shields.io/badge/Archeology-In%20Progress-brown?logo=fossil&logoColor=white)
![Danger](https://img.shields.io/badge/Dragons-Here%20Be-red?logo=firebase&logoColor=white)
![Support](https://img.shields.io/badge/Support-LOL%20No-critical?logo=stackexchange&logoColor=white)
![Version](https://img.shields.io/badge/Version-Jurassic-yellow?logo=github&logoColor=white)

## 👵 Welcome, digital archaeologist! 

 What Are These Ancient Relics?
 
 You've stumbled upon the fossilized remains of earlier Stationeers Server UI versions. These branches are preserved primarily for historical documentation, developer nightmares, and occasionally making @JacksonTheMaster cry.

## ⚠️ IMPORTANT: WHY YOU SHOULDN'T BE HERE ⚠️

### These branches are about as functional as a teapot in summer.

**SERIOUSLY, GO DOWNLOAD THE [LATEST RELEASE](https://github.com/JacksonTheMaster/StationeersServerUI/releases/latest) INSTEAD.**

Old branches will attempt to download files that have long since moved on to the great repository in the sky. They're maintained purely for:

# Stationeers Dedicated Server Control v1.1
Stationeers Dedicated Server Control is a web-based tool for managing a Stationeers dedicated server. It offers an intuitive UI and a robust REST API for server operations, configuration management, and backup handling.
I created this project to make it easier for me to manage my Stationeers server more efficiently, especially to restore backups.
I found that the Stationeers server was not very user-friendly, and I wanted to create a tool that would make it easier to manage the server.
Also I wanted my friends to be able to start, stop and manage the Server without having to ask me to restore the lastest backup because some base exploded. So here we are.

DISCLAIMER: PUBLISHING THIS UI TO THE WEB SHOULD ONLY BE DONE BEHIND A SECURE AUTHENTICATION MECHANISM; THIS SHOULD NOT UNDER ANY CIRCUMSTANCES BE PORT FORWARDED STRAIGHT OUT!

## Linux Support

For all my Linux people out there: yes, it’s coming. But, as of right now, I’m relying on Powershell to execute the game server, so... yeah. That’s pretty much it. Stay tuned for when the Linux fairies finally decide to bless this project.
I'd say you can expect a Docker release then aswell. Because who wouldn't make a Docker release then.

## Features
| UI Overview | Configuration | Backup Management |
|:-----------:|:-------------:|:-----------------:|
| ![UI Overview](media/UI-1.png) | ![Configuration](media/UI-2.png) | ![Backup Management](media/UI-3.png) |

- Start and stop the server (because sometimes, the server just needs a break too).
- View real-time server output (so you can stare at the logs and pretend you know what's going on).
- Manage server configurations
- List and restore backups
- Fully functional REST API for all operations (because who doesn’t love APIs, right..?).

### Coming *Soon™*:
- Discord integration with a shiny control channel, log channel, save channel, and player log channel (all the channels you never knew you needed).
- Granular backups with date-sorted directories, because losing unsaved progress is *so 2023*.
- Linux support (yes, it’s real, and yes, it’s happening—just not right now, *soon™*)

## Requirements
- Windows OS
- Downloaded and installed the Stationeers Dedicated Server.
- Administrative Privileges (Hostesly i havnt tested it without running as admin, but  it..should'nt run without..?!

## Quick Installation Instrcutions for Administrators & Server Operators

1. Download & Extract release ZIP from GitHub.
2. Move "startStatoneersServerUI.exe" and the "UIMod" folder to the server's executable directory.
3. Run "startStatoneersServerUI.exe". (If you start "UIMod/Stationeers-ServerUI.exe", the Server wont auto restart)
4. Access UI at `http://<server-ip>:8080`.
5. Open firewall ports 27015, 27016, 8080.
6. Check /config before starting the server.


## Detailed Installation Instrcutions for "Normal" Windows Users

1. Go to the link: https://github.com/jacksonthemaster/StationeersServerUI/releases.
2. Find the latest release and click to download the ZIP file.
3. Once downloaded, locate the ZIP file, right-click on it, and select "Extract All...".
4. Choose a folder where you want to save the extracted files and click "Extract".
5. Open the folder with the extracted files and locate "startStatoneersServerUI.exe".
6. Cut (Ctrl+X) or copy (Ctrl+C) "startStatoneersServerUI.exe".
7. Navigate to the folder where you have installed your Stationeers Dedicated Server.
8. Paste (Ctrl+V) "startStatoneersServerUI.exe" into this folder.
9. Go back to the extracted files folder and find the "UIMod" folder.
10. Cut (Ctrl+X) or copy (Ctrl+C) the "UIMod" folder.
11. Paste (Ctrl+V) the "UIMod" folder into the same folder where your Stationeers Dedicated Server executable is located.
13. Double-click "startStatoneersServerUI.exe" to run it. Do not run "UIMod/Stationeers-ServerUI.exe" unless you DONT want the server to auto restart.
14. Open your web browser and type `http://<IP-OF-YOUR-SERVER>:8080` in the address bar. Replace `<IP-OF-YOUR-SERVER>` with the actual IP address of your server. You can find this by opening the Command Prompt and typing `ipconfig`.
15. To allow other users to connect to your UI and the Server, open the Windows Firewall settings:
    - Go to Control Panel > System and Security > Windows Defender Firewall.
    - Click "Advanced settings" on the left.
    - In the Windows Firewall with Advanced Security window, click "Inbound Rules" on the left.
    - Click "New Rule..." on the right.
    - Select "Port" and click "Next".
    - Choose "TCP" and enter "27015, 27016, 8080" in the Specific local ports field. Click "Next".
    - Allow the connection and click "Next".
    - Select the network types to apply this rule (usually Domain, Private, and Public) and click "Next".
    - Name the rule something recognizable (e.g., "Stationeers Server Ports") and click "Finish".
    - __Note__:  Depending on your Setup, you might need to Port forward those ports on your router. For this, please consider using google or any other search engine exept bing to find a tutorial on how to do this.
16. Before starting your server, ensure the configuration files on the /config page are set up correctly.


## REST API Information

This server is based on Go, so it's basically a REST-API with some HTML files on top. All UI actions are API calls, so you can fully use the API to control the server.

### API Endpoints

- **Start Server**: `/start` (GET)
- **Stop Server**: `/stop` (GET)
- **Get Server Output**: `/output` (GET)
- **List Backups**: `/backups` (GET)
- **Restore Backup**: `/restore?index=<index>` (GET)
- **Edit Configuration**: `/config` (GET)
- **Save Configuration**: `/saveconfig` (POST Form Data)

### Form Data Explanation

- **SaveFileName**: The name of the save file to load. This is the name of the file without the extension. Example: `Mars`.
- **Settings**: The server settings. Use the UI to get the correct settings if you're unsure.

## UI

The web interface provides buttons to start and stop the server, edit configuration, and manage backups. The current server status and console output are displayed in real-time.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests!

## Acknowledgments

- [JacksonTheMaster](https://github.com/JacksonTheMaster) Developed with ❤️ and 💧 by J. Langisch.
- [Go](https://go.dev/) for the Go programming language.
- [RocketWerkz](https://github.com/RocketWerkz) for creating the Stationeers game.
# 🦖 Prehistoric Branches of Stationeers Server UI

# About historic releases
1. Historical reference
2. Code archaeology
3. Developer nostalgia
4. The following (sigh):

## 🤦‍♂️ "But I REALLY Need This Specific Old Version!"

Are you absolutely, positively sure? Fine. Here's what you'll need to do:

1. Clone this specific branch
2. Edit `src/config/config.Branch` to match your current branch name
3. Build using:
   - For ancient versions: `go run build.go`
   - For slightly less ancient versions: `go run build/build.go`
4. Accept that you are now on your own personal journey of pain and probably fixing already fixed bugs

## 🧟‍♂️ "I'm Still Going To Try This!"

Congratulations on your determination! When things go sideways (not if, WHEN), please follow these steps:

1. Acknowledge your choices
2. Download the latest release instead
3. Pretend this never happened

## 📜 Historical Context

These branches are NOT maintained and exist only for continuity of development history and historical reference. They represent the evolutionary path that led to the current codebase, much like how chickens evolved from dinosaurs. Do I like dinosaurs? Maybe.

## 🤡 Support Policy

Support for these branches is provided on a "point and laugh" basis only. 

My official, heartfelt recommendation is:

```
$ rm -rf old-branch/
$ wget https://github.com/JacksonTheMaster/StationeersServerUI/releases/latest/download/StationeersServerUI
```

## 📊 Success Rate Prediction

Your chance of successfully running an old branch without modification:

```
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░░░ 39.42%
```

## License

This project, including it's ancient branches, is licensed under the STATIONEERS SERVER UI LICENSE AGREEMENT - see the ([LICENSE](https://github.com/JacksonTheMaster/StationeersServerUI/blob/main/LICENSE)) file for details.
