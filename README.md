# predix-cli

A command line tool to interact with the Predix platform

## Features

- Login to the Predix Basic or Predix Select environment
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
@powershell -Command "(new-object net.webclient).DownloadFile('https://raw.githubusercontent.com/PredixDev/local-setup/master/setup-windows.bat','%TEMP%/setup-windows.bat')" && %TEMP%/setup-windows.bat /cf /predixcli
```

## Manual Installation instructions

The latest release is downloadable at https://github.com/PredixDev/predix-cli/releases

### Linux / Mac OS X
- Extract the file 'predix-cli.tar.gz'
- Navigate to the extracted folder and run 'sudo ./install'

### Windows
- Extract the file 'predix-cli.tar.gz'
- Copy predix.exe in bin/win64 in the extracted folder to somewhere on the PATH
- Make a symbolic short link spelled with px.exe pointing at predix.exe. e.g. mklink path-to-cli\px.exe path-to-cli\predix.exe- Autocompletion is not supported on Windows

[![Analytics](https://ga-beacon.appspot.com/UA-82773213-1/predixcli/readme?pixel)](https://github.com/PredixDev)
