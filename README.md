# predix-cli

A command line tool to interact with the Predix platform

## Features

- Login to the various Predix PoP environments (US West, US East, Frankfurt, Japan, etc)
- Define your own PoPs for internal or custom clouds.  See [below](#define-custom-cloud-login-pop-endpoints)
- Bash autocompletion for the Cloud Foundry CLI commands, parameters and arguments

## Installation

Use our one-click [local-setup](https://github.com/PredixDev/local-setup) installers

On Mac OS X

Run the command below in a terminal window to install Cloud Foundry CLI and the Predix CLI
```
bash <( curl https://raw.githubusercontent.com/PredixDev/local-setup/master/setup-mac.sh ) --cf --predixcli
```

On Windows

Open a Command Window as Administrator (Right click 'Run as Administrator') and run the command below
```
@powershell -Command "(new-object net.webclient).DownloadFile('https://raw.githubusercontent.com/PredixDev/local-setup/master/setup-windows.bat','%TEMP%\setup-windows.bat')" && %TEMP%/setup-windows.bat /cf /predixcli
```

## Manual Installation instructions

The latest release is downloadable at https://github.com/PredixDev/predix-cli/releases

### Linux / Mac OS X
- Extract the file 'predix-cli.tar.gz'
- Navigate to the extracted folder and run './install'

### Windows
- Extract the file 'predix-cli.tar.gz'
- Copy predix.exe in bin/win64 in the extracted folder to somewhere on the PATH
- Make a symbolic short link spelled with px.exe pointing at predix.exe. e.g. mklink path-to-cli\px.exe path-to-cli\predix.exe- Autocompletion is not supported on Windows

[![Analytics](https://ga-beacon.appspot.com/UA-82773213-1/predixcli/readme?pixel)](https://github.com/PredixDev)


### Define Custom Cloud Login PoP Endpoints

In the hidden directory ~/.predix create a pops.json file.
```
[
  {"name": "CF3", "url": "https://api.system.your-endpoint-here.ice.predix.io", "flag": "cf3", "usage": "Login to the Predix CF3 PoP"}
]
```
