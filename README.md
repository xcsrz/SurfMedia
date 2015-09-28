# SurfMedia

### Features
* recursively scans the current directory for images and videos and provides a web interface to browse through them
* compiles as a single executable for portability
* uses (forked version of) prettyPhoto to slide through videos

### Prereqs
* https://github.com/xcsrz/prettyphoto needs to be in the statics directory

### Customizing and Building
* update files in the statics directory as desired
* run `make files` to encode the files into the autoGenStaticFiles.go file
* `go run` or `go build`

### Known Issues
* 100 large images on a page load overloads the browser (17 seconds to load for me, and scrolling is painful)
* swapping "/file" for "/image" will scale thumbnails on the fly but results in 110 second page loads
* styling needs love

## CONTRIBUTIONS WELCOME ##

### Status
* slow and clunky but otherwise functional
* thinking the next step is to try the imagemagick library

### Goals / Todos
* maintain the build to a single distributable file
* add command line flags to specify preferences (images per page, root directory, port number, etc)
* make the port auto incrementing port so multiple instances can run simultaneously
* syscall the opening of the url in the default browser (needs to be os specific - anyone have a windows machine?)
* speed spreed speed: it does a quick job of creating the list of images, but browsing is still crazy slow
