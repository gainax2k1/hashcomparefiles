# hashcomparefiles
CLI tool that computes and compares file hashes, using SHA-256. 
This tool also makes it easy to delete duplicate files, move them to trash, or output a list of all duplicate files with their filesize. It uses SHA-256 to uniquely identify the file contents, so even if a duplicate file has a different name, it will still be flagged. The filesize is included for refrence, and for the remote chance of hash collision. 

* symlinks and empty files are ignored.

# Usage:

```python
hashcomparefiles (filename/directory)
```
- returns hash value of a single file or through directory and sub-directories, displaying lists of duplicate files with their size and their hash value

```python
hashcomparefiles -trash (directory)
```
- (Linux only) Scans through directory, moving all duplicate files to trash after the first found instance. Currently, only fully works on primary drive in Linux based systems. For non-Linux systems, a folder is created in the working directory, and files are move into that, with corresponding .trashinfo files being created to record original file location.
- (Currently, -trash uses os.Rename to trash the files, which might not work correctly on external mounts/devices. In these cases, it switches to a copy/delete process instead, which doesn not follow the FreeDesktop spec. I would like to improve this in the future.)
- In the future, will add ability to choose which file(s) to trash.


```python
hashcomparefiles -delete (directory)
```
- Scans through directory, deleting all duplicate files after the first found instance
- In the future, will add ability to choose which file(s) to delete.

```python
hashcomparefiles -p (directory)
```
- Scans through directory once without hashing to get total file count, then hashes the second run. Potentially useful for large runs (+1,000 files) to determine progress, at the cost of additional disk hits

```python
hashcomparefiles -v (directory)
```
- Verbose - outputs every hash value and filesize, even if not duplicated.

```python
hashcomparefiles -log (directory/logfilename) ...
```
- Creates a log file in the given directory/logfilename, default is current working directory

```python
hashcomparefiles --help
```
- Shows list of available flags and descriptions



