# Mpkgr
Mpkgr is the build system for the mpkg package manager. It's currently non-functional as a whole, but some functionality is complete. Pull requests are welcome, if you can read my poorly documented code.

# Progress
### Build
- [ ] Build packages in isolated chroot environment (not possible until I can bootstrap the distro)
- [X] Generate package.yml for new packages
- [X] Read package.yml to get package information
- [X] Build package based on package.yml
- [X] Parse install instructions
- [X] Describe package contents in header
- [X] ContentPayload for package files
- [X] Write everything to package.mpkg

### Install
- [X] Read package.mpkg
- [X] Determine structure of file
- [X] Run install instructions for ContentPayload
