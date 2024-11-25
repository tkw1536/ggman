package dirs

import (
	"io/fs"
	"os"
)

// NewModBits holds the default permission bits for a new directory.
const NewModBits fs.FileMode = os.ModePerm
