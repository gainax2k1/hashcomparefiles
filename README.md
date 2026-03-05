# hashcomparefiles
CLI tool that computes and compares file hashes, using SHA-256. 
This tool also makes it easy to delete duplicate files, move them to trash, or output a list of all duplicate files with their filesize. It uses SHA-256 to uniquely identify the file contents, so even if a duplicate file has a different name, it will still be flagged. The filesize is included for refrence, and for the remote chance of hash collision. 

* symlinks and empty files are ignored.

# Usage:

```python
hashcomparefiles -f (filename)
```
- returns hash value of a single file (filename)

```python
hashcomparefiles -d (directory)
```
- Scans through directory and sub-directories, displaying lists of duplicate files with their size and their hash value

```python
hashcomparefiles -trash (directory)
```
- (Linux only) Scans through directory, moving all duplicate files to trash after the first found instance. Currently, only works on primary drive in linux based systems. For non-linux systems, a folder is created in the working directory, and files are move into that, with corresponding .trashinfo files being created to record original file location.


```python
hashcomparefiles -delete (directory)
```
- Scans through directory, deleting all duplicate files after the first found instance


```python
hashcomparefiles -log (directory/logfilename) ...
```
- Creates a log file in the given directory/logfilename, default is current working directory

```python
hashcomparefiles --help
```
- Shows list of available flags and descriptions


