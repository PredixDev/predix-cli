# Predix CLI
The CLI is a command line utility meant to simplify interaction with the Predix Cloud. It is a wrapper over the CF CLI and provides commands combining multiple steps of interaction into a single command.
This is a beta release of the CLI.

## Wiki
https://github.com/PredixDev/predix-cli/wiki

## Installation
Download the latest release from https://github.com/PredixDev/predix-cli/releases

### Linux / Mac OS X
- Extract the file 'predix-cli.tar.gz'
- Navigate to the extracted folder and run 'sudo ./install'

### Windows
- Extract the file 'predix-cli.tar.gz'
- Copy predix.exe from bin/win64 in the extracted folder to somewhere on the PATH. Run `echo %PATH%` to see its value
- Make a symbolic short link spelled with px.exe pointing at predix.exe. e.g. mklink path-to-cli\px.exe path-to-cli\predix.exe
- Autocompletion is not supported on Windows
