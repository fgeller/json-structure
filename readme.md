# json-structure

Reads a JSON value from stdin and converts it to a JSON value that describes the structure of the original value.
For example:

```
~ % curl -s "https://api.openstreetmap.org/api/0.6/relation/10000000/full.json" | json-structure -dedupe -merge-objects | jq
{
  "attribution": "string",
  "copyright": "string",
  "elements": [
    {
      "changeset": "number",
      "id": "number",
      "lat": "number",
      "lon": "number",
      "members": [
        {
          "ref": "number",
          "role": "string",
          "type": "string"
        }
      ],
      "nodes": [
        "number"
      ],
      "tags": {
        "bicycle": "string",
        "colour:arrow": "string",
        "colour:back": "string",
        "colour:text": "string",
        "destination": "string",
        "destination:symbol": "string",
        "distance": "string",
        "foot": "string",
        "highway": "string",
        "horse": "string",
        "incline": "string",
        "information": "string",
        "lcn": "string",
        "motor_vehicle": "string",
        "ref": "string",
        "surface": "string",
        "tourism": "string",
        "tracktype": "string",
        "traffic_sign": "string",
        "type": "string"
      },
      "timestamp": "string",
      "type": "string",
      "uid": "number",
      "user": "string",
      "version": "number"
    }
  ],
  "generator": "string",
  "license": "string",
  "version": "string"
}
```

Original JSON value:

```json
{
  "version": "0.6",
  "generator": "CGImap 0.8.8 (1487389 spike-08.openstreetmap.org)",
  "copyright": "OpenStreetMap and contributors",
  "attribution": "http://www.openstreetmap.org/copyright",
  "license": "http://opendatacommons.org/licenses/odbl/1-0/",
  "elements": [
    {
      "type": "node",
      "id": 35352001,
      "lat": 50.2106234,
      "lon": 8.5856189,
      "timestamp": "2012-12-08T23:04:02Z",
      "version": 6,
      "changeset": 14206963,
      "user": "HoloDuke",
      "uid": 75317
    },
    {
      "type": "node",
      "id": 35352060,
      "lat": 50.2133772,
      "lon": 8.587342,
      "timestamp": "2022-01-23T16:42:23Z",
      "version": 3,
      "changeset": 116506987,
      "user": "PHerison",
      "uid": 28378
    },
    {
      "type": "node",
      "id": 35352061,
      "lat": 50.2129946,
      "lon": 8.5872135,
      "timestamp": "2022-01-23T16:42:23Z",
      "version": 3,
      "changeset": 116506987,
      "user": "PHerison",
      "uid": 28378
    },
    {
      "type": "node",
      "id": 35352062,
      "lat": 50.2126594,
      "lon": 8.5870583,
      "timestamp": "2022-01-23T16:42:23Z",
      "version": 3,
      "changeset": 116506987,
      "user": "PHerison",
      "uid": 28378
    },
    {
      "type": "node",
      "id": 35352063,
      "lat": 50.2121634,
      "lon": 8.5867429,
      "timestamp": "2022-01-23T16:42:23Z",
      "version": 5,
      "changeset": 116506987,
      "user": "PHerison",
      "uid": 28378
    },
    {
      "type": "node",
      "id": 2056819178,
      "lat": 50.2111796,
      "lon": 8.5860363,
      "timestamp": "2012-12-08T23:04:01Z",
      "version": 1,
      "changeset": 14206963,
      "user": "HoloDuke",
      "uid": 75317
    },
    {
      "type": "node",
      "id": 6772305769,
      "lat": 50.2106895,
      "lon": 8.5857061,
      "timestamp": "2019-09-04T15:05:00Z",
      "version": 1,
      "changeset": 74092935,
      "user": "Jonaes",
      "uid": 4433367,
      "tags": {
        "bicycle": "yes",
        "information": "guidepost",
        "ref": "OU.425.1",
        "tourism": "information"
      }
    },
    {
      "type": "way",
      "id": 5123563,
      "timestamp": "2021-10-06T17:09:44Z",
      "version": 8,
      "changeset": 112179717,
      "user": "toertsche",
      "uid": 23161,
      "nodes": [
        35352060,
        35352061,
        35352062,
        35352063,
        2056819178,
        35352001
      ],
      "tags": {
        "bicycle": "yes",
        "foot": "yes",
        "highway": "track",
        "horse": "yes",
        "incline": "down",
        "lcn": "yes",
        "motor_vehicle": "agricultural;forestry",
        "surface": "asphalt",
        "tracktype": "grade1",
        "traffic_sign": "DE:260,DE:1026-36"
      }
    },
    {
      "type": "relation",
      "id": 10000000,
      "timestamp": "2019-09-18T18:21:44Z",
      "version": 2,
      "changeset": 74641141,
      "user": "HoloDuke",
      "uid": 75317,
      "members": [
        {
          "type": "node",
          "ref": 6772305769,
          "role": "sign"
        },
        {
          "type": "node",
          "ref": 35352001,
          "role": "intersection"
        },
        {
          "type": "way",
          "ref": 5123563,
          "role": "to"
        }
      ],
      "tags": {
        "bicycle": "yes",
        "colour:arrow": "green",
        "colour:back": "white",
        "colour:text": "green",
        "destination": "Friedrichsdorf;Bad Homburg",
        "destination:symbol": ";train_station",
        "distance": "8.0;4.6",
        "type": "destination_sign"
      }
    }
  ]
}
```
