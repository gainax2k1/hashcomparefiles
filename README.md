# hash-file-compare
CLI tool that creates and compares file hashes, using SHA-256. 
This tool also makes it easy to delete duplicate files, move them to trash, or output a list of all duplicate files with their filesize. It uses SHA-256 to uniquely identify the file contents, so even if a duplicate file has a different name, it will still be flagged. The filesize is included for refrence, and for the remote chance of hash collision. 

(I developed this tool primarily to cleanup my data hoarding backup that had gotten out of hand with sloppy backups, included multiple copies of the same .iso image in multiple folders, repeated backups of cell phone pics, etc. I was able to quickly remove ~200 Gb of duplicate files on a 4 Tb drive. )

# Usage:

```python
hash-file-compare  (filename)
```
- returns hash value of (filename)


```python
hash-file-compare -d (directory)
```
- Scans through directory and sub-directories, displaying lists of duplicate files with their size and their hash value

```python
hash-file-compare -TRASH (directory)
```
- (in progress) Scans through directory, moving all duplicate files to trash after the first found instance


```python
hash-file-compare -DELETE (directory)
```
- Scans through directory, deleting all duplicate files after the first found instance
