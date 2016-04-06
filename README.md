Siemens MRI Raw Data Tools
---------------------------

A couple of simple tools to create a (somewhat) unique signature of `*.dat` files and DICOM files from siemens MRI systems

To install:

    go get github.com/hansenms/siemensraw/...


To get a signature from a datfile:

    datsignature <DATFILENAME>

To get a signature from a DICOM:

    dicomsignature <DICOMFILE>


If signatures from DICOM and `*.dat` file match, there is a pretty good chance that the images were generated from the `*.dat` file.

The package also provides a tool for extracting any buffer from Siemens `*.dat` files:

    datbuffer -b Phoenix <DATFILENAME>

Finally there are some tools for creating a filesystem based database for raw data files:

    addfile -b <BASEPATH TO BE REMOVED> -d <DESTINATION FOLDER> <DATFILENAME>

To add an entire folder with lots of `*.dat` files to the database:

   find <FOLDER TO SEARCH> -name "*.dat" -exec addfile -b /mnt/cnmc/ -d <DATABASE PATH> {} \;

To subsequently locate a particular `*.dat` file based on signature:

   findraw -d <DATABASE PATH> <SIGNATURE>

The output is in json format and you can use a tool like `jq` to extract info from the output:

    findraw -d <DATABASE PATH> <SIGNATURE> | jq .


