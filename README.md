<h1> hashcomparefiles</h1>
A robust CLI tool that computes file hashes to identify duplicate files regardless of filename, using SHA-256, and presents them to the user. This tool also makes it easy to selectively delete duplicate files, move them to trash, or output a list of all duplicate files with their filesize.
It runs in a multi-pass method, minimizing disk hits and improving efficiency.

* symlinks and empty files are ignored
* sub-folders are automatically walked and included
* compatible with piping in lists of folders/filenames for more customization
* The filesize is included for reference, and for the remote chance of hash collision.

This tool was developed and tested in a Linux environment (Pop!_OS 24.04 LTS), but I would like it to eventually fully operate in Mac and Windows environments as well. In paricular, the trashing processes in Mac and Windows is incompatible, and the input when selectively removing files currently uses TTY (which likely won't work in Windows, but is currently untested).

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
hashcomparefiles -log (directory/logfilename) ...
```
- Creates a log file in the given directory/logfilename, default is current working directory.

```python
hashcomparefiles -min (integer) -max (integer) ...
```
- Set minimum and/or maximum filesizes (in bytes) to process.

```python
hashcomparefiles --help
```
- Shows list of available flags and descriptions

```python
cat (filename) | hashcomparefiles -(flag)
```
- Pipe in list of files and or folders to compare against each other. Flags maintain functionality.



# Installation Instructions: 

- Download compiled binaries under "Releases" here: 
https://github.com/gainax2k1/hashcomparefiles/releases

  or, if you want to build your own:

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

# Examples:
<h2> Small run, no flags</h2>

```python
$ hashcomparefiles testdata/
2026/03/08 16:18:33 Filecount after first pass: 17
2026/03/08 16:18:33 Filecount after second pass: 12
2026/03/08 16:18:33 Filecount after third pass: 12
2026/03/08 16:18:33 Groups of duplicates after shrink: 4
2026/03/08 16:18:33 Files with hash: 7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f
2026/03/08 16:18:33  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileB.txt size: 20
2026/03/08 16:18:33  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt size: 20
2026/03/08 16:18:33  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt size: 20
2026/03/08 16:18:33  -- Duplicates: 3
2026/03/08 16:18:33 Files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/08 16:18:33  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileDDup.txt size: 30
2026/03/08 16:18:33  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/08 16:18:33  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt size: 30
2026/03/08 16:18:33  -- Duplicates: 3
2026/03/08 16:18:33 Files with hash: c0f5efbef0fe98aa90619444250b1a5eb23158d6686f0b190838f3d544ec85b9
2026/03/08 16:18:33  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt size: 10
2026/03/08 16:18:33  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt size: 10
2026/03/08 16:18:33  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt size: 10
2026/03/08 16:18:33  -- Duplicates: 3
2026/03/08 16:18:33 Files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/08 16:18:33  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileC.txt size: 22
2026/03/08 16:18:33  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/08 16:18:33  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/08 16:18:33  -- Duplicates: 3
2026/03/08 16:18:33 (Done)



```

<h2>Full output from simple run with log output.</h2>

```python
$ hashcomparefiles -log default  testdata/
2026/03/08 16:19:44 Filecount after first pass: 17
2026/03/08 16:19:44 Filecount after second pass: 12
2026/03/08 16:19:44 Filecount after third pass: 12
2026/03/08 16:19:44 Groups of duplicates after shrink: 4
2026/03/08 16:19:44 Files with hash: c0f5efbef0fe98aa90619444250b1a5eb23158d6686f0b190838f3d544ec85b9
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt size: 10
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt size: 10
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt size: 10
2026/03/08 16:19:44  -- Duplicates: 3
2026/03/08 16:19:44 Files with hash: 7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileB.txt size: 20
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt size: 20
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt size: 20
2026/03/08 16:19:44  -- Duplicates: 3
2026/03/08 16:19:44 Files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileC.txt size: 22
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/08 16:19:44  -- Duplicates: 3
2026/03/08 16:19:44 Files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileDDup.txt size: 30
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt size: 30
2026/03/08 16:19:44  -- Duplicates: 3
2026/03/08 16:19:44 (Done)
$ more log.log 
2026/03/08 16:19:44 Filecount after first pass: 17
2026/03/08 16:19:44 Filecount after second pass: 12
2026/03/08 16:19:44 Filecount after third pass: 12
2026/03/08 16:19:44 Groups of duplicates after shrink: 4
2026/03/08 16:19:44 Files with hash: c0f5efbef0fe98aa90619444250b1a5eb23158d6686f0b190838f3d544ec85b9
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt size: 10
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt size: 10
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt size: 10
2026/03/08 16:19:44  -- Duplicates: 3
2026/03/08 16:19:44 Files with hash: 7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileB.txt size: 20
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt size: 20
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt size: 20
2026/03/08 16:19:44  -- Duplicates: 3
2026/03/08 16:19:44 Files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileC.txt size: 22
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/08 16:19:44  -- Duplicates: 3
2026/03/08 16:19:44 Files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileDDup.txt size: 30
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/08 16:19:44  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt size: 30
2026/03/08 16:19:44  -- Duplicates: 3
2026/03/08 16:19:44 (Done)

```


<h2> Running with log on my home directory, note that the full run took less than 1 minute.</h2>

```python
$ hashcomparefiles fullhome.log /home/gainax2k1/
 \ Files processed: 412300
2026/03/08 16:21:37 Filecount after first pass: 412379
2026/03/08 16:21:55 Filecount after second pass: 382210
2026/03/08 16:22:20 Filecount after third pass: 298294
2026/03/08 16:22:20 Groups of duplicates after shrink: 83666

< cut for ammount of output >

$ more fullhome.log 
2026/03/08 16:21:37 Filecount after first pass: 412379
2026/03/08 16:21:55 Filecount after second pass: 382210
2026/03/08 16:22:20 Filecount after third pass: 298294
2026/03/08 16:22:20 Groups of duplicates after shrink: 83666

< cut for amount of output >

2026/03/08 16:22:24 Files with hash: 009a3b2574c24b179adfab8d8bfccc58b780d5c65a6cc306823175e5aa2eccf1
2026/03/08 16:22:24  - /home/gainax2k1/.local/share/Steam/steamapps/common/SteamLinuxRuntime_sniper/sniper_platform_3.0.20260119.200241/files/share/mime/image/cgm.xml size: 1326
2026/03/08 16:22:24  - /home/gainax2k1/.local/share/Steam/steamapps/common/SteamLinuxRuntime_sniper/var/tmp-20VPL3/usr/share/mime/image/cgm.xml size: 1326
2026/03/08 16:22:24  - /home/gainax2k1/.local/share/Steam/steamrt64/pv-runtime/steam-runtime-steamrt/steamrt3c_platform_3c.0.20251202.187499/files/share/mime/image/cgm.xml size: 1326
2026/03/08 16:22:24  - /home/gainax2k1/.local/share/Steam/steamrt64/pv-runtime/steam-runtime-steamrt/var/tmp-M3M0K3/usr/share/mime/image/cgm.xml size: 1326
2026/03/08 16:22:24  - /home/gainax2k1/.local/share/Steam/steamrt64/steam-runtime-steamrt/var/tmp-SCQ1H3/usr/share/mime/image/cgm.xml size: 1326
2026/03/08 16:22:24  -- Duplicates: 5
2026/03/08 16:22:24 Files with hash: 6ba12c6b4b5fec8df3b5e6771c108eb7ef080561ab5fc2f5206b2b35f9e498bb
2026/03/08 16:22:24  - /home/gainax2k1/.local/share/flatpak/repo/objects/6f/7b0ffebd1fee717dec30b2a6cc0b1869bee25e426c1dabfc7e12c966d902ea.file size: 2312
2026/03/08 16:22:24  - /home/gainax2k1/.local/share/flatpak/runtime/org.gnome.Platform/x86_64/48/c7cdf5b885ebc70910207b6a4d62ffca3458036edc14d0972ca090b3797b52c3/files/share/icons/Adwaita/symbolic/status/network-wireless-offline-symbolic.svg size: 2312
2026/03/08 16:22:24  -- Duplicates: 2
2026/03/08 16:22:24 Files with hash: 28e5515b7f5b09ac8794c315822028e2ca1aaac7f8adb261e2257339cfd6e03e
2026/03/08 16:22:24  - /home/gainax2k1/.local/share/flatpak/repo/objects/6d/e675ea1da1ddb2f90902e3566dab3ec8c724890efbb36f91e3ac83e5d4e2b4.file size: 24637
2026/03/08 16:22:24  - /home/gainax2k1/.local/share/flatpak/runtime/org.freedesktop.Platform/x86_64/25.08/6482ce412b0584ab2e2191db1c1de27b7072b8945c20e83a661d284b9c10e6d4/files/lib/python3.13/site-packages/setuptools/__pycache__/build_meta.cpython-313.pyc size: 24637
2026/03/08 16:22:24  -- Duplicates: 2
2026/03/08 16:22:24 (Done)
```

<h2> Piping in list of files and folders, then selectively removing them: </h2>

```python
$ cat testFilesList.txt | hashcomparefiles -remove
2026/03/08 16:32:55 Filecount after first pass: 11
2026/03/08 16:32:55 Filecount after second pass: 9
2026/03/08 16:32:55 Filecount after third pass: 9
2026/03/08 16:32:55 Groups of duplicates after shrink: 4
2026/03/08 16:32:55 Files with hash: c0f5efbef0fe98aa90619444250b1a5eb23158d6686f0b190838f3d544ec85b9
2026/03/08 16:32:55  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt size: 10
2026/03/08 16:32:55  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt size: 10
2026/03/08 16:32:55  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt size: 10
2026/03/08 16:32:55  -- Duplicates: 3
Delete file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > d
2026/03/08 16:33:01 Deleted duplicate file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt
Delete file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > t
2026/03/08 16:33:05 Deleted duplicate file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt
Delete file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > s
2026/03/08 16:33:07 Files with hash: 7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f
2026/03/08 16:33:07  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt size: 20
2026/03/08 16:33:07  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt size: 20
2026/03/08 16:33:07  -- Duplicates: 2
Delete file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > c
2026/03/08 16:33:08 Files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/08 16:33:08  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/08 16:33:08  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/08 16:33:08  -- Duplicates: 2
Delete file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > c
2026/03/08 16:33:13 Files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/08 16:33:13  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/08 16:33:13  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt size: 30
2026/03/08 16:33:13  -- Duplicates: 2
Delete file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > s
Delete file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > d
2026/03/08 16:33:20 Deleted duplicate file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt
2026/03/08 16:33:20 (Done)

```

