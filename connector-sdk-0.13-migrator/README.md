# Connector SDK 0.13 migrator

`connector-sdk-0.13-migrator` is a tool that migrates connectors written with
Connector SDK `v0.12` and before to `v0.13` and using `conn-sdk-cli`.

## Example usage

```shell
go run main.go <path/to/connector>
```

The migration is done by _migrators_: one for each step in the migration, like:

- updating the SDK
- migrating the source connector
- etc.

If only one migrator needs to be run, then use this:

```shell
go run main.go <path/to/connector> UpdateSourceGo
```
