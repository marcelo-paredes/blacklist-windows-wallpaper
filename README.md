# Blacklist Windows Wallpaper

This small utility will bring back your
current wallpaper if changed to a specific
blacklisted one (by filename).

```
IMPORTANT: This utility is meant for Windows only.
```


## Usage

To use a list of blacklisted wallpapers,
```powershell
blacklist-windows-wallpaper.exe --blacklist "path/to/wallpaper.jpg" --blacklist "path/to/other/wallpaper.jpg"
```

To use a list of blacklisted wallpapers from a file
```powershell
blacklis-windows-wallpaper.exe --blacklist "@path/to/blacklist.txt"
```

where the contents of `blacklist.txt` is something like
```text
wallpaper.jpg
other/wallpaper.jpg
```

## Possible Improvements so far

- [x] Add a command line argument to specify the blacklisted wallpaper name.
- [x] Add a command line argument to specify containing a list of blacklisted wallpapers.
- [ ] Find why the call to "setWallPaper" needs to be done twice and fix.
