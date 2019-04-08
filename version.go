package ipfs

// CurrentCommit is the current git commit, this is set as a ldflag in the Makefile
var CurrentCommit string = "feat/ut-v0.1.2"

// CurrentVersionNumber is the current application's version literal
const CurrentVersionNumber = "0.4.18"

const ApiVersion = "/go-ipfs/" + CurrentVersionNumber + "/"
