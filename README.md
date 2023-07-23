# Go Find Files

This program will recursively walk the directory it is launched from and copy the filetypes specified as command line arguments to generated directory named FOUND.

## Usage

Copy all JPEG images to the FOUND directory: `./find-files .jpeg .jpg`

If you are short on space, extract (move) the files to the FOUND directory instead: `./find-files -x .jpeg .jpg`
