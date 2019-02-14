# Module Update Service

Small server app to update modules

## Usage

Create config file with list of modules on the server. Example `module.json`.

Build server or run

```bash
go run cmd/server/server.go --path PATH
```

## API

### Check module versions

`[GET] /sync`

Body:

```json
  {
    "deviceId": "uniqueDeviceID",
    "installedModules": [
      {
        "id": "ads",
        "version": 8
      },
      {
        "id": "proxy",
        "version": 3
      }
    ]
  }
```

Will return list of outdated modules ID

### Get module update

`[GET] /module/{id}`

Will return meta information with link to download a latest version of a module
