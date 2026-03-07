<h1> hashcomparefiles</h1>
CLI tool that computes and compares file hashes, using SHA-256. 
This tool also makes it easy to delete duplicate files, move them to trash, or output a list of all duplicate files with their filesize. It uses SHA-256 to uniquely identify the file contents, so even if a duplicate file has a different name, it will still be flagged. The filesize is included for refrence, and for the remote chance of hash collision. 

* symlinks and empty files are ignored.

# Usage:

```python
hashcomparefiles (filename/directory)
```
- returns hash value of a single file or through directory and sub-directories, displaying lists of duplicate files with their size and their hash value




```python
hashcomparefiles -remove (directory)
```
- Processes all files, and goes through each group of hash-match files, allowing per-file deletion or trashing if desired. *Only fully supports FreeDesktop spec on primary drive in Linux based systems.* For non-Linux systems, or if running on an external mount in Linux, a folder is created in the working directory, and files are move into that, with corresponding .trashinfo files being created to record original file location. For this reason, it's highly recommended to run from the drive where files are stored. I would like to improve this in the future. 

```python
hashcomparefiles -p (directory)
```
- Scans through directory once without hashing to get total file count, then hashes the second run. Potentially useful for large runs (+1,000 files) to determine scale of run, at the cost of additional disk hits.


```python
hashcomparefiles -log (directory/logfilename) ...
```
- Creates a log file in the given directory/logfilename, default is current working directory.

```python
hashcomparefiles --help
```
- Shows list of available flags and descriptions

```python
cat (filename) | hashcomparefiles -(flag)
```
- Pipe in list of files and or folders to compare against each other. Flags maintain functionality.

# Examples:

```python
hashcomparefiles testdata/
2026/03/06 18:48:58 Duplicate files with hash: c0f5efbef0fe98aa90619444250b1a5eb23158d6686f0b190838f3d544ec85b9
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt size: 10
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt size: 10
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt size: 10
2026/03/06 18:48:58 Duplicate files with hash: 7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileB.txt size: 20
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt size: 20
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt size: 20
2026/03/06 18:48:58 Duplicate files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileC.txt size: 22
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/06 18:48:58 Duplicate files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileDDup.txt size: 30
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/06 18:48:58  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt size: 30
2026/03/06 18:48:58 Total files processed: 16
2026/03/06 18:48:59  (Done)


```


```python
hashcomparefiles -log default testdata/
2026/03/06 18:58:22 Duplicate files with hash: 7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileB.txt size: 20
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt size: 20
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt size: 20
2026/03/06 18:58:22 Duplicate files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileC.txt size: 22
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/06 18:58:22 Duplicate files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileDDup.txt size: 30
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt size: 30
2026/03/06 18:58:22 Duplicate files with hash: c0f5efbef0fe98aa90619444250b1a5eb23158d6686f0b190838f3d544ec85b9
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt size: 10
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt size: 10
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt size: 10
2026/03/06 18:58:22 Total files processed: 16
2026/03/06 18:58:23  (Done)
gainax2k1@pop-os:~/Documents/workspace/hashcomparefiles$ more log.log 
2026/03/06 18:58:22 Duplicate files with hash: 7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileB.txt size: 20
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt size: 20
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt size: 20
2026/03/06 18:58:22 Duplicate files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileC.txt size: 22
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/06 18:58:22 Duplicate files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileDDup.txt size: 30
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt size: 30
2026/03/06 18:58:22 Duplicate files with hash: c0f5efbef0fe98aa90619444250b1a5eb23158d6686f0b190838f3d544ec85b9
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt size: 10
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt size: 10
2026/03/06 18:58:22  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt size: 10
2026/03/06 18:58:22 Total files processed: 16
2026/03/06 18:58:23  (Done)
```



```python
hashcomparefiles -p /home/gainax2k1/
 / Files processed: 408900
2026/03/07 03:47:11 Total files to process: 408943
 \ Files processed: 73500

 <fast forward and edited for sheer volume of output>

2026/03/07 03:48:03 Files with hash: e263c7c51686a204417e8b670bb1525bae5dbf73407e678ba645b3f6f2e3e22f
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/common/Proton - Experimental/files/share/default_pfx/drive_c/windows/system32/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/common/Proton - Experimental/files/share/default_pfx/drive_c/windows/syswow64/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/common/Proton - Experimental/files/share/wine/nls/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/compatdata/1041920/pfx/drive_c/windows/system32/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/compatdata/1041920/pfx/drive_c/windows/syswow64/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/compatdata/1493710/pfx/drive_c/windows/system32/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/compatdata/1493710/pfx/drive_c/windows/syswow64/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/compatdata/1590760/pfx/drive_c/windows/system32/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/compatdata/1590760/pfx/drive_c/windows/syswow64/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/compatdata/3617780/pfx/drive_c/windows/system32/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/compatdata/3617780/pfx/drive_c/windows/syswow64/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/compatdata/3719980/pfx/drive_c/windows/system32/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/compatdata/3719980/pfx/drive_c/windows/syswow64/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/compatdata/4146670/pfx/drive_c/windows/system32/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/compatdata/4146670/pfx/drive_c/windows/syswow64/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/compatdata/4329470/pfx/drive_c/windows/system32/c_737.nls size: 66594
2026/03/07 03:48:03  - /home/gainax2k1/.local/share/Steam/steamapps/compatdata/4329470/pfx/drive_c/windows/syswow64/c_737.nls size: 66594

<edited again for sheer volume of output>


2026/03/07 03:51:03 Total unique hashes with numerous files in this batch: 83563
2026/03/07 03:51:04 (Done)

```

```python
cat testFilesList.txt | hashcomparefiles -p -remove
2026/03/07 04:30:19 Total files to process: 8
2026/03/07 04:30:19 Files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/07 04:30:19  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/07 04:30:19  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/07 04:30:19  -- Duplicates: 2
Delete file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > s
Delete file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > t
2026/03/07 04:30:25 Deleted duplicate file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt
2026/03/07 04:30:25 Files with hash: c0f5efbef0fe98aa90619444250b1a5eb23158d6686f0b190838f3d544ec85b9
2026/03/07 04:30:25  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt size: 10
2026/03/07 04:30:25  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt size: 10
2026/03/07 04:30:25  -- Duplicates: 2
Delete file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > c
2026/03/07 04:30:26 (Done)
```


# Installation Instructions: 
- Install Go, if not already installed, instuctions here:
https://go.dev/doc/install

- Download repo
- In root folder of repo, run:
```python
go build
go install
```

That should be everything! Any issues with the project, feel free to reach out to me. 
Thanks, and have a day. =^.^=
