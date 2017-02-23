Feature: Predix CLI
  Scenario: Invoke the CLI
    Given I successfully run `run`
    Then the output should contain all of these lines:
    | NAME:                                                                                 |
    |   predix - A command line tool to interact with the Predix platform                   |
    |                                                                                       |
    | USAGE:                                                                                |
    |    main [global options] command [command options] [arguments...]                     |
    |                                                                                       |
    | VERSION:                                                                              |
    |    BUILT_FROM_SOURCE                                                                  |
    |                                                                                       |
    | COMMANDS:                                                                             |
    |      cf                  Run Cloud Foundry CLI commands on the Predix Platform        |
    |      login, l            Log user in to the Predix Platform                           |
    |      cache               Manage Predix CLI cache                                      |
    |      create-service, cs  Create a service instance                                    |
    |      service-info, si    List info for a service instance                             |
    |      uaa                 Manage Predix UAA instance                                   |
    |      help                Shows a list of commands or help for one command             |
    |                                                                                       |
    | ENVIRONMENT VARIABLES:                                                                |
    |   PREDIX_NO_CACHE=true      Do not use the cache to lookup apps and services          |
    |   PREDIX_NO_CF_BYPASS=true  Do not try to run an unknown command as a CF CLI command  |
    |                                                                                       |
    | GLOBAL OPTIONS:                                                                       |
    |    --help, -h     show help                                                           |
    |    --version, -v  print the version                                                   |
