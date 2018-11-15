// package fsrepo
//
// TODO explain the package roadmap...
//
//   .udfs/
//   ├── client/
//   |   ├── client.lock          <------ protects client/ + signals its own pid
//   │   ├── udfs-client.cpuprof
//   │   └── udfs-client.memprof
//   ├── config
//   ├── daemon/
//   │   ├── daemon.lock          <------ protects daemon/ + signals its own address
//   │   ├── udfs-daemon.cpuprof
//   │   └── udfs-daemon.memprof
//   ├── datastore/
//   ├── repo.lock                <------ protects datastore/ and config
//   └── version
package fsrepo

// TODO prevent multiple daemons from running
