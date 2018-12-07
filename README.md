# safetar
Reads a possibly truncated tar file from stdin and writes a non-truncated tar to stdout

The following takes an archive where the last file is truncated (incomplete download?) 
and creates a new file with everything but that truncated file. 

```cat archive-incomplete.tar | safetar > archive-only-entire-files.tar```

The usefulness of this is admittedly narrow, the main purpose is to be used in pipelines
where (over)?writing a file that  is not complete to the output stream is undesirable.  
In my testing existing tar commands will write a full size file with as much data as 
possible and the rest of the data as null bytes which is not always what you want.
