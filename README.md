# Mpkgr
Mpkgr is the build system for the mpkg package manager. It's currently non-functional as a whole, but some functionality is complete. Pull requests are welcome, if you can read my poorly documented code.

# Progress
### Build
- [X] Generate package.yml for new packages
- [X] Read package.yml to get package information
- [ ] Build package based on package.yml
- [ ] Parse install instructions
- [X] Describe package contents in header
- [X] ContentPayload for package files
- [ ] Payload for install instructions
- [X] Write everything to package.mpkg

### Install
- [ ] Read package.mpkg
- [ ] Determine structure of file
- [ ] Run install instructions for ContentPayload
