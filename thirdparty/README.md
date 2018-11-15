thirdparty consists of Golang packages that contain no go-udfs dependencies and
may be vendored udfs/go-udfs at a later date.

packages in under this directory _must not_ import packages under
`udfs/go-udfs` that are not also under `thirdparty`.
